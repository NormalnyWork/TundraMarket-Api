package httptransport

import (
	"net/http"

	appstation "tundraMarket/internal/application/trading_station"
)

type TradingStationHandler struct {
	uc *appstation.UseCase
}

func NewTradingStationHandler(uc *appstation.UseCase) *TradingStationHandler {
	return &TradingStationHandler{uc: uc}
}

func (h *TradingStationHandler) List(w http.ResponseWriter, r *http.Request) {
	stations, error := h.uc.GetAll(r.Context())
	if error != nil {
		writeProtoError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeProto(w, http.StatusOK, appstation.ToProtoList(stations))
}
