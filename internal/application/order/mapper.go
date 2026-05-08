package apporder

import (
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
