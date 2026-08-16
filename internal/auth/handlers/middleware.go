package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/thamesblanq/healthconnect/internal/auth/ports"
)

type contextKey string

const userClaimsKey contextKey = "user_claims"

type AuthMiddleware struct {
	tokenVerifier ports.TokenVerifier
}

func NewAuthMiddleware(
	tokenVerifier ports.TokenVerifier,
) *AuthMiddleware {
	return &AuthMiddleware{
		tokenVerifier: tokenVerifier,
	}
}

func (m *AuthMiddleware) RequireAuth(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Get the Authorization header.
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(
				w,
				"missing authorization header",
				http.StatusUnauthorized,
			)
			return
		}

		// 2. Make sure it uses the Bearer scheme.
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(
				w,
				"invalid authorization header",
				http.StatusUnauthorized,
			)
			return
		}

		tokenString := strings.TrimSpace(parts[1])

		if tokenString == "" {
			http.Error(
				w,
				"missing token",
				http.StatusUnauthorized,
			)
			return
		}

		// 3. Verify the JWT.
		claims, err := m.tokenVerifier.Verify(
			r.Context(),
			tokenString,
		)

		if err != nil {
			http.Error(
				w,
				"invalid or expired token",
				http.StatusUnauthorized,
			)
			return
		}

		// 4. Put the authenticated user's claims
		//    into the request context.
		ctx := context.WithValue(
			r.Context(),
			userClaimsKey,
			claims,
		)

		// 5. Continue to the protected handler.
		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}
