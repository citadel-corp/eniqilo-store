package checkout

import (
	"net/http"

	"github.com/citadel-corp/eniqilo-store/internal/common/response"
	"github.com/gorilla/schema"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
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
