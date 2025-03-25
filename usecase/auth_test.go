package usecase_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/usecase/test/mock_repos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthUsecase_GetJWT(t *testing.T) {
	mockUserRepo := new(mock_repos.MockUserRepo)
	cfg := &config.Config{JWTSecret: "mysecret"}
	ctx := context.Background()

	authService := usecase.NewAuthUsecase(mockUserRepo, nil, cfg, nil, ctx)

	t.Run("Success", func(t *testing.T) {
		user := &domain.User{
			ID:         "user123",
			Email:      "user@example.com",
			Name:       "John",
			Surname:    "Doe",
			IsVerified: true,
		}

		jwtToken, err := authService.GetJWT(user)

		assert.NoError(t, err)
		assert.NotEmpty(t, jwtToken)

		parsedToken, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		assert.NoError(t, err)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok, "Token claims should be of type jwt.MapClaims")
		assert.Equal(t, user.ID, claims["id"])
		assert.Equal(t, user.Email, claims["email"])
		assert.Equal(t, user.Name, claims["name"])
		assert.Equal(t, user.Surname, claims["surname"])
		assert.Equal(t, true, claims["is_verified"])

		expTime := time.Unix(int64(claims["exp"].(float64)), 0)

		assert.True(t, expTime.After(time.Now()), "JWT expiration time should be in the future")

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error_SigningFailed", func(t *testing.T) {
		user := &domain.User{
			ID: "user123",
		}

		cfg.JWTSecret = ""

		jwtToken, err := authService.GetJWT(user)

		assert.Error(t, err)
		assert.Empty(t, jwtToken)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthUsecase_LoginWithCredentials(t *testing.T) {
	mockUserRepo := new(mock_repos.MockUserRepo)
	cfg := &config.Config{JWTSecret: "mysecret"}
	ctx := context.Background()

	authService := usecase.NewAuthUsecase(mockUserRepo, nil, cfg, nil, ctx)

	t.Run("Success", func(t *testing.T) {
		req := dto.LoginWithCredentialsRequestDTO{
			Email:    "user@example.com",
			Password: "test_password",
		}

		userPassword := "$argon2id$v=19$m=65536,t=3,p=4$FXmLBDy38gqYZbnq93O3dQ$nGnJi8zr0nrWAcx1kXBQA62pBbKvWygShzVCjexi7RI"

		expectedUser := &domain.User{
			ID:         "user123",
			Email:      "user@example.com",
			Password:   &userPassword,
			IsVerified: true,
		}

		mockUserRepo.On("FindByEmail", ctx, req.Email).Return(expectedUser, nil)

		resp, err := authService.LoginWithCredentials(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.AccessToken)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error_UserNotFound", func(t *testing.T) {
		req := dto.LoginWithCredentialsRequestDTO{
			Email:    "user@example.com",
			Password: "password123",
		}

		mockUserRepo.On("FindByEmail", ctx, req.Email).Return(nil, exception.ErrUserNotFound)

		resp, err := authService.LoginWithCredentials(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error_InvalidPassword", func(t *testing.T) {
		req := dto.LoginWithCredentialsRequestDTO{
			Email:    "user@example.com",
			Password: "wrongpassword",
		}

		userPassword := "hashedpassword123"

		expectedUser := &domain.User{
			ID:         "user123",
			Email:      "user@example.com",
			Password:   &userPassword, // Assume hashed password is used
			IsVerified: true,
		}

		mockUserRepo.On("FindByEmail", ctx, req.Email).Return(expectedUser, nil)

		resp, err := authService.LoginWithCredentials(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthUsecase_Register(t *testing.T) {
	mockUserRepo := new(mock_repos.MockUserRepo)
	cfg := &config.Config{}
	ctx := context.Background()

	authService := usecase.NewAuthUsecase(mockUserRepo, nil, cfg, nil, ctx)

	t.Run("Success", func(t *testing.T) {
		req := dto.RegisterUserRequestDTO{
			Email:      "newuser@example.com",
			Password:   "password123",
			Name:       "John Doe",
			MiddleName: nil,
			Surname:    "Test",
		}

		mockUserRepo.On("FindByEmail", ctx, req.Email).Return(nil, exception.ErrUserNotFound)
		mockUserRepo.On("Create", ctx, mock.Anything).Return(nil)

		err := authService.Register(ctx, req)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error_UserExists", func(t *testing.T) {
		req := dto.RegisterUserRequestDTO{
			Email:    "existinguser@example.com",
			Password: "password123",
			Name:     "John Doe",
		}

		mockUserRepo.On("FindByEmail", ctx, req.Email).Return(&domain.User{Email: "existinguser@example.com"}, nil)

		err := authService.Register(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, exception.ErrUserAlreadyExist, err)
		mockUserRepo.AssertExpectations(t)
	})
}
