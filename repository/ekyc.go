package repository

import (
	"errors"
	"mime/multipart"

	"github.com/hewpao/hewpao-backend/dto"
)

type EKYCRepository interface {
	Verify(file *multipart.FileHeader) (*dto.EKYCResponseDTO, error)
}

type EKYCRepositoryFactory interface {
	Register(provider string, repo EKYCRepository)
	GetRepository(provider string) (EKYCRepository, error)
}

type eKYCRepositoryFactory struct {
	repos map[string]EKYCRepository
}

func NewEKYCRepositoryFactory() EKYCRepositoryFactory {
	return &eKYCRepositoryFactory{repos: make(map[string]EKYCRepository)}
}

func (f *eKYCRepositoryFactory) Register(provider string, repo EKYCRepository) {
	f.repos[provider] = repo
}

func (f *eKYCRepositoryFactory) GetRepository(provider string) (EKYCRepository, error) {
	repo, exists := f.repos[provider]
	if !exists {
		return nil, errors.New("unsupported ekyc provider")
	}
	return repo, nil
}
