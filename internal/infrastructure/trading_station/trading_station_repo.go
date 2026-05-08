package trading_station

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

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
			postgres.NumericToFloat32(row.Longitude),
			postgres.NumericToFloat32(row.Latitude),
		)
	}
	return result, nil
}

func (r *TradingStationRepo) GetByID(ctx context.Context, id int32) (*tradingstation.TradingStation, error) {
	row, err := r.q.GetTradingStationByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tradingstation.ErrNotFound
		}
		return nil, err
	}
	return tradingstation.New(
		row.ID,
		postgres.TextToString(row.Name),
		postgres.TextToStringPtr(row.Phone),
		postgres.NumericToFloat32(row.Longitude),
		postgres.NumericToFloat32(row.Latitude),
	), nil
}

func (r *TradingStationRepo) SetPhone(ctx context.Context, id int32, phone string) (*tradingstation.TradingStation, error) {
	row, err := r.q.SetTradingStationPhone(ctx, sqlcdb.SetTradingStationPhoneParams{
		ID: id,
		Phone: pgtype.Text{
			String: phone,
			Valid:  true,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tradingstation.ErrNotFound
		}
		return nil, err
	}
	return tradingstation.New(
		row.ID,
		postgres.TextToString(row.Name),
		postgres.TextToStringPtr(row.Phone),
		postgres.NumericToFloat32(row.Longitude),
		postgres.NumericToFloat32(row.Latitude),
	), nil
}
