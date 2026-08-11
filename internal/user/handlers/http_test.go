package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thamesblanq/healthconnect/internal/user/application"
	"github.com/thamesblanq/healthconnect/internal/user/domain"
	"github.com/thamesblanq/healthconnect/internal/user/ports"
)

// --------------------------------------------------
// Fake Repository
// --------------------------------------------------

type fakeRepository struct {
	user *domain.User
}

func (f *fakeRepository) Create(
	ctx context.Context,
	user *domain.User,
) error {
	f.user = user
	return nil
}

func (f *fakeRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {
	return nil, nil
}

func (f *fakeRepository) FindByID(
	ctx context.Context,
	id string,
) (*domain.User, error) {
	return nil, nil
}

// Make sure the fake implements the repository port.
var _ ports.UserRepository = (*fakeRepository)(nil)

// --------------------------------------------------
// Fake Password Hasher
// --------------------------------------------------

type fakeHasher struct{}

func (f *fakeHasher) Hash(password string) (string, error) {
	return "hashed-password", nil
}

func (f *fakeHasher) Compare(password string, hash string) error {
	return nil
}

// Make sure the fake implements the password hasher port.
var _ ports.PasswordHasher = (*fakeHasher)(nil)

// --------------------------------------------------
// Tests
// --------------------------------------------------

func TestRegisterUserHandler_Success(t *testing.T) {
	repository := &fakeRepository{}
	hasher := &fakeHasher{}

	useCase := application.NewRegisterUserUseCase(
		repository,
		hasher,
	)

	handler := NewHandler(useCase)

	body := `{
		"email": "test@healthconnect.com",
		"password": "TestPassword123!"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/users",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	handler.RegisterUser(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}

	var response registerUserResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if response.Email != "test@healthconnect.com" {
		t.Fatalf(
			"expected email %q, got %q",
			"test@healthconnect.com",
			response.Email,
		)
	}

	if response.Role != "user" {
		t.Fatalf(
			"expected role %q, got %q",
			"user",
			response.Role,
		)
	}

	if !response.IsActive {
		t.Fatal("expected user to be active")
	}

	if repository.user == nil {
		t.Fatal("expected user to be saved")
	}

	if repository.user.PasswordHash != "hashed-password" {
		t.Fatalf(
			"expected password hash %q, got %q",
			"hashed-password",
			repository.user.PasswordHash,
		)
	}
}

func TestRegisterUserHandler_InvalidEmail(t *testing.T) {
	useCase := application.NewRegisterUserUseCase(
		&fakeRepository{},
		&fakeHasher{},
	)

	handler := NewHandler(useCase)

	body := `{
		"email": "not-an-email",
		"password": "TestPassword123!"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/users",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	handler.RegisterUser(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	var response errorResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if response.Error != "invalid email address" {
		t.Fatalf(
			"expected invalid email error, got %q",
			response.Error,
		)
	}
}

func TestRegisterUserHandler_InvalidJSON(t *testing.T) {
	useCase := application.NewRegisterUserUseCase(
		&fakeRepository{},
		&fakeHasher{},
	)

	handler := NewHandler(useCase)

	request := httptest.NewRequest(
		http.MethodPost,
		"/users",
		strings.NewReader(`{"email":`),
	)

	recorder := httptest.NewRecorder()

	handler.RegisterUser(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}

func TestRegisterUserHandler_MethodNotAllowed(t *testing.T) {
	useCase := application.NewRegisterUserUseCase(
		&fakeRepository{},
		&fakeHasher{},
	)

	handler := NewHandler(useCase)

	request := httptest.NewRequest(
		http.MethodGet,
		"/users",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.RegisterUser(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			recorder.Code,
		)
	}
}
