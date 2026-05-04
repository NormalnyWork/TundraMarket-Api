package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	domainauth "tundraMarket/internal/domain/auth"
)

type TokenIssuer struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenIssuer(secret string, ttl time.Duration) *TokenIssuer {
	return &TokenIssuer{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (i *TokenIssuer) Issue(claims domainauth.TokenClaims) (string, error) {
	now := time.Now()
	payload := tokenPayload{
		Role:             claims.Role,
		Phone:            claims.Phone,
		NomadID:          claims.NomadID,
		TradingStationID: claims.TradingStationID,
		IssuedAt:         now.Unix(),
		ExpiresAt:        now.Add(i.ttl).Unix(),
	}

	if claims.NomadID != nil {
		payload.Subject = fmt.Sprintf("nomad:%d", *claims.NomadID)
	}
	if claims.TradingStationID != nil {
		payload.Subject = fmt.Sprintf("trading_station:%d", *claims.TradingStationID)
	}

	header, err := encodeJSON(tokenHeader{Algorithm: "HS256", Type: "JWT"})
	if err != nil {
		return "", err
	}
	body, err := encodeJSON(payload)
	if err != nil {
		return "", err
	}

	signingInput := strings.Join([]string{header, body}, ".")
	signature := sign(signingInput, i.secret)
	return strings.Join([]string{signingInput, signature}, "."), nil
}

type tokenHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type tokenPayload struct {
	Subject          string `json:"sub"`
	Role             string `json:"role"`
	Phone            string `json:"phone"`
	NomadID          *int32 `json:"nomad_id,omitempty"`
	TradingStationID *int32 `json:"trading_station_id,omitempty"`
	IssuedAt         int64  `json:"iat"`
	ExpiresAt        int64  `json:"exp"`
}

func encodeJSON(value any) (string, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func sign(value string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
