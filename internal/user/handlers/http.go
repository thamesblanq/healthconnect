package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/thamesblanq/healthconnect/internal/user/application"
)

type Handler struct {
	registerUser *application.RegisterUserUseCase
}

func NewHandler(
	registerUser *application.RegisterUserUseCase,
) *Handler {
	return &Handler{
		registerUser: registerUser,
	}
}

type registerUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerUserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	var request registerUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	user, err := h.registerUser.Execute(
		r.Context(),
		application.RegisterUserInput{
			Email:    request.Email,
			Password: request.Password,
		},
	)

	if err != nil {
		handleRegisterUserError(w, err)
		return
	}

	response := registerUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func handleRegisterUserError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, application.ErrEmailRequired):
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

	case errors.Is(err, application.ErrPasswordRequired):
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

	case errors.Is(err, application.ErrUserAlreadyExists):
		writeError(
			w,
			http.StatusConflict,
			err.Error(),
		)

	case errors.Is(err, application.ErrInvalidEmail):
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

	case errors.Is(err, application.ErrPasswordTooShort):
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

	default:
		writeError(
			w,
			http.StatusInternalServerError,
			"internal server error",
		)
	}
}

func writeError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(errorResponse{
		Error: message,
	})
}
