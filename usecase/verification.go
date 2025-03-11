package usecase

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/hewpao/hewpao-backend/types"
	"github.com/hewpao/hewpao-backend/util"
)

type VerificationUsecase interface {
	VerifyWithKYC(reader io.Reader, file *multipart.FileHeader, userID string, provider string) error
	GetVerificationInfo(instructorEmail string, info *domain.Verification, verificationID uint) error
	UpdateIsVerified(req *dto.UpdateUserVerificationDTO, userEmail string, instructorEmail string) error
}

type verificationService struct {
	minioRepo        repository.S3Repository
	ctx              context.Context
	cfg              config.Config
	userRepo         repository.UserRepository
	verificationRepo repository.VerificationRepository
	ekycRepoFactory  repository.EKYCRepositoryFactory
}

func NewVerificationService(minioRepo repository.S3Repository, ctx context.Context, cfg config.Config, userRepo repository.UserRepository, verificationRepo repository.VerificationRepository, ekycRepoFactory repository.EKYCRepositoryFactory) VerificationUsecase {
	return &verificationService{
		minioRepo:        minioRepo,
		ctx:              ctx,
		cfg:              cfg,
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		ekycRepoFactory:  ekycRepoFactory,
	}
}

func (v *verificationService) UpdateIsVerified(req *dto.UpdateUserVerificationDTO, userEmail string, instructorEmail string) error {
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
	err = v.userRepo.UpdateVerification(v.ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func newVerification(res *dto.EKYCResponseDTO, userID string, cardImageURI string) *domain.Verification {
	return &domain.Verification{
		UserID:      userID,
		CardImage:   &cardImageURI,
		IDNumber:    res.IDNumber,
		FirstNameTh: res.ThFName,
		LastNameTh:  res.ThLName,
		FirstNameEn: res.EnFName,
		LastNameEn:  res.EnLName,
		Gender:      res.Gender,
		DOBTh:       res.ThDOB,
		DOBEn:       res.EnDOB,
		ExpireTh:    res.ThExpire,
		ExpireEn:    res.EnExpire,
		IssueTh:     res.ThIssue,
		IssueEn:     res.EnIssue,
		Address:     res.Address,
		SubDistrict: res.SubDistrict,
		District:    res.District,
		Province:    res.Province,
		PostalCode:  res.PostalCode,
	}
}

func (v *verificationService) VerifyWithKYC(reader io.Reader, file *multipart.FileHeader, userID string, provider string) error {
	ekyc, err := v.ekycRepoFactory.GetRepository(provider)
	if err != nil {
		return err
	}

	user, err := v.userRepo.FindByID(v.ctx, userID)
	if err != nil {
		return err
	}

	uploadInfo, err := v.minioRepo.UploadFile(v.ctx, file.Filename, reader, file.Size, file.Header.Get("Content-Type"), "verification-images")
	if err != nil {
		return err
	}

	uri := uploadInfo.Bucket + "/" + uploadInfo.Key

	res, err := ekyc.Verify(file)
	if err != nil {
		return err
	}

	verification := newVerification(res, user.ID, uri)
	if err := v.verificationRepo.Create(verification); err != nil {
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

	images := []string{*verification.CardImage}
	urls, err := util.GetUrls(v.minioRepo, v.ctx, &v.cfg, images)
	if err != nil {
		return err
	}

	*info = *verification
	*info.CardImage = urls[0]

	return nil
}
