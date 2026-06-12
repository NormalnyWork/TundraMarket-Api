package admin

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	domainadmin "tundraMarket/internal/domain/admin"
	sqlcdb "tundraMarket/internal/infrastructure/postgres/sqlc"
)

type AdminRepo struct {
	q *sqlcdb.Queries
}

func NewAdminRepo(q *sqlcdb.Queries) *AdminRepo {
	return &AdminRepo{q: q}
}

func (r *AdminRepo) GetByLogin(ctx context.Context, login string) (*domainadmin.Admin, error) {
	row, err := r.q.GetAdminByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}

	return domainadmin.New(
		row.ID,
		row.Login,
		row.Password,
	), nil
}
