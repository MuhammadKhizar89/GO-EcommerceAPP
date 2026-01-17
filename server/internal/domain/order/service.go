package orders

import (
	"context"
	"fmt"
	repo "server/internal/adapters/postgresql/sqlc"

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
		return repo.Order{}, fmt.Errorf("customer id not found")
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
		if item.Quantity == 0 {
			return repo.Order{}, fmt.Errorf("quantity is required")
		}
		if item.Quantity > 10 {
			return repo.Order{}, fmt.Errorf("quantity cannot be greater than 10")
		}
		if product.Quantity < int32(item.Quantity) {
			return repo.Order{}, fmt.Errorf("product %d is out of stock", item.ProductID)
		}

		price := int32(0)
		if product.Price.Valid && !product.Price.NaN && product.Price.Int != nil {
			price = int32(product.Price.Int.Int64())
		}
		_, err = qtx.UpdateProduct(ctx, repo.UpdateProductParams{
			ID:       int32(item.ProductID),
			Price:    product.Price,
			Name:     product.Name,
			Image:    product.Image,
			Quantity: product.Quantity - int32(item.Quantity),
		})
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

	orders, err := s.repo.GetOrdersByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("fetching orders: %w", err)
	}

	productIDsMap := make(map[int32]struct{})
	orderItemsMap := make(map[int32][]repo.GetOrderItemsByOrderIDRow)
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

	productsList, err := s.repo.FindProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("fetching products: %w", err)
	}

	productMap := make(map[int32]Product)
	for _, p := range productsList {
		productMap[p.ID] = Product{
			ID:    p.ID,
			Name:  p.Name,
			Image: p.Image,
		}
	}

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
