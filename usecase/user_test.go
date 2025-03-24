package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/types"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/usecase/test/mock_repos"

	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	mockUserRepo := new(mock_repos.MockUserRepo)
	userUsecase := usecase.NewUserUsecase(mockUserRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		
		userID := "user123"
		now := time.Now()
		middleName := "Middle"
		phoneNumber := "1234567890"
		password := "hashedpassword"
		
		expectedUser := &domain.User{
			ID:          userID,
			Name:        "Test",
			MiddleName:  &middleName,
			Surname:     "User",
			Email:       "test@example.com",
			Password:    &password,
			PhoneNumber: &phoneNumber,
			CreatedAt:   &now,
			UpdatedAt:   &now,
			Role:        types.Role("User"),
			IsVerified:  true,
		}

		mockUserRepo.On("FindByID", ctx, userID).Return(expectedUser, nil).Once()

		user, err := userUsecase.GetUserByID(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser, user)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, "Test", user.Name)
		assert.Equal(t, "User", user.Surname)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, middleName, *user.MiddleName)
		assert.Equal(t, phoneNumber, *user.PhoneNumber)
		assert.Equal(t, types.Role("User"), user.Role)
		assert.True(t, user.IsVerified)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error_UserNotFound", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		
		userID := "nonexistent"
		expectedError := errors.New("user not found")
		
		emptyUser := &domain.User{
			ID: "",
		}

		mockUserRepo.On("FindByID", ctx, userID).Return(emptyUser, expectedError).Once()

		_, err := userUsecase.GetUserByID(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		
		mockUserRepo.AssertExpectations(t)
	})
}

func TestEditProfile(t *testing.T) {
	mockUserRepo := new(mock_repos.MockUserRepo)
	userUsecase := usecase.NewUserUsecase(mockUserRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		
		userID := "user123"
		middleName := "Updated Middle"
		phoneNumber := "9876543210"
		
		editProfileDTO := dto.EditProfileDTO{
			Name:        "Updated Name",
			MiddleName:  &middleName,
			Surname:     "Updated Surname",
			PhoneNumber: &phoneNumber,
		}

		mockUserRepo.On("EditProfile", ctx, userID, editProfileDTO).Return(nil).Once()

		err := userUsecase.EditProfile(ctx, userID, editProfileDTO)

		assert.NoError(t, err)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error_UpdateFailed", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		
		userID := "user123"
		middleName := "Updated Middle"
		phoneNumber := "9876543210"
		
		editProfileDTO := dto.EditProfileDTO{
			Name:        "Updated Name",
			MiddleName:  &middleName,
			Surname:     "Updated Surname",
			PhoneNumber: &phoneNumber,
		}
		expectedError := errors.New("database error")

		mockUserRepo.On("EditProfile", ctx, userID, editProfileDTO).Return(expectedError).Once()

		err := userUsecase.EditProfile(ctx, userID, editProfileDTO)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockUserRepo.AssertExpectations(t)
	})
}