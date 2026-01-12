package products

import (
	"net/http"
	"server/internal/request"
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

	response.WriteJson(w, http.StatusOK, response.GernalResponse{Success: true, Message: "Products fetched successfully", Data: products})
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var tempProduct CreateProductParams
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
