package product

import (
	commonv1 "tundraMarket/gen/common/v1"
	orderv1 "tundraMarket/gen/order/v1"
	domainproduct "tundraMarket/internal/domain/product"
)

func ToProtoCatalog(products []*domainproduct.Product) *orderv1.UserCatalogOut {
	result := make([]*commonv1.Product, len(products))
	for i, p := range products {
		result[i] = &commonv1.Product{
			Id:      p.ID(),
			Name:    p.Name(),
			Details: p.Details(),
			Weight:  p.Weight(),
			Volume:  p.Volume(),
		}
	}
	return &orderv1.UserCatalogOut{
		Result: result,
	}
}
