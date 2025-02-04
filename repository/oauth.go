package repository

import (
	"context"
	"errors"
	"log"

	"github.com/hewpao/hewpao-backend/dto"
)

type OAuthRepository interface {
	VerifyIDToken(ctx context.Context, idToken string) (*dto.OAuthClaims, error)
}

type OAuthRepositoryFactory struct {
	repos map[string]OAuthRepository
}

func NewOAuthRepositoryFactory() OAuthRepositoryFactory {
	return OAuthRepositoryFactory{repos: make(map[string]OAuthRepository)}
}

func (f *OAuthRepositoryFactory) Register(provider string, repo OAuthRepository) {
	log.Println("[🔐] Registered", provider, "OAuth repository!")
	f.repos[provider] = repo
}

func (f *OAuthRepositoryFactory) GetRepository(provider string) (OAuthRepository, error) {
	repo, exists := f.repos[provider]
	if !exists {
		return nil, errors.New("unsupported OAuth provider")
	}
	return repo, nil
}
