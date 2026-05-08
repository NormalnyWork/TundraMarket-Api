package apporder

import (
	"context"
	"errors"

	domainorder "tundraMarket/internal/domain/order"
)

var ErrEmptyCart = errors.New("cart is empty")

type CreateInput struct {
	NomadID          int32
	TradingStationID int32
	Comment          string
	Products         []ProductCountInput
	Longitude        float32
	Latitude         float32
}

type ProductCountInput struct {
	ProductID int32
	Quantity  int32
}

type UseCase struct {
	repo domainorder.OrderRepository
}

func NewUseCase(repo domainorder.OrderRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Create(ctx context.Context, in CreateInput) (*domainorder.Order, error) {
	if len(in.Products) == 0 {
		return nil, ErrEmptyCart
	}

	items := make([]domainorder.ProductCount, len(in.Products))
	for i, p := range in.Products {
		items[i] = domainorder.ProductCount{
			ProductID: p.ProductID,
			Quantity:  p.Quantity,
		}
	}

	o, err := domainorder.New(in.NomadID, in.TradingStationID, in.Comment, in.Longitude, in.Latitude, items)
	if err != nil {
		return nil, err
	}

	return uc.repo.Save(ctx, o)
}
