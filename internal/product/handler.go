package product

import (
	"net/http"

	"github.com/citadel-corp/eniqilo-store/internal/common/request"
	"github.com/citadel-corp/eniqilo-store/internal/common/response"
	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductPayload

	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	err = req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}

	userResp, err := h.service.Create(r.Context(), req)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusCreated, response.ResponseBody{
		Message: "Product created successfully",
		Data:    userResp,
	})
}

func (h *Handler) EditProduct(w http.ResponseWriter, r *http.Request) {
	var req EditProductPayload

	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	params := mux.Vars(r)
	req.ID = params["id"]

	err = req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}

	err = h.service.Edit(r.Context(), req)
	if err == ErrProductNotFound {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Not found",
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
		Message: "Product edited successfully",
	})
}
