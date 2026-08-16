package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/thamesblanq/healthconnect/internal/auth/application"
)

type Handler struct {
	loginUserUseCase *application.LoginUserUseCase
}

func NewHandler(
	loginUserUseCase *application.LoginUserUseCase,
) *Handler {
	return &Handler{
		loginUserUseCase: loginUserUseCase,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request loginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	output, err := h.loginUserUseCase.Execute(
		r.Context(),
		application.LoginUserInput{
			Email:    request.Email,
			Password: request.Password,
		},
	)

	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidCredentials):
			http.Error(
				w,
				"invalid email or password",
				http.StatusUnauthorized,
			)

		case errors.Is(err, application.ErrUserInactive):
			http.Error(
				w,
				"user account is inactive",
				http.StatusForbidden,
			)

		default:
			http.Error(
				w,
				"internal server error",
				http.StatusInternalServerError,
			)
		}

		return
	}

	response := loginResponse{
		AccessToken: output.AccessToken,
		UserID:      output.UserID,
		Email:       output.Email,
		Role:        output.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}
