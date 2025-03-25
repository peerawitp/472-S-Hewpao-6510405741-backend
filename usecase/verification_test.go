package usecase_test

import (
	"context"
	"io"
	"mime/multipart"
	"testing"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/types"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/usecase/test/mock_repos"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVerifyWithKYC(t *testing.T) {
	ctx := context.Background()
	mockMinioRepo := new(mock_repos.MockS3Repository)
	mockUserRepo := new(mock_repos.MockUserRepo)
	mockVerificationRepo := new(mock_repos.MockVerificationRepository)
	mockEKYCFactory := new(mock_repos.MockEKYCRepositoryFactory)
	mockEKYCRepo := new(mock_repos.MockEKYCRepository)

	cfg := config.Config{}
	service := usecase.NewVerificationService(mockMinioRepo, ctx, cfg, mockUserRepo, mockVerificationRepo, mockEKYCFactory)

	userID := "12345"
	provider := "mockProvider"
	fileHeader := &multipart.FileHeader{Filename: "test.jpg", Size: 1024}
	reader := io.NopCloser(nil)

	mockUserRepo.On("FindByID", ctx, userID).Return(&domain.User{ID: userID}, nil)
	mockMinioRepo.On("UploadFile", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, "verification-images").Return(minio.UploadInfo{}, nil)
	mockEKYCFactory.On("GetRepository", provider).Return(mockEKYCRepo, nil)
	mockEKYCRepo.On("Verify", fileHeader).Return(&dto.EKYCResponseDTO{IDNumber: "123456789"}, nil)
	mockVerificationRepo.On("Create", mock.Anything).Return(nil)

	err := service.VerifyWithKYC(reader, fileHeader, userID, provider)
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockMinioRepo.AssertExpectations(t)
	mockEKYCFactory.AssertExpectations(t)
	mockEKYCRepo.AssertExpectations(t)
	mockVerificationRepo.AssertExpectations(t)
}

func TestGetVerificationInfo(t *testing.T) {
	ctx := context.Background()
	mockMinioRepo := new(mock_repos.MockS3Repository)
	mockUserRepo := new(mock_repos.MockUserRepo)
	mockVerificationRepo := new(mock_repos.MockVerificationRepository)
	cfg := config.Config{S3Expiration: "2h30m"}
	service := usecase.NewVerificationService(mockMinioRepo, ctx, cfg, mockUserRepo, mockVerificationRepo, &mock_repos.MockEKYCRepositoryFactory{})

	instructorEmail := "admin@example.com"
	verificationID := uint(1)
	verification := &domain.Verification{CardImage: new(string)}
	*verification.CardImage = "hewpao-s3/test_image.png"

	mockUserRepo.On("FindByEmail", ctx, instructorEmail).Return(&domain.User{Role: types.Admin}, nil)
	mockVerificationRepo.On("FindByID", verificationID).Return(verification, nil)
	mockMinioRepo.On("GetSignedURL", ctx, mock.Anything, mock.Anything, mock.Anything).Return("https://example.com/bucket/key", nil)

	var info domain.Verification
	err := service.GetVerificationInfo(instructorEmail, &info, verificationID)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/bucket/key", *info.CardImage)
	mockUserRepo.AssertExpectations(t)
	mockVerificationRepo.AssertExpectations(t)
	mockMinioRepo.AssertExpectations(t)
}

func TestUpdateIsVerified(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(mock_repos.MockUserRepo)
	cfg := config.Config{}
	mockVerificationRepo := new(mock_repos.MockVerificationRepository)
	service := usecase.NewVerificationService(nil, ctx, cfg, mockUserRepo, mockVerificationRepo, &mock_repos.MockEKYCRepositoryFactory{})

	userEmail := "user@example.com"
	instructorEmail := "admin@example.com"
	req := &dto.UpdateUserVerificationDTO{Isverified: true}

	mockUserRepo.On("FindByEmail", ctx, instructorEmail).Return(&domain.User{Role: types.Admin}, nil)
	mockUserRepo.On("FindByEmail", ctx, userEmail).Return(&domain.User{}, nil)
	mockUserRepo.On("UpdateVerification", ctx, mock.Anything).Return(nil)

	err := service.UpdateIsVerified(req, userEmail, instructorEmail)
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}
