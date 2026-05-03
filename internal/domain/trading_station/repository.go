package trading_station

import "context"

type TradingStationRepository interface {
	GetAll(ctx context.Context) ([]*TradingStation, error)
}
