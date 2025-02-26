package usecase

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/hewpao/hewpao-backend/types"
)

type VerificationUsecase interface {
	VerifyWithKYC(reader io.Reader, file *multipart.FileHeader, email string) error
	GetVerificationInfo(userEmail string, instructorEmail string, info *dto.GetUserVerificationDTO) error
	UpdateVerificationInfo(req *dto.UpdateUserVerificationDTO, userEmail string, instructorEmail string) error
}

type verificationService struct {
	minioRepo repository.S3Repository
	ctx       context.Context
	cfg       config.Config
	userRepo  repository.UserRepository
}

func NewVerificationService(minioRepo repository.S3Repository, ctx context.Context, cfg config.Config, userRepo repository.UserRepository) VerificationUsecase {
	return &verificationService{
		minioRepo: minioRepo,
		ctx:       ctx,
		cfg:       cfg,
		userRepo:  userRepo,
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

	err = v.userRepo.Update(v.ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (v *verificationService) VerifyWithKYC(reader io.Reader, file *multipart.FileHeader, email string) error {
	user, err := v.userRepo.FindByEmail(v.ctx, email)
	if err != nil {
		return err
	}

	uploadInfo, err := v.minioRepo.UploadFile(v.ctx, file.Filename, reader, file.Size, file.Header.Get("Content-Type"), "verification-images")
	if err != nil {
		return err
	}

	uri := uploadInfo.Bucket + "/" + uploadInfo.Key
	user.CardImage = &uri
	v.userRepo.Update(v.ctx, user)

	return nil
}

func (v *verificationService) GetVerificationInfo(userEmail string, instructorEmail string, info *dto.GetUserVerificationDTO) error {
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

	info.Email = user.Email
	info.Name = user.Name
	info.MiddleName = user.MiddleName
	info.Surname = user.Surname
	info.PhoneNumber = user.PhoneNumber
	info.Role = user.Role
	info.IsVerified = user.IsVerified
	info.CardImage = user.CardImage

	return nil
}
