package trading_station

import (
	"context"
	tradingstation "tundraMarket/internal/domain/trading_station"
	"tundraMarket/internal/infrastructure/postgres"
	sqlcdb "tundraMarket/internal/infrastructure/postgres/sqlc"
)

type TradingStationRepo struct {
	q *sqlcdb.Queries
}

func NewTradingStationRepo(q *sqlcdb.Queries) *TradingStationRepo {
	return &TradingStationRepo{q: q}
}

func (r *TradingStationRepo) GetAll(ctx context.Context) ([]*tradingstation.TradingStation, error) {
	rows, err := r.q.GetAllTradingStations(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*tradingstation.TradingStation, len(rows))
	for i, row := range rows {
		result[i] = tradingstation.New(
			row.ID,
			postgres.TextToString(row.Name),
			postgres.TextToStringPtr(row.Phone),
			postgres.NumericToFloat64(row.Longitude),
			postgres.NumericToFloat64(row.Latitude),
		)
	}
	return result, nil
}
