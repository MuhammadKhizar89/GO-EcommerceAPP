package orders

import (
	"net/http"
	"server/internal/request"
	"server/internal/response"
	"strconv"

	"github.com/go-chi/chi/v5"
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
func (h *handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {

	customerIDStr := chi.URLParam(r, "customerId")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GernalResponse{
			Success: false,
			Message: "invalid customer id",
			Data:    nil,
		})
		return
	}

	orders, err := h.service.GetOrdersByCustomerID(r.Context(), int32(customerID))
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GernalResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	response.WriteJson(w, http.StatusOK, response.GernalResponse{
		Success: true,
		Message: "success",
		Data:    orders,
	})
}
