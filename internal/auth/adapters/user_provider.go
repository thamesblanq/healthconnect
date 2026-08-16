package adapters

import (
	"context"

	"github.com/thamesblanq/healthconnect/internal/auth/ports"
	userports "github.com/thamesblanq/healthconnect/internal/user/ports"
)

type UserProvider struct {
	userRepository userports.UserRepository
}

func NewUserProvider(
	userRepository userports.UserRepository,
) *UserProvider {
	return &UserProvider{
		userRepository: userRepository,
	}
}

func (p *UserProvider) FindByEmail(
	ctx context.Context,
	email string,
) (*ports.UserAuthData, error) {
	user, err := p.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return &ports.UserAuthData{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		IsActive:     user.IsActive,
	}, nil
}

var _ ports.UserProvider = (*UserProvider)(nil)
