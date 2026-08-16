package ports

import "context"

type UserAuthData struct {
	ID           string
	Email        string
	PasswordHash string
	Role         string
	IsActive     bool
}

type UserProvider interface {
	FindByEmail(ctx context.Context, email string) (*UserAuthData, error)
}
