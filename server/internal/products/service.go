package products

import "context"

type Service interface {
	ListProducts(ctx context.Context) error
}

// upper vala tou interface hy
// neechy vala struct hy aur vo struct interface implement kr rhy hn
// so may ab handler may service ko use kru ga
type svc struct {
}

// service svc is liye return kr pa rha hy bcoz svc nay service kay funcs ko implement kia hua h
func NewService() Service {
	return &svc{}
}

func (s *svc) ListProducts(ctx context.Context) error {
	return nil
}
