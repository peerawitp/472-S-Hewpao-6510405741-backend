package mock_repos

import (
	"mime/multipart"

	"github.com/hewpao/hewpao-backend/dto"
	"github.com/stretchr/testify/mock"
)

type MockEKYCRepository struct {
	mock.Mock
}

func (m *MockEKYCRepository) Verify(file *multipart.FileHeader) (*dto.EKYCResponseDTO, error) {
	args := m.Called(file)
	return args.Get(0).(*dto.EKYCResponseDTO), args.Error(1)
}
