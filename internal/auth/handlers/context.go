package handlers

import (
	"context"

	"github.com/thamesblanq/healthconnect/internal/auth/ports"
)

func UserClaimsFromContext(
	ctx context.Context,
) (*ports.Identity, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*ports.Identity)

	return claims, ok
}
