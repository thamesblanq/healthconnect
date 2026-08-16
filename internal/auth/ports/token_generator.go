package ports

import "context"

type TokenGenerator interface {
	Generate(
		ctx context.Context,
		userID string,
		email string,
		role string,
	) (string, error)
}
