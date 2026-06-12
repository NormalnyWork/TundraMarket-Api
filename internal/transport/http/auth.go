package httptransport

import (
	"encoding/json"
	"errors"
	"net/http"

	authv1 "tundraMarket/gen/auth/v1"
	appauth "tundraMarket/internal/application/auth"
)

type AuthHandler struct {
	uc *appauth.UseCase
}

func NewAuthHandler(uc *appauth.UseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var input authv1.UserAuthIn
	if err := readProto(r, &input); err != nil {
		writeProtoError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.uc.Auth(r.Context(), appauth.Input{
		Phone:            input.GetPhone(),
		TradingStationID: input.TradingStationId,
	})
	if err != nil {
		switch {
		case errors.Is(err, appauth.ErrInvalidInput):
			writeProtoError(w, http.StatusBadRequest, "invalid auth input")
		case errors.Is(err, appauth.ErrUnauthorized):
			writeProtoError(w, http.StatusUnauthorized, "unauthorized")
		default:
			writeProtoError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	writeProto(w, http.StatusOK, &authv1.UserAuthOut{Token: token})
}

func (h *AuthHandler) AdminAuth(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeProtoError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY")
		return
	}

	token, err := h.uc.AuthAdmin(r.Context(), req.Login, req.Password)
	if err != nil {
		writeProtoError(w, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
