package auth

import (
	"context"
	"errors"
	"strings"

	"tundraMarket/internal/domain/nomad"
	tradingstation "tundraMarket/internal/domain/trading_station"
)

const (
	RoleNomad          = "nomad"
	RoleTradingStation = "trading_station"
)

var (
	ErrInvalidInput = errors.New("invalid auth input")
	ErrUnauthorized = errors.New("unauthorized")
)

type Input struct {
	Phone            string
	TradingStationID *int32
}

type TokenClaims struct {
	Role             string
	Phone            string
	NomadID          *int32
	TradingStationID *int32
}

type TokenIssuer interface {
	Issue(claims TokenClaims) (string, error)
}

type UseCase struct {
	nomads   nomad.Repository
	stations tradingstation.TradingStationRepository
	tokens   TokenIssuer
}

func NewUseCase(
	nomads nomad.Repository,
	stations tradingstation.TradingStationRepository,
	tokens TokenIssuer,
) *UseCase {
	return &UseCase{
		nomads:   nomads,
		stations: stations,
		tokens:   tokens,
	}
}

func (uc *UseCase) Auth(ctx context.Context, input Input) (string, error) {
	phone := strings.TrimSpace(input.Phone)
	if phone == "" {
		return "", ErrInvalidInput
	}

	if input.TradingStationID != nil {
		return uc.authTradingStation(ctx, phone, *input.TradingStationID)
	}

	return uc.authNomad(ctx, phone)
}

func (uc *UseCase) authNomad(ctx context.Context, phone string) (string, error) {
	n, err := uc.nomads.GetByPhone(ctx, phone)
	if err != nil {
		if !errors.Is(err, nomad.ErrNotFound) {
			return "", err
		}
		n, err = uc.nomads.Create(ctx, phone)
		if err != nil {
			return "", err
		}
	}

	nomadID := n.ID()
	return uc.tokens.Issue(TokenClaims{
		Role:    RoleNomad,
		Phone:   phone,
		NomadID: &nomadID,
	})
}

func (uc *UseCase) authTradingStation(ctx context.Context, phone string, stationID int32) (string, error) {
	station, err := uc.stations.GetByID(ctx, stationID)
	if err != nil {
		if errors.Is(err, tradingstation.ErrNotFound) {
			return "", ErrUnauthorized
		}
		return "", err
	}

	stationPhone := station.Phone()
	if stationPhone == nil || strings.TrimSpace(*stationPhone) == "" {
		station, err = uc.stations.SetPhone(ctx, stationID, phone)
		if err != nil {
			return "", err
		}
		stationPhone = station.Phone()
	}
	if stationPhone == nil || strings.TrimSpace(*stationPhone) != phone {
		return "", ErrUnauthorized
	}

	return uc.tokens.Issue(TokenClaims{
		Role:             RoleTradingStation,
		Phone:            phone,
		TradingStationID: &stationID,
	})
}
