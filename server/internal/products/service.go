package products

import (
	"context"
	"errors"
	"math/big"
	repo "server/internal/adapters/postgresql/sqlc"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	CreateProduct(
		ctx context.Context,
		tempProduct CreateProductParams,
	) (repo.Product, error)
}

// upper vala tou interface hy
// neechy vala struct hy aur vo struct interface implement kr rhy hn
// so may ab handler may service ko use kru ga
type svc struct {
	repo repo.Querier
}

// service svc is liye return kr pa rha hy bcoz svc nay service kay funcs ko implement kia hua h
func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.GetProducts(ctx)
}

func (s *svc) CreateProduct(
	ctx context.Context,
	tempProduct CreateProductParams,
) (repo.Product, error) {

	if strings.TrimSpace(tempProduct.Name) == "" {
		return repo.Product{}, errors.New("name is required")
	}

	if tempProduct.Price <= 0 {
		return repo.Product{}, errors.New("price must be greater than 0")
	}

	if tempProduct.Quantity < 0 {
		return repo.Product{}, errors.New("quantity cannot be negative")
	}

	price := pgtype.Numeric{
		Int:   big.NewInt(int64(tempProduct.Price)),
		Exp:   0,
		Valid: true,
	}

	return s.repo.CreateProduct(ctx, repo.CreateProductParams{
		Name:     tempProduct.Name,
		Price:    price,
		Quantity: tempProduct.Quantity,
	})
}
