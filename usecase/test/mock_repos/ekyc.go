package mock_repos

import (
	"mime/multipart"

	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/stretchr/testify/mock"
)

type MockEKYCRepository struct {
	mock.Mock
}

func (m *MockEKYCRepository) Verify(file *multipart.FileHeader) (*dto.EKYCResponseDTO, error) {
	args := m.Called(file)
	return args.Get(0).(*dto.EKYCResponseDTO), args.Error(1)
}

type MockEKYCRepositoryFactory struct {
	mock.Mock
}

func (m *MockEKYCRepositoryFactory) Register(provider string, repo repository.EKYCRepository) {
	m.Called(provider, repo)
}

func (f *MockEKYCRepositoryFactory) GetRepository(provider string) (repository.EKYCRepository, error) {
	args := f.Called(provider)
	if repo, ok := args.Get(0).(repository.EKYCRepository); ok {
		return repo, args.Error(1)
	}
	return nil, args.Error(1)
}
