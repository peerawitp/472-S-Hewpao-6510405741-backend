package domain

import (
	"context"

	"github.com/hewpao/hewpao-backend/dto"
)

type OAuthService interface {
	VerifyIDToken(ctx context.Context, idToken string) (*dto.OAuthClaims, error)
}
