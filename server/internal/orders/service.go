package orders

import (
	"context"
	"fmt"
	repo "server/internal/adapters/postgresql/sqlc"

	"github.com/jackc/pgx/v5"
)

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder CreateOrderParams) (repo.Order, error)
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
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     int32(item.Price),
		})
		if err != nil {
			return repo.Order{}, err
		}
	}
	tx.Commit(ctx)
	return order, nil
}
