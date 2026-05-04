package product

import (
	"context"

	domainproduct "tundraMarket/internal/domain/product"
	"tundraMarket/internal/infrastructure/postgres"
	sqlcdb "tundraMarket/internal/infrastructure/postgres/sqlc"
)

type Repo struct {
	q *sqlcdb.Queries
}

func NewRepo(q *sqlcdb.Queries) *Repo {
	return &Repo{q: q}
}

func (r *Repo) GetAll(ctx context.Context) ([]*domainproduct.Product, error) {
	rows, err := r.q.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*domainproduct.Product, len(rows))
	for i, row := range rows {
		result[i] = domainproduct.New(
			row.ID,
			row.Name,
			postgres.TextToStringPtr(row.Details),
			postgres.Int4ToInt32(row.Weight),
			postgres.Int4ToInt32(row.Volume),
		)
	}
	return result, nil
}
