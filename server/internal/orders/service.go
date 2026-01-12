package orders

import (
	"context"
	"fmt"
	repo "server/internal/adapters/postgresql/sqlc"
	"server/internal/products"

	"github.com/jackc/pgx/v5"
)

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder CreateOrderParams) (repo.Order, error)
	GetOrdersByCustomerID(ctx context.Context, customerID int32) ([]OrderWithItems, error)
}

type svc struct {
	repo repo.Querier
	db   *pgx.Conn
}

func NewService(repo repo.Querier, db *pgx.Conn) Service {
	return &svc{repo: repo, db: db}
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder CreateOrderParams) (repo.Order, error) {

	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer id is required")
	}
	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("items are required")
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)
	qtx := repo.New(tx)
	order, err := qtx.CreateOrder(ctx, int32(tempOrder.CustomerID))
	if err != nil {
		return repo.Order{}, err
	}

	for _, item := range tempOrder.Items {
		product, err := qtx.FindProductByID(ctx, int32(item.ProductID))
		if err != nil {
			return repo.Order{}, err
		}
		if product.Quantity < int32(item.Quantity) {
			return repo.Order{}, fmt.Errorf("product %d is out of stock", item.ProductID)
		}

		price := int32(0)
		if product.Price.Valid && !product.Price.NaN && product.Price.Int != nil {
			price = int32(product.Price.Int.Int64())
		}
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     int32(item.Quantity) * int32(price),
		})
		if err != nil {
			return repo.Order{}, err
		}
	}
	tx.Commit(ctx)
	return order, nil
}

func (s *svc) GetOrdersByCustomerID(ctx context.Context, customerID int32) ([]OrderWithItems, error) {
	// 1️⃣ Get all orders
	orders, err := s.repo.GetOrdersByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("fetching orders: %w", err)
	}

	// 2️⃣ Collect all product IDs
	productIDsMap := make(map[int32]struct{})
	orderItemsMap := make(map[int32][]repo.GetOrderItemsByOrderIDRow) // map[orderID]items
	for _, o := range orders {
		itemsRows, err := s.repo.GetOrderItemsByOrderID(ctx, o.ID)
		if err != nil {
			return nil, fmt.Errorf("fetching items for order %d: %w", o.ID, err)
		}
		orderItemsMap[o.ID] = itemsRows
		for _, i := range itemsRows {
			productIDsMap[i.ProductID] = struct{}{}
		}
	}

	productIDs := make([]int32, 0, len(productIDsMap))
	for id := range productIDsMap {
		productIDs = append(productIDs, id)
	}

	// 3️⃣ Fetch all products in bulk
	productsList, err := s.repo.FindProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("fetching products: %w", err)
	}

	productMap := make(map[int32]products.Product)
	for _, p := range productsList {
		price := int32(0)
		if p.Price.Valid && !p.Price.NaN && p.Price.Int != nil {
			price = int32(p.Price.Int.Int64())
		}
		productMap[p.ID] = products.Product{
			ID:       p.ID,
			Name:     p.Name,
			Price:    price,
			Quantity: p.Quantity,
			Image:    p.Image,
		}
	}

	// 4️⃣ Build response
	result := make([]OrderWithItems, 0, len(orders))
	for _, o := range orders {
		itemsRows := orderItemsMap[o.ID]
		items := make([]OrderItem, 0, len(itemsRows))
		for _, i := range itemsRows {
			p, ok := productMap[i.ProductID]
			if !ok {
				return nil, fmt.Errorf("product %d not found", i.ProductID)
			}
			items = append(items, OrderItem{
				ID:       i.ID,
				Product:  p,
				Quantity: i.Quantity,
				Price:    i.Price,
			})
		}

		result = append(result, OrderWithItems{
			ID:        o.ID,
			CreatedAt: o.CreatedAt.Time,
			Items:     items,
		})
	}

	return result, nil
}
