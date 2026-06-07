package order

import "context"

type OrderRepository interface {
	Save(ctx context.Context, o *Order) (*Order, error)
	GetByID(ctx context.Context, id int32) (*Order, error)
	ChangeStatus(ctx context.Context, id int32, status Status, comment *string) (*Order, error)
	ListByNomad(ctx context.Context, nomadID int32, category OrderCategory, anchor, pageSize int32) ([]*Order, error)
	ListByTradingStation(ctx context.Context, tradingStationID int32, category OrderCategory, anchor, pageSize int32) ([]*Order, error)
	GetUpdatesByNomad(ctx context.Context, nomadID int32, afterUnix int64) ([]*Order, error)
	GetUpdatesByTradingStation(ctx context.Context, tradingStationID int32, afterUnix int64) ([]*Order, error)
}
