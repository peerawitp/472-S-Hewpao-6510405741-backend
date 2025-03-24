package mock_repos

import (
	"context"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (mu *MockUserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	args := mu.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (mu *MockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := mu.Called(ctx, email)
	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mu *MockUserRepo) Create(ctx context.Context, req dto.CreateUserDTO) error {
	args := mu.Called(ctx, req)
	return args.Error(0)
}

func (mu *MockUserRepo) UpdateVerification(ctx context.Context, req *domain.User) error {
	args := mu.Called(ctx, req)
	return args.Error(0)
}

func (mu *MockUserRepo) EditProfile(ctx context.Context, userID string, req dto.EditProfileDTO) error {
	args := mu.Called(ctx, userID, req)
	return args.Error(0)
}
