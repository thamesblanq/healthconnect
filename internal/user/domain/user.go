package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Role         string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(email string, passwordHash string) *User {
	now := time.Now()

	return &User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "user",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
