package application

import (
	"context"
	"errors"
	"testing"

	"github.com/thamesblanq/healthconnect/internal/user/domain"
	"github.com/thamesblanq/healthconnect/internal/user/ports"
)

// --------------------------------------------------
// Fake User Repository
// --------------------------------------------------

type fakeUserRepository struct {
	users       map[string]*domain.User
	createError error
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (f *fakeUserRepository) Create(
	ctx context.Context,
	user *domain.User,
) error {
	if f.createError != nil {
		return f.createError
	}

	f.users[user.Email] = user

	return nil
}

func (f *fakeUserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {
	user, exists := f.users[email]

	if !exists {
		return nil, nil
	}

	return user, nil
}

func (f *fakeUserRepository) FindByID(
	ctx context.Context,
	id string,
) (*domain.User, error) {
	for _, user := range f.users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, nil
}

// Make sure the fake implements the repository port.
var _ ports.UserRepository = (*fakeUserRepository)(nil)

// --------------------------------------------------
// Fake Password Hasher
// --------------------------------------------------

type fakePasswordHasher struct {
	hashResult string
	hashError  error
}

func (f *fakePasswordHasher) Hash(password string) (string, error) {
	if f.hashError != nil {
		return "", f.hashError
	}

	return f.hashResult, nil
}

func (f *fakePasswordHasher) Compare(password string, hash string) error {
	return nil
}

// Make sure the fake implements the password hasher port.
var _ ports.PasswordHasher = (*fakePasswordHasher)(nil)

// --------------------------------------------------
// Tests
// --------------------------------------------------

func TestRegisterUser_Success(t *testing.T) {
	repository := newFakeUserRepository()

	hasher := &fakePasswordHasher{
		hashResult: "hashed-password",
	}

	useCase := NewRegisterUserUseCase(
		repository,
		hasher,
	)

	user, err := useCase.Execute(
		context.Background(),
		RegisterUserInput{
			Email:    "TEST@HealthConnect.com ",
			Password: "TestPassword123!",
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.Email != "test@healthconnect.com" {
		t.Fatalf(
			"expected normalized email %q, got %q",
			"test@healthconnect.com",
			user.Email,
		)
	}

	if user.PasswordHash != "hashed-password" {
		t.Fatalf(
			"expected password hash %q, got %q",
			"hashed-password",
			user.PasswordHash,
		)
	}

	if user.Role != "user" {
		t.Fatalf(
			"expected role %q, got %q",
			"user",
			user.Role,
		)
	}

	if !user.IsActive {
		t.Fatal("expected user to be active")
	}

	if _, exists := repository.users["test@healthconnect.com"]; !exists {
		t.Fatal("expected user to be saved in repository")
	}
}

func TestRegisterUser_EmailRequired(t *testing.T) {
	useCase := NewRegisterUserUseCase(
		newFakeUserRepository(),
		&fakePasswordHasher{
			hashResult: "hashed-password",
		},
	)

	_, err := useCase.Execute(
		context.Background(),
		RegisterUserInput{
			Email:    "",
			Password: "TestPassword123!",
		},
	)

	if !errors.Is(err, ErrEmailRequired) {
		t.Fatalf(
			"expected ErrEmailRequired, got %v",
			err,
		)
	}
}

func TestRegisterUser_InvalidEmail(t *testing.T) {
	useCase := NewRegisterUserUseCase(
		newFakeUserRepository(),
		&fakePasswordHasher{
			hashResult: "hashed-password",
		},
	)

	_, err := useCase.Execute(
		context.Background(),
		RegisterUserInput{
			Email:    "not-an-email",
			Password: "TestPassword123!",
		},
	)

	if !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf(
			"expected ErrInvalidEmail, got %v",
			err,
		)
	}
}

func TestRegisterUser_PasswordRequired(t *testing.T) {
	useCase := NewRegisterUserUseCase(
		newFakeUserRepository(),
		&fakePasswordHasher{
			hashResult: "hashed-password",
		},
	)

	_, err := useCase.Execute(
		context.Background(),
		RegisterUserInput{
			Email:    "test@healthconnect.com",
			Password: "",
		},
	)

	if !errors.Is(err, ErrPasswordRequired) {
		t.Fatalf(
			"expected ErrPasswordRequired, got %v",
			err,
		)
	}
}

func TestRegisterUser_PasswordTooShort(t *testing.T) {
	useCase := NewRegisterUserUseCase(
		newFakeUserRepository(),
		&fakePasswordHasher{
			hashResult: "hashed-password",
		},
	)

	_, err := useCase.Execute(
		context.Background(),
		RegisterUserInput{
			Email:    "test@healthconnect.com",
			Password: "1234567",
		},
	)

	if !errors.Is(err, ErrPasswordTooShort) {
		t.Fatalf(
			"expected ErrPasswordTooShort, got %v",
			err,
		)
	}
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	repository := newFakeUserRepository()

	repository.users["test@healthconnect.com"] = &domain.User{
		ID:           "existing-id",
		Email:        "test@healthconnect.com",
		PasswordHash: "existing-hash",
		Role:         "user",
		IsActive:     true,
	}

	useCase := NewRegisterUserUseCase(
		repository,
		&fakePasswordHasher{
			hashResult: "hashed-password",
		},
	)

	_, err := useCase.Execute(
		context.Background(),
		RegisterUserInput{
			Email:    "test@healthconnect.com",
			Password: "TestPassword123!",
		},
	)

	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf(
			"expected ErrUserAlreadyExists, got %v",
			err,
		)
	}
}
