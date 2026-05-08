package httptransport

import (
	"context"
	"net/http"
	"strings"

	domainauth "tundraMarket/internal/domain/auth"
)

type contextKey string

const claimsKey contextKey = "claims"

func JWTMiddleware(verifier domainauth.TokenIssuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeProtoError(w, http.StatusUnauthorized, "missing token")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeProtoError(w, http.StatusUnauthorized, "invalid token format")
				return
			}

			claims, err := verifier.Verify(parts[1])
			if err != nil {
				writeProtoError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClaimsFromContext(ctx context.Context) *domainauth.TokenClaims {
	claims, _ := ctx.Value(claimsKey).(*domainauth.TokenClaims)
	return claims
}
