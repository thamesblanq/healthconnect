package application

import (
	"context"
	"net/mail"
	"strings"

	"github.com/thamesblanq/healthconnect/internal/user/domain"
	"github.com/thamesblanq/healthconnect/internal/user/ports"
)

type RegisterUserInput struct {
	Email    string
	Password string
}

type RegisterUserUseCase struct {
	userRepository ports.UserRepository
	passwordHasher ports.PasswordHasher
}

func NewRegisterUserUseCase(
	userRepository ports.UserRepository,
	passwordHasher ports.PasswordHasher,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
	}
}

func (uc *RegisterUserUseCase) Execute(
	ctx context.Context,
	input RegisterUserInput,
) (*domain.User, error) {

	// 1. Normalize the email.
	email := strings.ToLower(strings.TrimSpace(input.Email))

	// 2. Validate the input.
	if email == "" {
		return nil, ErrEmailRequired
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return nil, ErrInvalidEmail
	}

	if input.Password == "" {
		return nil, ErrPasswordRequired
	}

	if len(input.Password) < 8 {
		return nil, ErrPasswordTooShort
	}

	// 3. Check whether the email is already registered.
	existingUser, err := uc.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 4. Hash the password.
	passwordHash, err := uc.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	// 5. Create the domain User.
	user := domain.NewUser(
		email,
		passwordHash,
	)

	// 6. Save the user through the repository port.
	if err := uc.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	// 7. Return the newly created user.
	return user, nil
}
