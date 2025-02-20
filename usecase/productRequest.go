package usecase

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/minio/minio-go/v7"
)

type ProductRequestUsecase interface {
	CreateProductRequest(productRequest *domain.ProductRequest, files []*multipart.FileHeader, readers []io.Reader) error
}

type productRequestService struct {
	repo      repository.ProductRequestRepository
	minioRepo repository.S3Repository
	ctx       context.Context
}

func NewProductRequestService(repo repository.ProductRequestRepository, minioRepo repository.S3Repository, ctx context.Context) ProductRequestUsecase {
	return &productRequestService{
		repo:      repo,
		minioRepo: minioRepo,
		ctx:       ctx,
	}
}

func (pr *productRequestService) CreateProductRequest(productRequest *domain.ProductRequest, files []*multipart.FileHeader, readers []io.Reader) error {
	uploadInfos := []minio.UploadInfo{}
	for i, file := range files {
		reader := readers[i]

		uploadInfo, err := pr.minioRepo.UploadFile(pr.ctx, file.Filename, reader, file.Size, file.Header.Get("Content-Type"), "product-request-images")
		if err != nil {
			return err
		}

		uploadInfos = append(uploadInfos, uploadInfo)
	}

	uris := []string{}

	for _, uploadInfo := range uploadInfos {
		uri := uploadInfo.Bucket + "/" + uploadInfo.Key
		uris = append(uris, uri)
	}

	productRequest.Images = uris

	err := pr.repo.Create(productRequest)
	if err != nil {
		return err
	}
	return nil
}
