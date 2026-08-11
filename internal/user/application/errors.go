package application

import "errors"

var (
	ErrEmailRequired     = errors.New("email is required")
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrPasswordRequired  = errors.New("password is required")
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters")
	ErrUserAlreadyExists = errors.New("user with this email already exists")
)
