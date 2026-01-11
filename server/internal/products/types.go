package products

type Product struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Price    int32  `json:"price"`
	Quantity int32  `json:"quantity"`
}
type CreateProductParams struct {
	Name     string `json:"name"`
	Price    int32  `json:"price"`
	Quantity int32  `json:"quantity"`
}
