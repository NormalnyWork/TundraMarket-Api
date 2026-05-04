package nomad

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	domainnomad "tundraMarket/internal/domain/nomad"
	sqlcdb "tundraMarket/internal/infrastructure/postgres/sqlc"
)

type Repo struct {
	q *sqlcdb.Queries
}

func NewRepo(q *sqlcdb.Queries) *Repo {
	return &Repo{q: q}
}

func (r *Repo) GetByPhone(ctx context.Context, phone string) (*domainnomad.Nomad, error) {
	row, err := r.q.GetNomadByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainnomad.ErrNotFound
		}
		return nil, err
	}
	return domainnomad.New(row.ID, row.Phone), nil
}

func (r *Repo) Create(ctx context.Context, phone string) (*domainnomad.Nomad, error) {
	row, err := r.q.CreateNomad(ctx, phone)
	if err != nil {
		return nil, err
	}
	return domainnomad.New(row.ID, row.Phone), nil
}
