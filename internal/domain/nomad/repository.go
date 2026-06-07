package nomad

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("nomad not found")

type Repository interface {
	GetByPhone(ctx context.Context, phone string) (*Nomad, error)
	Create(ctx context.Context, phone string) (*Nomad, error)
	GetByID(ctx context.Context, id int32) (*Nomad, error)
}
