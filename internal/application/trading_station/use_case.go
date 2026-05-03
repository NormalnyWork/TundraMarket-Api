package trading_station

import (
	"context"

	tradingstation "tundraMarket/internal/domain/trading_station"
)

type UseCase struct {
	repository tradingstation.TradingStationRepository
}

func NewUseCase(repository tradingstation.TradingStationRepository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

func (uc *UseCase) GetAll(ctx context.Context) ([]*tradingstation.TradingStation, error) {
	return uc.repository.GetAll(ctx)
}
