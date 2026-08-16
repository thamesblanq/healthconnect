package jwt

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/thamesblanq/healthconnect/internal/auth/ports"
)

type TokenVerifier struct {
	secret []byte
}

func NewTokenVerifier(secret string) *TokenVerifier {
	return &TokenVerifier{
		secret: []byte(secret),
	}
}

func (v *TokenVerifier) Verify(
	ctx context.Context,
	tokenString string,
) (*ports.Identity, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			// Make sure the token uses HMAC.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return v.secret, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return nil, errors.New("invalid user ID claim")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return nil, errors.New("invalid email claim")
	}

	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return nil, errors.New("invalid role claim")
	}

	return &ports.Identity{
		UserID: userID,
		Email:  email,
		Role:   role,
	}, nil
}
