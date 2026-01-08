package products

import (
	"net/http"
	"server/internal/response"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GernalResponse{Success: false, Message: err.Error(), Data: nil})
		return
	}

	response.WriteJson(w, http.StatusOK, response.GernalResponse{Success: true, Message: "success", Data: products})
}
