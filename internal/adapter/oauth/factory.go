package oauth

import (
	"errors"
	"log"

	"github.com/hewpao/hewpao-backend/domain"
)

type OAuthServiceFactory struct {
	services map[string]domain.OAuthService
}

func NewOAuthServiceFactory() OAuthServiceFactory {
	return OAuthServiceFactory{services: make(map[string]domain.OAuthService)}
}

func (f *OAuthServiceFactory) Register(provider string, service domain.OAuthService) {
	log.Println("[🔐] Registered OAuth service for", provider)
	f.services[provider] = service
}

func (f *OAuthServiceFactory) GetService(provider string) (domain.OAuthService, error) {
	service, exists := f.services[provider]
	if !exists {
		return nil, errors.New("unsupported OAuth provider")
	}
	return service, nil
}
