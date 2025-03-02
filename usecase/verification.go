package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

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
	UpdateIsVerified(req *dto.UpdateUserVerificationDTO, userEmail string, instructorEmail string) error
}

type verificationService struct {
	minioRepo        repository.S3Repository
	ctx              context.Context
	cfg              config.Config
	userRepo         repository.UserRepository
	verificationRepo repository.VerificationRepository
	httpCli          *http.Client
}

func NewVerificationService(minioRepo repository.S3Repository, ctx context.Context, cfg config.Config, userRepo repository.UserRepository, verificationRepo repository.VerificationRepository, httpCli *http.Client) VerificationUsecase {
	return &verificationService{
		minioRepo:        minioRepo,
		ctx:              ctx,
		cfg:              cfg,
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		httpCli:          httpCli,
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

func eKYC(file *multipart.FileHeader, cfg *config.Config, httpCli *http.Client) (*dto.EKYCResponseDTO, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return nil, err
	}

	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	_, err = io.Copy(part, fileReader)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", cfg.KYCApiUrl, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("apikey", cfg.KYCApiKey)

	res, err := httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result dto.EKYCResponseDTO
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
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

func (v *verificationService) VerifyWithKYC(reader io.Reader, file *multipart.FileHeader, userID string) error {
	user, err := v.userRepo.FindByID(v.ctx, userID)
	if err != nil {
		return err
	}

	uploadInfo, err := v.minioRepo.UploadFile(v.ctx, file.Filename, reader, file.Size, file.Header.Get("Content-Type"), "verification-images")
	if err != nil {
		return err
	}

	uri := uploadInfo.Bucket + "/" + uploadInfo.Key

	res, err := eKYC(file, &v.cfg, v.httpCli)
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

	*info = *verification

	return nil
}
