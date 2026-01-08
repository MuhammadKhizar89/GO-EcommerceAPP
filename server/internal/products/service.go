package products

import (
	"context"
	repo "server/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
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
