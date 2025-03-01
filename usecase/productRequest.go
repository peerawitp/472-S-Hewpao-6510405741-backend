package usecase

import (
	"context"
	"html/template"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/hewpao/hewpao-backend/types"
	"github.com/minio/minio-go/v7"
)

type ProductRequestUsecase interface {
	CreateProductRequest(productRequest *domain.ProductRequest, files []*multipart.FileHeader, readers []io.Reader) error
	GetDetailByID(id int) (*dto.DetailOfProductRequestResponseDTO, error)
	GetBuyerProductRequestsByUserID(id string) ([]dto.DetailOfProductRequestResponseDTO, error)
	GetPaginatedProductRequests(page, limit int) (*dto.PaginationGetProductRequestResponse[dto.DetailOfProductRequestResponseDTO], error)
	UpdateProductRequest(req *dto.UpdateProductRequestDTO, prID int, userID string) error
	UpdateProductRequestStatus(req *dto.UpdateProductRequestStatusDTO, prID int, userID string) error
}

type productRequestService struct {
	repo                     repository.ProductRequestRepository
	minioRepo                repository.S3Repository
	ctx                      context.Context
	offerRepo                repository.OfferRepository
	userRepo                 repository.UserRepository
	gmailNotificationUsecase NotificationUsecase
	cfg                      *config.Config
}

func NewProductRequestService(repo repository.ProductRequestRepository, minioRepo repository.S3Repository, ctx context.Context, offerRepo repository.OfferRepository, userRepo repository.UserRepository, gmailNotificationUsecase NotificationUsecase, cfg *config.Config) ProductRequestUsecase {
	return &productRequestService{
		repo:                     repo,
		minioRepo:                minioRepo,
		ctx:                      ctx,
		offerRepo:                offerRepo,
		userRepo:                 userRepo,
		gmailNotificationUsecase: gmailNotificationUsecase,
		cfg:                      cfg,
	}
}

func (pr *productRequestService) UpdateProductRequestStatus(req *dto.UpdateProductRequestStatusDTO, prID int, userID string) error {
	productRequest, err := pr.repo.FindByID(prID)
	if err != nil {
		return err
	}

	if productRequest.SelectedOfferID == nil {
		return exception.ErrCouldNotUpdateStatus
	}

	user, err := pr.userRepo.FindByID(pr.ctx, userID)
	if err != nil {
		return err
	}

	offer := new(domain.Offer)
	offer.ID = *productRequest.SelectedOfferID
	err = pr.offerRepo.GetByID(offer)
	if err != nil {
		return err
	}

	if user.Role != types.Admin {
		if offer.UserID != userID {
			return exception.ErrPermissionDenied
		}

		if req.DeliveryStatus != types.Purchased {
			return exception.ErrPermissionDenied
		}
	}

	productRequest.DeliveryStatus = req.DeliveryStatus

	err = pr.repo.Update(productRequest)
	if err != nil {
		return err
	}

	var content strings.Builder

	err = notifyUpdate(productRequest, pr.cfg, &content)
	if err != nil {
		return err
	}

	travelerNotify := dto.NotificationDTO{
		ToID:    offer.UserID,
		Subject: "[HEWPAO] Product Request Current Status Report",
		Content: content.String(),
	}

	err = pr.gmailNotificationUsecase.Notify(&travelerNotify)
	if err != nil {
		return err
	}

	return nil
}

func notifyUpdate(productRequest *domain.ProductRequest, cfg *config.Config, content *strings.Builder) error {
	// Prepare the notification data
	data := dto.NotificationDataDTO{
		RecipientName: productRequest.User.Name,
		CompanyName:   "HEWPAO",
		ProductID:     productRequest.ID,
		ProductStatus: productRequest.DeliveryStatus,
		SupportEmail:  cfg.EmailUser,
		Year:          time.Now().Year(),
	}

	// Parse the template file
	tmpl, err := template.ParseFiles("./assets/emailTemplate.html")
	if err != nil {
		return err
	}

	// Clear any existing content in the builder before executing
	content.Reset()

	// Execute the template and build content
	err = tmpl.Execute(content, data)
	if err != nil {
		return err
	}

	return nil
}

func (pr *productRequestService) UpdateProductRequest(req *dto.UpdateProductRequestDTO, prID int, userID string) error {
	productRequest, err := pr.repo.FindByID(prID)
	if err != nil {
		return err
	}

	if *productRequest.UserID != userID {
		return exception.ErrPermissionDenied
	}

	offer := new(domain.Offer)
	offer.ID = req.SelectedOfferID
	err = pr.offerRepo.GetByID(offer)
	if err != nil {
		return err
	}

	found := false
	for _, o := range productRequest.Offers {
		if o.ID == offer.ID {
			found = true
			break
		}
	}

	if !found {
		return exception.ErrPermissionDenied
	}

	productRequest.Name = req.Name
	productRequest.Desc = req.Desc
	productRequest.Budget = req.Budget
	productRequest.Category = req.Category
	productRequest.Quantity = req.Quantity
	productRequest.SelectedOfferID = &req.SelectedOfferID
	productRequest.SelectedOffer = offer

	err = pr.repo.Update(productRequest)
	if err != nil {
		return err
	}
	return nil
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
		ID:           productRequest.ID,
		Desc:         productRequest.Desc,
		Category:     productRequest.Category,
		Images:       productRequest.Images,
		Budget:       productRequest.Budget,
		Quantity:     productRequest.Quantity,
		UserID:       productRequest.UserID,
		Offers:       productRequest.Offers,
		Transactions: productRequest.Transactions,
		CreatedAt:    productRequest.CreatedAt,
		UpdatedAt:    productRequest.UpdatedAt,
		DeletedAt:    &productRequest.DeletedAt.Time,
	}

	return &res, nil
}

func (pr *productRequestService) GetBuyerProductRequestsByUserID(id string) ([]dto.DetailOfProductRequestResponseDTO, error) {
	productRequests, err := pr.repo.FindByUserID(id)
	if err != nil {
		return nil, err
	}

	res := []dto.DetailOfProductRequestResponseDTO{}

	for _, productRequest := range productRequests {
		productRequestRes := dto.DetailOfProductRequestResponseDTO{
			ID:        productRequest.ID,
			Desc:      productRequest.Desc,
			Category:  productRequest.Category,
			Images:    productRequest.Images,
			Budget:    productRequest.Budget,
			Quantity:  productRequest.Quantity,
			UserID:    productRequest.UserID,
			Offers:    productRequest.Offers,
			CreatedAt: productRequest.CreatedAt,
			UpdatedAt: productRequest.UpdatedAt,
			DeletedAt: &productRequest.DeletedAt.Time,
		}

		res = append(res, productRequestRes)
	}

	return res, nil
}

func (pr *productRequestService) GetPaginatedProductRequests(page, limit int) (*dto.PaginationGetProductRequestResponse[dto.DetailOfProductRequestResponseDTO], error) {
	productRequests, totalRows, err := pr.repo.FindPaginatedProductRequests(page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := (int(totalRows) + limit - 1) / limit

	var dest []dto.DetailOfProductRequestResponseDTO

	for _, productRequest := range productRequests {
		productRequestRes := dto.DetailOfProductRequestResponseDTO{
			ID:        productRequest.ID,
			Desc:      productRequest.Desc,
			Category:  productRequest.Category,
			Images:    productRequest.Images,
			Budget:    productRequest.Budget,
			Quantity:  productRequest.Quantity,
			UserID:    productRequest.UserID,
			Offers:    productRequest.Offers,
			CreatedAt: productRequest.CreatedAt,
			UpdatedAt: productRequest.UpdatedAt,
			DeletedAt: &productRequest.DeletedAt.Time,
		}
		dest = append(dest, productRequestRes)
	}

	res := dto.PaginationGetProductRequestResponse[dto.DetailOfProductRequestResponseDTO]{
		Data:       dest,
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}

	return &res, nil
}
