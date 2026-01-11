// payload types that are used in the request and response bodies
package orders

type orderItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
	Price     int `json:"price"`
}
type CreateOrderParams struct {
	CustomerID int         `json:"customerId"`
	Items      []orderItem `json:"items"`
}
