package trading_station

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("trading station not found")

type TradingStationRepository interface {
	GetAll(ctx context.Context) ([]*TradingStation, error)
	GetByID(ctx context.Context, id int32) (*TradingStation, error)
	SetPhone(ctx context.Context, id int32, phone string) (*TradingStation, error)
}
