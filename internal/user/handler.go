package user

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

func (h *Handler) CreateStaff(w http.ResponseWriter, r *http.Request) {
	var req CreateStaffPayload

	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}
	userResp, err := h.service.CreateStaff(r.Context(), req)
	if errors.Is(err, ErrPhoneNumberAlreadyExists) {
		response.JSON(w, http.StatusConflict, response.ResponseBody{
			Message: "User already exists",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrValidationFailed) {
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
	response.JSON(w, http.StatusCreated, response.ResponseBody{
		Message: "User registered successfully",
		Data:    userResp,
	})
}

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req CreateCustomerPayload

	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}
	userResp, err := h.service.CreateCustomer(r.Context(), req)
	if errors.Is(err, ErrPhoneNumberAlreadyExists) {
		response.JSON(w, http.StatusConflict, response.ResponseBody{
			Message: "User already exists",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrValidationFailed) {
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
	response.JSON(w, http.StatusCreated, response.ResponseBody{
		Message: "User registered successfully",
		Data:    userResp,
	})
}

func (h *Handler) StaffLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginPayload

	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}
	userResp, err := h.service.StaffLogin(r.Context(), req)
	if errors.Is(err, ErrUserNotFound) {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Not found",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrWrongPassword) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrValidationFailed) {
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
		Message: "User logged successfully",
		Data:    userResp,
	})
}

func (h *Handler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	newSchema := schema.NewDecoder()
	newSchema.IgnoreUnknownKeys(true)

	var req ListCustomerPayload
	if err := newSchema.Decode(&req, r.URL.Query()); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{})
		return
	}

	res, err := h.service.ListCustomers(r.Context(), req)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, response.ResponseBody{
		Message: "success",
		Data:    res,
	})
}
