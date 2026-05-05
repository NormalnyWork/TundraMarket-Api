package trading_station

import (
	commonv1 "tundraMarket/gen/common/v1"
	tradingstationv1 "tundraMarket/gen/trading_station/v1"
	tradingstation "tundraMarket/internal/domain/trading_station"
)

func ToProtoList(trading_station []*tradingstation.TradingStation) *tradingstationv1.TradingStationListOut {
	result := make([]*commonv1.TradingStation, len(trading_station))
	for i, ts := range trading_station {
		result[i] = &commonv1.TradingStation{
			Id:    ts.ID(),
			Name:  ts.Name(),
			Phone: ts.Phone(),
			Location: &commonv1.Location{
				Longitude: float32(ts.Longitude()),
				Latitude:  float32(ts.Latitude()),
			},
		}
	}
	return &tradingstationv1.TradingStationListOut{
		Result: result,
	}
}
