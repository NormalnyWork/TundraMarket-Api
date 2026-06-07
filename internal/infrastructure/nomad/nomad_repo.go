package nomad

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	domainnomad "tundraMarket/internal/domain/nomad"
	sqlcdb "tundraMarket/internal/infrastructure/postgres/sqlc"
)

type NomadRepo struct {
	q *sqlcdb.Queries
}

func NewNomadRepo(q *sqlcdb.Queries) *NomadRepo {
	return &NomadRepo{q: q}
}

func (r *NomadRepo) GetByPhone(ctx context.Context, phone string) (*domainnomad.Nomad, error) {
	row, err := r.q.GetNomadByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainnomad.ErrNotFound
		}
		return nil, err
	}
	return domainnomad.New(row.ID, row.Phone), nil
}

func (r *NomadRepo) Create(ctx context.Context, phone string) (*domainnomad.Nomad, error) {
	row, err := r.q.CreateNomad(ctx, phone)
	if err != nil {
		return nil, err
	}
	return domainnomad.New(row.ID, row.Phone), nil
}

func (r *NomadRepo) GetByID(ctx context.Context, id int32) (*domainnomad.Nomad, error) {
	row, err := r.q.GetNomadByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainnomad.ErrNotFound
		}
		return nil, err
	}
	return domainnomad.New(row.ID, row.Phone), nil
}
