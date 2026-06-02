package product

import (
	"context"

	domainproduct "tundraMarket/internal/domain/product"
	"tundraMarket/internal/infrastructure/postgres"
	sqlcdb "tundraMarket/internal/infrastructure/postgres/sqlc"
)

type ProductRepo struct {
	q *sqlcdb.Queries
}

func NewProductRepo(q *sqlcdb.Queries) *ProductRepo {
	return &ProductRepo{q: q}
}

func (r *ProductRepo) GetAll(ctx context.Context) ([]*domainproduct.Product, error) {
	rows, err := r.q.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	return productsToDomain(rows), nil
}

func (r *ProductRepo) GetByIDs(ctx context.Context, ids []int32) ([]*domainproduct.Product, error) {
	rows, err := r.q.GetProductsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return productsToDomain(rows), nil
}

func productsToDomain(rows []sqlcdb.Product) []*domainproduct.Product {
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
	return result
}
