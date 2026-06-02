package apporder

import (
	commonv1 "tundraMarket/gen/common/v1"
	orderv1 "tundraMarket/gen/order/v1"
	domainorder "tundraMarket/internal/domain/order"
)

func FromCreateProto(req *orderv1.OrderCreateIn) CreateInput {
	products := make([]ProductCountInput, len(req.Products))
	for i, p := range req.Products {
		products[i] = ProductCountInput{
			ProductID: p.Id,
			Quantity:  p.Count,
		}
	}

	var longitude, latitude float32
	if len(req.Location) > 0 {
		longitude = req.Location[0].Longitude
		latitude = req.Location[0].Latitude
	}

	return CreateInput{
		TradingStationID: req.GetTradingStationId(),
		Comment:          req.GetComment(),
		Longitude:        longitude,
		Latitude:         latitude,
		Products:         products,
	}
}

func ToCreateProto(o *domainorder.Order) *orderv1.OrderCreateOut {
	return &orderv1.OrderCreateOut{
		OrderId: o.ID(),
	}
}

func ProtoStatusToDomain(status commonv1.Status) (domainorder.Status, bool) {
	switch status {
	case commonv1.Status_STATUS_CREATED:
		return domainorder.StatusCreated, true
	case commonv1.Status_STATUS_PROCESSING:
		return domainorder.StatusProcessing, true
	case commonv1.Status_STATUS_SENT:
		return domainorder.StatusSent, true
	case commonv1.Status_STATUS_COMPLETED:
		return domainorder.StatusCompleted, true
	case commonv1.Status_STATUS_CANCELLED:
		return domainorder.StatusCancelled, true
	case commonv1.Status_STATUS_DENIED:
		return domainorder.StatusDenied, true
	default:
		return "", false
	}
}

func ProtoCategoryToDomain(category commonv1.OrderCategory) (domainorder.OrderCategory, bool) {
	switch category {
	case commonv1.OrderCategory_ORDER_CATEGORY_NEW:
		return domainorder.OrderCategoryNew, true
	case commonv1.OrderCategory_ORDER_CATEGORY_PROCESSING:
		return domainorder.OrderCategoryProcessing, true
	case commonv1.OrderCategory_ORDER_CATEGORY_HISTORY:
		return domainorder.OrderCategoryHistory, true
	default:
		return "", false
	}
}

func ToChangeStatusProto(changedAt int64) *orderv1.OrderChangeStatusOut {
	return &orderv1.OrderChangeStatusOut{
		Time: changedAt,
	}
}

func ToListProto(orders []*domainorder.Order) *orderv1.GetOrderListOut {
	result := make([]*commonv1.Order, len(orders))
	for i, order := range orders {
		result[i] = ToOrderProto(order)
	}

	return &orderv1.GetOrderListOut{
		Orders: result,
	}
}

func ToUpdatesProto(orders []*domainorder.Order) *orderv1.GetOrderUpdatesResponse {
	result := make([]*commonv1.Order, len(orders))
	for i, order := range orders {
		result[i] = ToOrderProto(order)
	}

	return &orderv1.GetOrderUpdatesResponse{
		Orders: result,
	}
}

func ToOrderProto(order *domainorder.Order) *commonv1.Order {
	history := make([]*commonv1.StatusHistory, len(order.History()))
	for i, item := range order.History() {
		history[i] = &commonv1.StatusHistory{
			Status: domainStatusToProto(item.Status()),
			Time:   item.CreatedAt().Unix(),
		}
	}

	products := make([]*commonv1.ProductCount, len(order.Products()))
	for i, product := range order.Products() {
		products[i] = &commonv1.ProductCount{
			Id:    product.ProductID,
			Count: product.Quantity,
		}
	}

	return &commonv1.Order{
		Id:               order.ID(),
		NomadId:          order.NomadID(),
		TradingStationId: order.TradingStationID(),
		Status:           domainStatusToProto(order.Status()),
		History:          history,
		Comment:          order.Comment(),
		Card:             products,
		Location: &commonv1.Location{
			Longitude: order.Longitude(),
			Latitude:  order.Latitude(),
		},
	}
}

func domainStatusToProto(status domainorder.Status) commonv1.Status {
	switch status {
	case domainorder.StatusProcessing:
		return commonv1.Status_STATUS_PROCESSING
	case domainorder.StatusSent:
		return commonv1.Status_STATUS_SENT
	case domainorder.StatusCompleted:
		return commonv1.Status_STATUS_COMPLETED
	case domainorder.StatusCancelled:
		return commonv1.Status_STATUS_CANCELLED
	case domainorder.StatusDenied:
		return commonv1.Status_STATUS_DENIED
	default:
		return commonv1.Status_STATUS_CREATED
	}
}
