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

	if err = qtx.AddStatusHistory(ctx, sqlcdb.AddStatusHistoryParams{
		OrdersID: row.ID,
		Status:   statusToNullStatus(domainorder.StatusCreated),
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, row.ID)
}

func (r *OrderRepo) GetByID(ctx context.Context, id int32) (*domainorder.Order, error) {
	row, err := r.q.GetOrderById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainorder.ErrInvalidId
		}
		return nil, err
	}

	return r.rowToDomain(ctx, row)
}

func (r *OrderRepo) ChangeStatus(ctx context.Context, id int32, status domainorder.Status, comment *string) (*domainorder.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.q.WithTx(tx)

	row, err := qtx.UpdateOrderStatus(ctx, sqlcdb.UpdateOrderStatusParams{
		ID:     id,
		Status: statusToSQLC(status),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainorder.ErrInvalidId
		}
		return nil, err
	}

	if err = qtx.AddStatusHistory(ctx, sqlcdb.AddStatusHistoryParams{
		OrdersID: row.ID,
		Status:   statusToNullStatus(status),
		Comment:  textFromPointer(comment),
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, row.ID)
}

func (r *OrderRepo) ListByNomad(ctx context.Context, nomadID int32, category domainorder.OrderCategory, anchor, pageSize int32) ([]*domainorder.Order, error) {
	rows, err := r.q.GetOrdersByNomadIDAndCategory(ctx, sqlcdb.GetOrdersByNomadIDAndCategoryParams{
		NomadID:       pgtype.Int4{Int32: nomadID, Valid: true},
		Limit:         pageSize,
		OrderCategory: string(category),
		Anchor:        anchor,
	})
	if err != nil {
		return nil, err
	}

	return r.rowsToDomain(ctx, rows)
}

func (r *OrderRepo) ListByTradingStation(ctx context.Context, tradingStationID int32, category domainorder.OrderCategory, anchor, pageSize int32) ([]*domainorder.Order, error) {
	rows, err := r.q.GetOrdersByStationAndCategory(ctx, sqlcdb.GetOrdersByStationAndCategoryParams{
		TradingStationID: pgtype.Int4{Int32: tradingStationID, Valid: true},
		Limit:            pageSize,
		OrderCategory:    string(category),
		Anchor:           anchor,
	})
	if err != nil {
		return nil, err
	}

	return r.rowsToDomain(ctx, rows)
}

func (r *OrderRepo) GetUpdatesByNomad(ctx context.Context, nomadID int32, afterUnix int64) ([]*domainorder.Order, error) {
	rows, err := r.q.GetOrdersByNomadIDUpdatedAfter(ctx, sqlcdb.GetOrdersByNomadIDUpdatedAfterParams{
		NomadID:     pgtype.Int4{Int32: nomadID, Valid: true},
		ToTimestamp: float64(afterUnix),
	})
	if err != nil {
		return nil, err
	}

	return r.rowsToDomain(ctx, rows)
}

func (r *OrderRepo) GetUpdatesByTradingStation(ctx context.Context, tradingStationID int32, afterUnix int64) ([]*domainorder.Order, error) {
	rows, err := r.q.GetOrdersByStationUpdatedAfter(ctx, sqlcdb.GetOrdersByStationUpdatedAfterParams{
		TradingStationID: pgtype.Int4{Int32: tradingStationID, Valid: true},
		ToTimestamp:      float64(afterUnix),
	})
	if err != nil {
		return nil, err
	}

	return r.rowsToDomain(ctx, rows)
}

func (r *OrderRepo) rowsToDomain(ctx context.Context, rows []sqlcdb.Order) ([]*domainorder.Order, error) {
	orders := make([]*domainorder.Order, len(rows))
	for i, row := range rows {
		order, err := r.rowToDomain(ctx, row)
		if err != nil {
			return nil, err
		}
		orders[i] = order
	}

	return orders, nil
}

func (r *OrderRepo) rowToDomain(ctx context.Context, row sqlcdb.Order) (*domainorder.Order, error) {
	products, err := r.q.GetProductsByOrderID(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	historyRows, err := r.q.GetStatusHistoryByOrderID(ctx, row.ID)
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

	history := make([]domainorder.StatusHistory, len(historyRows))
	for i, item := range historyRows {
		history[i] = domainorder.NewStatusHistory(domainorder.Status(item.Status.Status), item.CreatedAt.Time, textToPointer(item.Comment))
	}

	return domainorder.Restore(
		row.ID,
		postgres.Int4ToInt32(row.NomadID),
		postgres.Int4ToInt32(row.TradingStationID),
		domainorder.Status(row.Status),
		postgres.TextToString(row.Comment),
		postgres.NumericToFloat32(row.Longitude),
		postgres.NumericToFloat32(row.Latitude),
		items,
		history,
		row.CreatedAt.Time,
	), nil
}

func statusToSQLC(status domainorder.Status) sqlcdb.Status {
	return sqlcdb.Status(status)
}

func statusToNullStatus(status domainorder.Status) sqlcdb.NullStatus {
	return sqlcdb.NullStatus{
		Status: statusToSQLC(status),
		Valid:  true,
	}
}

func textFromPointer(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}

	return pgtype.Text{
		String: *value,
		Valid:  true,
	}
}

func textToPointer(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}

	return &value.String
}

func (r *OrderRepo) GetCurrentByNomadID(ctx context.Context, nomadID int32) (*domainorder.Order, error) {
	row, err := r.q.GetCurrentOrderByNomadID(ctx, pgtype.Int4{Int32: nomadID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainorder.ErrInvalidId
		}
		return nil, err
	}

	return r.rowToDomain(ctx, row)
}
