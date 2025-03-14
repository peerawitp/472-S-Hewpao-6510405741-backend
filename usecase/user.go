package usecase

import (
	"context"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
)

type UserUsecase interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	EditProfile(ctx context.Context, userID string, req dto.EditProfileDTO) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) EditProfile(ctx context.Context, userID string, req dto.EditProfileDTO) error {
	if err := u.userRepo.EditProfile(ctx, userID, req); err != nil {
		return err
	}

	return nil
}
