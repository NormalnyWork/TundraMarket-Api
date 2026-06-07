package apporder

import (
	"context"
	"errors"
	"math"
	"time"

	domainauth "tundraMarket/internal/domain/auth"
	domainorder "tundraMarket/internal/domain/order"
	domainproduct "tundraMarket/internal/domain/product"
	domainstation "tundraMarket/internal/domain/trading_station"
)

const maxDeliveryDistanceKm = 50

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

type Actor struct {
	Role             string
	NomadID          *int32
	TradingStationID *int32
}

type ChangeStatusInput struct {
	Actor     Actor
	OrderID   int32
	NewStatus domainorder.Status
	Comment   *string
}

type ListInput struct {
	Actor    Actor
	Anchor   int32
	PageSize int32
	Category domainorder.OrderCategory
}

type UpdatesInput struct {
	Actor Actor
	Time  int64
}

type UseCase struct {
	repo               domainorder.OrderRepository
	tradingStationRepo domainstation.TradingStationRepository
	productRepo        domainproduct.Repository
}

func NewUseCase(repo domainorder.OrderRepository, tradingStationRepo domainstation.TradingStationRepository, productRepo domainproduct.Repository) *UseCase {
	return &UseCase{
		repo:               repo,
		tradingStationRepo: tradingStationRepo,
		productRepo:        productRepo,
	}
}

func (uc *UseCase) Create(ctx context.Context, in CreateInput) (*domainorder.Order, error) {
	if len(in.Products) == 0 {
		return nil, domainorder.ErrEmptyCart
	}
	if in.NomadID <= 0 || in.TradingStationID <= 0 {
		return nil, domainorder.ErrInvalidId
	}

	station, err := uc.tradingStationRepo.GetByID(ctx, in.TradingStationID)
	if err != nil {
		if errors.Is(err, domainstation.ErrNotFound) {
			return nil, domainorder.ErrInvalidId
		}
		return nil, err
	}

	items := make([]domainorder.ProductCount, len(in.Products))
	productIDs := make([]int32, 0, len(in.Products))
	seenProductIDs := make(map[int32]struct{}, len(in.Products))
	for i, p := range in.Products {
		if p.ProductID <= 0 || p.Quantity <= 0 {
			return nil, domainorder.ErrInvalidId
		}
		items[i] = domainorder.ProductCount{
			ProductID: p.ProductID,
			Quantity:  p.Quantity,
		}
		if _, ok := seenProductIDs[p.ProductID]; !ok {
			seenProductIDs[p.ProductID] = struct{}{}
			productIDs = append(productIDs, p.ProductID)
		}
	}

	products, err := uc.productRepo.GetByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}
	if len(products) != len(productIDs) {
		return nil, domainorder.ErrInvalidId
	}

	if distanceKm(in.Latitude, in.Longitude, station.Latitude(), station.Longitude()) > maxDeliveryDistanceKm {
		return nil, domainorder.ErrDistanceTooFar
	}

	o, err := domainorder.New(in.NomadID, in.TradingStationID, in.Comment, in.Longitude, in.Latitude, items)
	if err != nil {
		return nil, err
	}

	return uc.repo.Save(ctx, o)
}

func distanceKm(lat1, lon1, lat2, lon2 float32) float64 {
	const earthRadiusKm = 6371

	lat1Rad := degreesToRadians(float64(lat1))
	lon1Rad := degreesToRadians(float64(lon1))
	lat2Rad := degreesToRadians(float64(lat2))
	lon2Rad := degreesToRadians(float64(lon2))

	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	return earthRadiusKm * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func (uc *UseCase) ChangeStatus(ctx context.Context, in ChangeStatusInput) (int64, error) {
	if in.OrderID <= 0 {
		return 0, domainorder.ErrInvalidId
	}
	if !isKnownStatus(in.NewStatus) || in.NewStatus == domainorder.StatusCreated {
		return 0, domainorder.ErrUnknownStatus
	}

	order, err := uc.repo.GetByID(ctx, in.OrderID)
	if err != nil {
		return 0, err
	}

	if !canAccess(in.Actor, order) {
		return 0, domainorder.ErrForbidden
	}
	if !canChangeStatus(in.Actor, order, in.NewStatus) {
		return 0, domainorder.ErrIllegalStatusChange
	}

	if _, err := uc.repo.ChangeStatus(ctx, in.OrderID, in.NewStatus, in.Comment); err != nil {
		return 0, err
	}

	return time.Now().Unix(), nil
}

func (uc *UseCase) List(ctx context.Context, in ListInput) ([]*domainorder.Order, error) {
	if in.Anchor < 0 || in.PageSize <= 0 {
		return nil, domainorder.ErrInvalidId
	}
	if !isKnownCategory(in.Category) {
		return nil, domainorder.ErrUnknownCategory
	}

	switch in.Actor.Role {
	case domainauth.RoleNomad:
		if in.Actor.NomadID == nil {
			return nil, domainorder.ErrForbidden
		}
		return uc.repo.ListByNomad(ctx, *in.Actor.NomadID, in.Category, in.Anchor, in.PageSize)
	case domainauth.RoleTradingStation:
		if in.Actor.TradingStationID == nil {
			return nil, domainorder.ErrForbidden
		}
		return uc.repo.ListByTradingStation(ctx, *in.Actor.TradingStationID, in.Category, in.Anchor, in.PageSize)
	default:
		return nil, domainorder.ErrForbidden
	}
}

func (uc *UseCase) Updates(ctx context.Context, in UpdatesInput) ([]*domainorder.Order, error) {
	switch in.Actor.Role {
	case domainauth.RoleNomad:
		if in.Actor.NomadID == nil {
			return nil, domainorder.ErrForbidden
		}
		return uc.repo.GetUpdatesByNomad(ctx, *in.Actor.NomadID, in.Time)
	case domainauth.RoleTradingStation:
		if in.Actor.TradingStationID == nil {
			return nil, domainorder.ErrForbidden
		}
		return uc.repo.GetUpdatesByTradingStation(ctx, *in.Actor.TradingStationID, in.Time)
	default:
		return nil, domainorder.ErrForbidden
	}
}

func canAccess(actor Actor, order *domainorder.Order) bool {
	switch actor.Role {
	case domainauth.RoleNomad:
		return actor.NomadID != nil && *actor.NomadID == order.NomadID()
	case domainauth.RoleTradingStation:
		return actor.TradingStationID != nil && *actor.TradingStationID == order.TradingStationID()
	default:
		return false
	}
}

func canChangeStatus(actor Actor, order *domainorder.Order, next domainorder.Status) bool {
	switch actor.Role {
	case domainauth.RoleNomad:
		return order.Status() == domainorder.StatusCreated && next == domainorder.StatusCancelled
	case domainauth.RoleTradingStation:
		switch order.Status() {
		case domainorder.StatusCreated:
			return next == domainorder.StatusProcessing || next == domainorder.StatusDenied
		case domainorder.StatusProcessing:
			return next == domainorder.StatusSent
		case domainorder.StatusSent:
			return next == domainorder.StatusCompleted
		default:
			return false
		}
	default:
		return false
	}
}

func isKnownStatus(status domainorder.Status) bool {
	switch status {
	case domainorder.StatusCreated,
		domainorder.StatusProcessing,
		domainorder.StatusSent,
		domainorder.StatusCompleted,
		domainorder.StatusCancelled,
		domainorder.StatusDenied:
		return true
	default:
		return false
	}
}

func isKnownCategory(category domainorder.OrderCategory) bool {
	switch category {
	case domainorder.OrderCategoryNew,
		domainorder.OrderCategoryProcessing,
		domainorder.OrderCategoryHistory:
		return true
	default:
		return false
	}
}
