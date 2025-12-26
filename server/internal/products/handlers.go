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
	products := struct {
		Products []string `json:"products"`
	}{}
	response.WriteJson(w, http.StatusOK, response.GernalResponse{Success: true, Message: "success", Data: products})
}
