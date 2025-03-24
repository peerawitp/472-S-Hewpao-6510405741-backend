package usecase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/hewpao/hewpao-backend/util"
)

type AuthUsecase interface {
	GetJWT(user *domain.User) (string, error)
	LoginWithCredentials(ctx context.Context, req dto.LoginWithCredentialsRequestDTO) (*dto.LoginResponseDTO, error)
	LoginWithOAuth(ctx context.Context, req dto.LoginWithOAuthRequestDTO) (*dto.LoginResponseDTO, error)
	Register(ctx context.Context, req dto.RegisterUserRequestDTO) error
}

type authService struct {
	userRepo         repository.UserRepository
	oauthRepoFactory *repository.OAuthRepositoryFactory
	cfg              *config.Config
	minioRepo        repository.S3Repository
	ctx              context.Context
}

func NewAuthUsecase(userRepo repository.UserRepository, oauthRepoFactory *repository.OAuthRepositoryFactory, cfg *config.Config, minioRepo repository.S3Repository, ctx context.Context) AuthUsecase {
	return &authService{
		userRepo:         userRepo,
		oauthRepoFactory: oauthRepoFactory,
		cfg:              cfg,
		minioRepo:        minioRepo,
		ctx:              ctx,
	}
}

func (a *authService) GetJWT(user *domain.User) (string, error) {
	if a.cfg.JWTSecret == "" {
		return "", exception.ErrJWTSecretIsEmpty
	}

	expiredAt := time.Now().Add(time.Hour * 24) // 1 day
	claims := jwt.MapClaims{
		"id":          user.ID,
		"email":       user.Email,
		"name":        user.Name,
		"middle_name": user.MiddleName,
		"surname":     user.Surname,
		"is_verified": user.IsVerified,
		"exp":         expiredAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *authService) LoginWithCredentials(ctx context.Context, req dto.LoginWithCredentialsRequestDTO) (*dto.LoginResponseDTO, error) {
	user, err := a.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user.Password == nil {
		return nil, exception.ErrUserNoPassword
	}

	same, err := util.VerifyPassword(req.Password, *user.Password)
	if err != nil {
		return nil, err
	}

	if !same {
		return nil, exception.ErrInvalidPassword
	}

	tokenString, err := a.GetJWT(user)
	if err != nil {
		return nil, err
	}

	res := &dto.LoginResponseDTO{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		MiddleName:  user.MiddleName,
		Surname:     user.Surname,
		IsVerified:  user.IsVerified,
		AccessToken: tokenString,
	}

	return res, nil
}

func (a *authService) Register(ctx context.Context, req dto.RegisterUserRequestDTO) error {
	exist, err := a.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && err != exception.ErrUserNotFound {
		return err
	}

	if exist != nil {
		return exception.ErrUserAlreadyExist
	}

	encodedPassword, err := util.HashPassword(req.Password, util.DefaultArgon2Params)
	if err != nil {
		return err
	}

	req.Password = encodedPassword

	createErr := a.userRepo.Create(ctx, dto.CreateUserDTO{
		Email:      req.Email,
		Name:       req.Name,
		MiddleName: req.MiddleName,
		Surname:    req.Surname,
		Password:   &req.Password,
	})
	if createErr != nil {
		return createErr
	}

	return nil
}

func (a *authService) LoginWithOAuth(ctx context.Context, req dto.LoginWithOAuthRequestDTO) (*dto.LoginResponseDTO, error) {
	provider, err := a.oauthRepoFactory.GetRepository(req.Provider)
	if err != nil {
		return nil, err
	}

	userClaims, err := provider.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		return nil, err
	}

	user, err := a.userRepo.FindByEmail(ctx, userClaims.Email)
	if err != nil && err != exception.ErrUserNotFound {
		return nil, err
	}

	// Create the user if not found
	if err == exception.ErrUserNotFound {
		createErr := a.userRepo.Create(ctx, dto.CreateUserDTO{
			Email:      userClaims.Email,
			Name:       userClaims.Name,
			MiddleName: userClaims.MiddleName,
			Surname:    userClaims.Surname,
		})
		if createErr != nil {
			return nil, createErr
		}

		user, err = a.userRepo.FindByEmail(ctx, userClaims.Email)
		if err != nil {
			return nil, err
		}
	}

	tokenString, err := a.GetJWT(user)
	if err != nil {
		return nil, err
	}

	res := &dto.LoginResponseDTO{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		MiddleName:  user.MiddleName,
		Surname:     user.Surname,
		IsVerified:  user.IsVerified,
		AccessToken: tokenString,
	}

	return res, nil
}
