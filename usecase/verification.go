package usecase

import (
	"context"
	"io"
	"log"
	"mime/multipart"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/hewpao/hewpao-backend/types"
)

type VerificationUsecase interface {
	VerifyWithKYC(reader io.Reader, file *multipart.FileHeader, userID string) error
	GetVerificationInfo(instructorEmail string, info *domain.Verification, verificationID uint) error
	UpdateVerificationInfo(req *dto.UpdateUserVerificationDTO, userEmail string, instructorEmail string) error
}

type verificationService struct {
	minioRepo        repository.S3Repository
	ctx              context.Context
	cfg              config.Config
	userRepo         repository.UserRepository
	verificationRepo repository.VerificationRepository
}

func NewVerificationService(minioRepo repository.S3Repository, ctx context.Context, cfg config.Config, userRepo repository.UserRepository, verificationRepo repository.VerificationRepository) VerificationUsecase {
	return &verificationService{
		minioRepo:        minioRepo,
		ctx:              ctx,
		cfg:              cfg,
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
	}
}

func (v *verificationService) UpdateVerificationInfo(req *dto.UpdateUserVerificationDTO, userEmail string, instructorEmail string) error {
	instructor, err := v.userRepo.FindByEmail(v.ctx, instructorEmail)
	if err != nil {
		return err
	}

	if instructor.Role != types.Admin {
		return exception.ErrPermissionDenied
	}

	user, err := v.userRepo.FindByEmail(v.ctx, userEmail)
	if err != nil {
		return err
	}

	user.IsVerified = req.Isverified
	log.Println(user.IsVerified)
	err = v.userRepo.UpdateVerification(v.ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (v *verificationService) VerifyWithKYC(reader io.Reader, file *multipart.FileHeader, userID string) error {
	user, err := v.userRepo.FindByID(v.ctx, userID)
	if err != nil {
		return err
	}

	verification := domain.Verification{}
	verification.UserID = user.ID
	// TODO: try ekyc here

	uploadInfo, err := v.minioRepo.UploadFile(v.ctx, file.Filename, reader, file.Size, file.Header.Get("Content-Type"), "verification-images")
	if err != nil {
		return err
	}

	uri := uploadInfo.Bucket + "/" + uploadInfo.Key
	verification.CardImage = &uri

	err = v.verificationRepo.Create(&verification)
	if err != nil {
		return err
	}

	return nil
}

func (v *verificationService) GetVerificationInfo(instructorEmail string, info *domain.Verification, verificationID uint) error {
	instructor, err := v.userRepo.FindByEmail(v.ctx, instructorEmail)
	if err != nil {
		return err
	}

	if instructor.Role != types.Admin {
		return exception.ErrPermissionDenied
	}

	verification, err := v.verificationRepo.FindByID(verificationID)
	if err != nil {
		return err
	}

	*info = *verification

	return nil
}
