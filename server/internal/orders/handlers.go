package orders

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

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var tempOrder CreateOrderParams
	if err := request.ReadJSON(r, &tempOrder); err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GernalResponse{Success: false, Message: err.Error(), Data: nil})
		return
	}
	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GernalResponse{Success: false, Message: err.Error(), Data: nil})
		return
	}

	response.WriteJson(w, http.StatusOK, response.GernalResponse{Success: true, Message: "success", Data: createdOrder})
}
