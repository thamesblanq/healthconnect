package ports

import "context"

type TokenVerifier interface {
	Verify(
		ctx context.Context,
		tokenString string,
	) (*Identity, error)
}
