package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	domainorder "tundraMarket/internal/domain/order"
	"tundraMarket/internal/infrastructure/postgres"
	sqlcdb "tundraMarket/internal/infrastructure/postgres/sqlc"
)

type OrderRepo struct {
	pool *pgxpool.Pool
	q    *sqlcdb.Queries
}

func NewOrderRepo(pool *pgxpool.Pool, q *sqlcdb.Queries) *OrderRepo {
	return &OrderRepo{pool: pool, q: q}
}

func (r *OrderRepo) Save(ctx context.Context, o *domainorder.Order) (*domainorder.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.q.WithTx(tx)

	row, err := qtx.CreateOrder(ctx, sqlcdb.CreateOrderParams{
		NomadID:          pgtype.Int4{Int32: o.NomadID(), Valid: true},
		TradingStationID: pgtype.Int4{Int32: o.TradingStationID(), Valid: true},
		Longitude:        postgres.Float32ToNumeric(o.Longitude()),
		Latitude:         postgres.Float32ToNumeric(o.Latitude()),
		Comment:          pgtype.Text{String: o.Comment(), Valid: o.Comment() != ""},
	})
	if err != nil {
		return nil, err
	}

	for _, p := range o.Products() {
		err = qtx.AddProductToOrder(ctx, sqlcdb.AddProductToOrderParams{
			OrdersID:  row.ID,
			ProductID: p.ProductID,
			Quantity:  p.Quantity,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return rowToDomain(row, o.Products()), nil
}

func (r *OrderRepo) GetByID(ctx context.Context, id int32) (*domainorder.Order, error) {
	row, err := r.q.GetOrderById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainorder.ErrInvalidId
		}
		return nil, err
	}

	products, err := r.q.GetProductsByOrderID(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	items := make([]domainorder.ProductCount, len(products))
	for i, p := range products {
		items[i] = domainorder.ProductCount{
			ProductID: p.ID,
			Quantity:  p.Quantity,
		}
	}

	return rowToDomain(row, items), nil
}

func rowToDomain(row sqlcdb.Order, products []domainorder.ProductCount) *domainorder.Order {
	return domainorder.Restore(
		row.ID,
		postgres.Int4ToInt32(row.NomadID),
		postgres.Int4ToInt32(row.TradingStationID),
		domainorder.Status(row.Status),
		postgres.TextToString(row.Comment),
		postgres.NumericToFloat32(row.Longitude),
		postgres.NumericToFloat32(row.Latitude),
		products,
		row.CreatedAt.Time,
	)
}
