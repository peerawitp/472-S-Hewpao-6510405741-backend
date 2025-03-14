package gorm

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"gorm.io/gorm"
)

type UserGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) repository.UserRepository {
	return &UserGormRepository{db: db}
}

func (u *UserGormRepository) UpdateVerification(ctx context.Context, user *domain.User) error {
	result := u.db.Model(&domain.User{}).Where("email = ?", user.Email).Select("is_verified").Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserGormRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user := domain.User{}
	if err := u.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserGormRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := domain.User{}
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserGormRepository) Create(ctx context.Context, req dto.CreateUserDTO) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	user := domain.User{
		ID:          id.String(),
		Email:       req.Email,
		Password:    req.Password,
		Name:        req.Name,
		MiddleName:  req.MiddleName,
		Surname:     req.Surname,
		PhoneNumber: req.PhoneNumber,
	}

	createErr := u.db.Create(&user).Error
	if createErr != nil {
		return createErr
	}
	return nil
}

func (u *UserGormRepository) EditProfile(ctx context.Context, userID string, req dto.EditProfileDTO) error {
	user := domain.User{}
	err := u.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return err
	}

	user.Name = req.Name
	user.MiddleName = req.MiddleName
	user.Surname = req.Surname
	user.PhoneNumber = req.PhoneNumber

	err = u.db.Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}
