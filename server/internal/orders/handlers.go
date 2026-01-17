package orders

import (
	"net/http"
	"server/internal/middleware"
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
	customerID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.WriteJson(w, http.StatusUnauthorized,
			response.GernalResponse{Success: false, Message: "unauthorized", Data: nil})
		return
	}
	var tempOrder CreateOrderParams
	if err := request.ReadJSON(r, &tempOrder); err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GernalResponse{Success: false, Message: err.Error(), Data: nil})
		return
	}
	tempOrder.CustomerID = int(customerID)

	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GernalResponse{Success: false, Message: err.Error(), Data: nil})
		return
	}
	response.WriteJson(w, http.StatusOK, response.GernalResponse{Success: true, Message: "Order placed successfully", Data: createdOrder})
}
func (h *handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {

	customerID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.WriteJson(w, http.StatusUnauthorized,
			response.GernalResponse{Success: false, Message: "unauthorized", Data: nil})
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
		Message: "Orders fetched successfully",
		Data:    orders,
	})
}
