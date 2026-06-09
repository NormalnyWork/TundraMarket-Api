package httptransport

import (
	"errors"
	"log"
	"net/http"
	nomadv1 "tundraMarket/gen/nomad/v1"

	orderv1 "tundraMarket/gen/order/v1"
	apporder "tundraMarket/internal/application/order"
	domainorder "tundraMarket/internal/domain/order"
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
		log.Printf("order create error: %v", err)
		handleOrderError(w, err)
		return
	}

	writeProto(w, http.StatusCreated, apporder.ToCreateProto(order))
}

func (h *OrderHandler) CreateForNomad(w http.ResponseWriter, r *http.Request) {
	actor, ok := orderActorFromRequest(r)
	if !ok {
		writeProtoError(w, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}

	var req orderv1.OrderCreateForNomadIn
	if err := readProto(r, &req); err != nil {
		writeProtoError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY")
		return
	}

	input := apporder.FromCreateForNomadProto(&req)
	input.Actor = actor

	order, err := h.uc.CreateForNomad(r.Context(), input)
	if err != nil {
		log.Printf("order create for nomad error: %v", err)
		handleOrderError(w, err)
		return
	}

	writeProto(w, http.StatusCreated, apporder.ToCreateProto(order))
}

func (h *OrderHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	actor, ok := orderActorFromRequest(r)
	if !ok {
		writeProtoError(w, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}

	var req orderv1.OrderChangeStatusIn
	if err := readProto(r, &req); err != nil {
		writeProtoError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY")
		return
	}

	status, ok := apporder.ProtoStatusToDomain(req.GetNewStatus())
	if !ok {
		writeProtoError(w, http.StatusBadRequest, "UNKNOWN_STATUS")
		return
	}

	changedAt, err := h.uc.ChangeStatus(r.Context(), apporder.ChangeStatusInput{
		Actor:     actor,
		OrderID:   req.GetOrderId(),
		NewStatus: status,
		Comment:   req.Comment,
	})
	if err != nil {
		handleOrderError(w, err)
		return
	}

	writeProto(w, http.StatusOK, apporder.ToChangeStatusProto(changedAt))
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	actor, ok := orderActorFromRequest(r)
	if !ok {
		writeProtoError(w, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}

	var req orderv1.GetOrderListIn
	if err := readProto(r, &req); err != nil {
		writeProtoError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY")
		return
	}

	category, ok := apporder.ProtoCategoryToDomain(req.GetOrderCategory())
	if !ok {
		writeProtoError(w, http.StatusBadRequest, "UNKNOWN_CATEGORY")
		return
	}

	var anchor int32
	if req.Anchor != nil {
		anchor = req.GetAnchor()
	}

	orders, err := h.uc.List(r.Context(), apporder.ListInput{
		Actor:    actor,
		Anchor:   anchor,
		PageSize: req.GetPageSize(),
		Category: category,
	})
	if err != nil {
		handleOrderError(w, err)
		return
	}

	writeProto(w, http.StatusOK, apporder.ToListProto(orders))
}

func (h *OrderHandler) Updates(w http.ResponseWriter, r *http.Request) {
	actor, ok := orderActorFromRequest(r)
	if !ok {
		writeProtoError(w, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}

	var req orderv1.GetOrderUpdatesRequest
	if err := readProtoAllowEmpty(r, &req); err != nil {
		writeProtoError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY")
		return
	}

	orders, err := h.uc.Updates(r.Context(), apporder.UpdatesInput{
		Actor: actor,
		Time:  req.GetTime(),
	})
	if err != nil {
		handleOrderError(w, err)
		return
	}

	writeProto(w, http.StatusOK, apporder.ToUpdatesProto(orders))
}

func orderActorFromRequest(r *http.Request) (apporder.Actor, bool) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		return apporder.Actor{}, false
	}

	return apporder.Actor{
		Role:             claims.Role,
		NomadID:          claims.NomadID,
		TradingStationID: claims.TradingStationID,
	}, true
}

func handleOrderError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domainorder.ErrEmptyCart):
		writeProtoError(w, http.StatusBadRequest, "EMPTY_CART")
	case errors.Is(err, domainorder.ErrInvalidId):
		writeProtoError(w, http.StatusBadRequest, "INVALID_ID")
	case errors.Is(err, domainorder.ErrDistanceTooFar):
		writeProtoError(w, http.StatusBadRequest, "DISTANCE_TOO_FAR")
	case errors.Is(err, domainorder.ErrUnknownStatus):
		writeProtoError(w, http.StatusBadRequest, "UNKNOWN_STATUS")
	case errors.Is(err, domainorder.ErrIllegalStatusChange):
		writeProtoError(w, http.StatusBadRequest, "ILLEGAL_STATUS_CHANGE")
	case errors.Is(err, domainorder.ErrUnknownCategory):
		writeProtoError(w, http.StatusBadRequest, "UNKNOWN_CATEGORY")
	case errors.Is(err, domainorder.ErrForbidden):
		writeProtoError(w, http.StatusForbidden, "FORBIDDEN")
	case errors.Is(err, domainorder.ErrInvalidPhone):
		writeProtoError(w, http.StatusBadRequest, "INVALID_PHONE")
	case errors.Is(err, domainorder.ErrNomadNotFound):
		writeProtoError(w, http.StatusNotFound, "NOMAD_NOT_FOUND")
	default:
		writeProtoError(w, http.StatusInternalServerError, "internal error")
	}
}

func (h *OrderHandler) CurrentOrder(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		writeProtoError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if claims.NomadID == nil {
		writeProtoError(w, http.StatusForbidden, "nomad only")
		return
	}

	order, err := h.uc.GetCurrentOrder(r.Context(), *claims.NomadID)
	if err != nil {
		if errors.Is(err, domainorder.ErrInvalidId) {
			writeProto(w, http.StatusOK, &nomadv1.UserCurrentOrderOut{})
			return
		}
		log.Printf("current order error: %v", err)
		writeProtoError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeProto(w, http.StatusOK, apporder.ToCurrentOrderProto(order))
}

func (h *OrderHandler) CheckStatus(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		writeProtoError(w, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}
	if claims.NomadID == nil {
		writeProtoError(w, http.StatusForbidden, "NOMAD_ONLY")
		return
	}

	var req orderv1.OrderCheckStatusIn
	if err := readProtoAllowEmpty(r, &req); err != nil {
		writeProtoError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY")
		return
	}

	out, err := h.uc.CheckStatus(r.Context(), apporder.CheckStatusInput{
		NomadID: *claims.NomadID,
		After:   req.GetTime(),
	})
	if err != nil {
		handleOrderError(w, err)
		return
	}

	writeProto(w, http.StatusOK, apporder.ToCheckStatusProto(out))
}
