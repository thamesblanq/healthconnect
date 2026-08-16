package application

import (
	"context"
	"errors"
	"strings"

	"github.com/thamesblanq/healthconnect/internal/auth/ports"
	securityports "github.com/thamesblanq/healthconnect/internal/security/ports"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user account is inactive")
)

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	AccessToken string
	UserID      string
	Email       string
	Role        string
}

type LoginUserUseCase struct {
	userProvider   ports.UserProvider
	passwordHasher securityports.PasswordHasher
	tokenGenerator ports.TokenGenerator
}

func NewLoginUserUseCase(
	userProvider ports.UserProvider,
	passwordHasher securityports.PasswordHasher,
	tokenGenerator ports.TokenGenerator,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		userProvider:   userProvider,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

func (uc *LoginUserUseCase) Execute(
	ctx context.Context,
	input LoginUserInput,
) (*LoginUserOutput, error) {

	// 1. Normalize the email.
	email := strings.ToLower(strings.TrimSpace(input.Email))

	// 2. Validate the input.
	if email == "" || input.Password == "" {
		return nil, ErrInvalidCredentials
	}

	// 3. Find the user.
	user, err := uc.userProvider.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// We deliberately return the same error whether the
	// email doesn't exist or the password is wrong.
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// 4. Check whether the account is active.
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// 5. Compare the supplied password against the stored hash.
	if err := uc.passwordHasher.Compare(
		user.PasswordHash,
		input.Password,
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := uc.tokenGenerator.Generate(
		ctx,
		user.ID,
		user.Email,
		user.Role,
	)
	if err != nil {
		return nil, err
	}

	// 6. Authentication succeeded.
	return &LoginUserOutput{
		AccessToken: token,
		UserID:      user.ID,
		Email:       user.Email,
		Role:        user.Role,
	}, nil
}
