package checkout

import (
	"errors"
	"net/http"

	"github.com/citadel-corp/eniqilo-store/internal/common/request"
	"github.com/citadel-corp/eniqilo-store/internal/common/response"
	"github.com/gorilla/schema"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CheckoutProducts(w http.ResponseWriter, r *http.Request) {
	var req CheckoutRequest

	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	err = h.service.CheckoutProducts(r.Context(), req)
	if errors.Is(err, ErrCustomerNotFound) ||
		errors.Is(err, ErrProductNotFound) {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Not found",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrValidationFailed) ||
		errors.Is(err, ErrProductUnavailable) ||
		errors.Is(err, ErrProductStockNotEnough) ||
		errors.Is(err, ErrNotEnoughMoney) ||
		errors.Is(err, ErrWrongChange) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, response.ResponseBody{
		Message: "success",
	})
}

func (h *Handler) ListCheckoutHistories(w http.ResponseWriter, r *http.Request) {
	newSchema := schema.NewDecoder()
	newSchema.IgnoreUnknownKeys(true)

	var req ListCheckoutHistoriesPayload
	if err := newSchema.Decode(&req, r.URL.Query()); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{})
		return
	}

	histories, err := h.service.ListCheckoutHistories(r.Context(), req)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, response.ResponseBody{
		Message: "success",
		Data:    histories,
	})
}
