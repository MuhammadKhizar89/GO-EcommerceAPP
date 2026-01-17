package handlers

import (
	"net/http"
	products "server/internal/domain/product"
	"server/internal/util/request"
	"server/internal/util/response"
)

type productHandler struct {
	service products.Service
}

func NewProductHandler(service products.Service) *productHandler {
	return &productHandler{service}
}

func (h *productHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GernalResponse{Success: false, Message: err.Error(), Data: nil})
		return
	}

	response.WriteJson(w, http.StatusOK, response.GernalResponse{Success: true, Message: "Products fetched successfully", Data: products})
}

func (h *productHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var tempProduct products.CreateProductParams
	if err := request.ReadJSON(r, &tempProduct); err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GernalResponse{Success: false, Message: err.Error(), Data: nil})
		return
	}
	product, err := h.service.CreateProduct(r.Context(), tempProduct)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GernalResponse{Success: false, Message: err.Error(), Data: nil})
		return
	}
	response.WriteJson(w, http.StatusOK, response.GernalResponse{Success: true, Message: "Product created successfully", Data: product})
}
