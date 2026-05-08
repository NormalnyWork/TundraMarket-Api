package httptransport

import (
	"errors"
	"log"
	"net/http"

	orderv1 "tundraMarket/gen/order/v1"
	apporder "tundraMarket/internal/application/order"
)

type OrderHandler struct {
	uc *apporder.UseCase
}

func NewOrderHandler(uc *apporder.UseCase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		writeProtoError(w, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}
	if claims.NomadID == nil {
		writeProtoError(w, http.StatusForbidden, "NOMAD_ONLY")
		return
	}

	var req orderv1.OrderCreateIn
	if err := readProto(r, &req); err != nil {
		writeProtoError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY")
		return
	}

	input := apporder.FromCreateProto(&req)
	input.NomadID = *claims.NomadID

	order, err := h.uc.Create(r.Context(), input)
	if err != nil {
		log.Printf("order create error: %v", err) // ← добавь
		switch {
		case errors.Is(err, apporder.ErrEmptyCart):
			writeProtoError(w, http.StatusBadRequest, "CART_EMPTY")
		default:
			writeProtoError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	writeProto(w, http.StatusCreated, apporder.ToCreateProto(order))
}
