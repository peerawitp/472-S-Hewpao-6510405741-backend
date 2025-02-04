package oauth

import (
	"context"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"google.golang.org/api/idtoken"
)

type GoogleOAuthRepository struct {
	cfg *config.Config
}

func NewGoogleOAuthRepository(cfg *config.Config) repository.OAuthRepository {
	return &GoogleOAuthRepository{
		cfg: cfg,
	}
}

func (r *GoogleOAuthRepository) VerifyIDToken(ctx context.Context, idToken string) (*dto.OAuthClaims, error) {
	payload, err := idtoken.Validate(ctx, idToken, r.cfg.GoogleClientID)
	if err != nil {
		return nil, err
	}

	var surname string
	if payload.Claims["family_name"] == nil {
		surname = ""
	}

	claims := &dto.OAuthClaims{
		Name:    payload.Claims["given_name"].(string),
		Surname: surname,
		Email:   payload.Claims["email"].(string),
		Picture: payload.Claims["picture"].(string),
	}

	return claims, nil
}
