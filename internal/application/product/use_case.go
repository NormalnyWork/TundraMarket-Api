package product

import (
	"context"

	domainproduct "tundraMarket/internal/domain/product"
)

type UseCase struct {
	repository domainproduct.Repository
}

func NewUseCase(repository domainproduct.Repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

func (uc *UseCase) GetAll(ctx context.Context) ([]*domainproduct.Product, error) {
	return uc.repository.GetAll(ctx)
}
