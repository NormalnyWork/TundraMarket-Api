package order

import "context"

type OrderRepository interface {
	Save(ctx context.Context, o *Order) (*Order, error)
	GetByID(ctx context.Context, id int32) (*Order, error)
}
