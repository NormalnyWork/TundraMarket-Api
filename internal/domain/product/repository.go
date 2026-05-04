package product

import "context"

type Repository interface {
	GetAll(ctx context.Context) ([]*Product, error)
}
