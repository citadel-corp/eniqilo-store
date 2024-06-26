package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/citadel-corp/eniqilo-store/internal/checkout"
	"github.com/citadel-corp/eniqilo-store/internal/common/db"
	"github.com/citadel-corp/eniqilo-store/internal/common/middleware"
	"github.com/citadel-corp/eniqilo-store/internal/product"
	"github.com/citadel-corp/eniqilo-store/internal/user"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// Connect to database
	// env := os.Getenv("ENV")
	// sslMode := "disable"
	// if env == "production" {
	// 	sslMode = "verify-full sslrootcert=ap-southeast-1-bundle.pem"
	// }
	// connStr := "postgres://[user]:[password]@[neon_hostname]/[dbname]?sslmode=require"
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_PARAMS"))
	// dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	// 	os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), sslMode)
	db, err := db.Connect(connStr)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Cannot connect to database: %v", err))
		os.Exit(1)
	}

	// Create migrations
	// err = db.UpMigration()
	// if err != nil {
	// 	log.Error().Msg(fmt.Sprintf("Up migration failed: %v", err))
	// 	os.Exit(1)
	// }

	// initialize user domain
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	// initialize product domain
	productRepository := product.NewRepository(db)
	productService := product.NewService(productRepository)
	productHandler := product.NewHandler(productService)

	// initialize checkout domain
	checkoutRepository := checkout.NewRepository(db)
	checkoutService := checkout.NewService(checkoutRepository, userRepository, productRepository)
	checkoutHandler := checkout.NewHandler(checkoutService)

	r := mux.NewRouter()
	r.Use(middleware.Logging)
	r.Use(middleware.PanicRecoverer)
	v1 := r.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text")
		io.WriteString(w, "Service ready")
	})

	// staff routes
	sr := v1.PathPrefix("/staff").Subrouter()
	sr.HandleFunc("/register", userHandler.CreateStaff).Methods(http.MethodPost)
	sr.HandleFunc("/login", userHandler.StaffLogin).Methods(http.MethodPost)

	// product routes
	pr := v1.PathPrefix("/product").Subrouter()
	pr.HandleFunc("/customer", middleware.Authenticate(productHandler.ListProductForCustomer)).Methods(http.MethodGet)
	pr.HandleFunc("", middleware.Authorized(productHandler.CreateProduct)).Methods(http.MethodPost)
	pr.HandleFunc("/{id}", middleware.Authorized(productHandler.EditProduct)).Methods(http.MethodPut)
	pr.HandleFunc("/{id}", middleware.Authorized(productHandler.DeleteProduct)).Methods(http.MethodDelete)
	pr.HandleFunc("", middleware.Authorized(productHandler.ListProduct)).Methods(http.MethodGet)

	// product checkout routes
	pcr := pr.PathPrefix("/checkout").Subrouter()
	pcr.HandleFunc("", middleware.Authorized(checkoutHandler.CheckoutProducts)).Methods(http.MethodPost)
	pcr.HandleFunc("/history", middleware.Authorized(checkoutHandler.ListCheckoutHistories)).Methods(http.MethodGet)

	// customer routes
	cr := v1.PathPrefix("/customer").Subrouter()
	cr.HandleFunc("/register", middleware.Authorized(userHandler.CreateCustomer)).Methods(http.MethodPost)
	cr.HandleFunc("", middleware.Authorized(userHandler.ListCustomers)).Methods(http.MethodGet)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Info().Msg(fmt.Sprintf("HTTP server listening on %s", httpServer.Addr))
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Error().Msg(fmt.Sprintf("HTTP server error: %v", err))
		}
		log.Info().Msg("Stopped serving new connections.")
	}()

	// Listen for the termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until termination signal received
	<-stop
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	log.Info().Msg(fmt.Sprintf("Shutting down HTTP server listening on %s", httpServer.Addr))
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Msg(fmt.Sprintf("HTTP server shutdown error: %v", err))
	}
	log.Info().Msg("Shutdown complete.")
}
