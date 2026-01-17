// payload types that are used in the request and response bodies
package orders

import (
	"time"
)

type orderItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}
type CreateOrderParams struct {
	CustomerID int         `json:"customerId"`
	Items      []orderItem `json:"items"`
}

type Product struct {
	ID    int32   `json:"id"`
	Name  string  `json:"name"`
	Image *string `json:"image"`
}

// OrderItem represents one item in an order
type OrderItem struct {
	ID       int32   `json:"id"`
	Product  Product `json:"product"`
	Quantity int32   `json:"quantity"`
	Price    int32   `json:"price"`
}

// OrderWithItems represents an order with all its items
type OrderWithItems struct {
	ID        int32       `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	Items     []OrderItem `json:"items"`
}
