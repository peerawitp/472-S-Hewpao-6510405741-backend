package repository

import (
	"errors"
	"mime/multipart"

	"github.com/hewpao/hewpao-backend/dto"
)

type EKYCRepository interface {
	Verify(file *multipart.FileHeader) (*dto.EKYCResponseDTO, error)
}

type EKYCRepositoryFactory struct {
	repos map[string]EKYCRepository
}

func NewEKYCRepositoryFactory() EKYCRepositoryFactory {
	return EKYCRepositoryFactory{repos: make(map[string]EKYCRepository)}
}

func (f *EKYCRepositoryFactory) Register(provider string, repo EKYCRepository) {
	f.repos[provider] = repo
}

func (f *EKYCRepositoryFactory) GetRepository(provider string) (EKYCRepository, error) {
	repo, exists := f.repos[provider]
	if !exists {
		return nil, errors.New("unsupported ekyc provider")
	}
	return repo, nil
}
