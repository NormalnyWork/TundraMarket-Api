package product

import "context"

type Repository interface {
	GetAll(ctx context.Context) ([]*Product, error)
	GetByIDs(ctx context.Context, ids []int32) ([]*Product, error)
}
