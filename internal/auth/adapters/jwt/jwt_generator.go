package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/thamesblanq/healthconnect/internal/auth/ports"
)

type TokenGenerator struct {
	secret     []byte
	expiration time.Duration
}

func NewTokenGenerator(
	secret string,
	expiration time.Duration,
) *TokenGenerator {
	return &TokenGenerator{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

func (g *TokenGenerator) Generate(
	ctx context.Context,
	userID string,
	email string,
	role string,
) (string, error) {

	if err := ctx.Err(); err != nil {
		return "", err
	}

	if len(g.secret) == 0 {
		return "", errors.New("JWT secret is not configured")
	}

	now := time.Now()

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  role,
		"iat":   now.Unix(),
		"exp":   now.Add(g.expiration).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString(g.secret)
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}

	return signedToken, nil
}

// Compile-time check that our adapter implements the port.
var _ ports.TokenGenerator = (*TokenGenerator)(nil)
