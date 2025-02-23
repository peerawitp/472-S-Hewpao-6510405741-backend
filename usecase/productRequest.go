package usecase

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/minio/minio-go/v7"
)

type ProductRequestUsecase interface {
	CreateProductRequest(productRequest *domain.ProductRequest, files []*multipart.FileHeader, readers []io.Reader) error
	GetDetailByID(id int) (*dto.DetailOfProductRequestResponseDTO, error)
	GetPaginatedProductRequests(page, limit int) (*dto.PaginationGetProductRequestResponse[domain.ProductRequest], error)
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

func (pr *productRequestService) GetDetailByID(id int) (*dto.DetailOfProductRequestResponseDTO, error) {
	productRequest, err := pr.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	res := dto.DetailOfProductRequestResponseDTO{
		ID:        productRequest.ID,
		Desc:      productRequest.Desc,
		Category:  productRequest.Category,
		Images:    productRequest.Images,
		Budget:    productRequest.Budget,
		Quantity:  productRequest.Quantity,
		UserID:    productRequest.UserID,
		User:      productRequest.User,
		Offers:    productRequest.Offers,
		CreatedAt: productRequest.CreatedAt,
		UpdatedAt: productRequest.UpdatedAt,
		DeletedAt: &productRequest.DeletedAt.Time,
	}

	return &res, nil
}

func (pr *productRequestService) GetPaginatedProductRequests(page, limit int) (*dto.PaginationGetProductRequestResponse[domain.ProductRequest], error) {
	productRequests, totalRows, err := pr.repo.FindPaginatedProductRequests(page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := (int(totalRows) + limit - 1) / limit

	res := dto.PaginationGetProductRequestResponse[domain.ProductRequest]{
		Data:       productRequests,
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}

	return &res, nil
}
