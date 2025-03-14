package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/hewpao/hewpao-backend/config"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/hewpao/hewpao-backend/types"
	"github.com/hewpao/hewpao-backend/util"
	"github.com/minio/minio-go/v7"
	"gopkg.in/gomail.v2"
)

type ProductRequestUsecase interface {
	CreateProductRequest(productRequest *domain.ProductRequest, files []*multipart.FileHeader, readers []io.Reader) error
	GetDetailByID(id int) (*dto.DetailOfProductRequestResponseDTO, error)
	GetBuyerProductRequestsByUserID(id string) ([]dto.DetailOfProductRequestResponseDTO, error)
	GetTravelerProductRequestsByUserID(id string) ([]dto.DetailOfProductRequestResponseDTO, error)
	GetPaginatedProductRequests(page, limit int) (*dto.PaginationGetProductRequestResponse[dto.DetailOfProductRequestResponseDTO], error)
	UpdateProductRequest(req *dto.UpdateProductRequestDTO, prID int, userID string) error
	UpdateProductRequestStatus(req *dto.UpdateProductRequestStatusDTO, prID int, userID string) (*domain.ProductRequest, error)
	UpdateProductRequestStatusAfterPaid(prID int) error
}

type productRequestService struct {
	repo      repository.ProductRequestRepository
	minioRepo repository.S3Repository
	ctx       context.Context
	offerRepo repository.OfferRepository
	userRepo  repository.UserRepository
	chatRepo  repository.ChatRepository
	cfg       *config.Config
	message   *gomail.Message
}

func NewProductRequestService(repo repository.ProductRequestRepository, minioRepo repository.S3Repository, ctx context.Context, offerRepo repository.OfferRepository, userRepo repository.UserRepository, chatRepo repository.ChatRepository, cfg *config.Config, message *gomail.Message) ProductRequestUsecase {
	return &productRequestService{
		repo:      repo,
		minioRepo: minioRepo,
		ctx:       ctx,
		offerRepo: offerRepo,
		userRepo:  userRepo,
		cfg:       cfg,
		message:   message,
		chatRepo:  chatRepo,
	}
}

func (pr *productRequestService) UpdateProductRequestStatus(req *dto.UpdateProductRequestStatusDTO, prID int, userID string) (*domain.ProductRequest, error) {
	productRequest, err := pr.repo.FindByID(prID)
	if err != nil {
		return nil, err
	}

	user, err := pr.userRepo.FindByID(pr.ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Role != types.Admin {
		switch user.IsVerified {
		case true: // traveler > purchase + cancel
			if productRequest.SelectedOfferID == nil {
				return nil, exception.ErrCouldNotUpdateStatus
			}

			offer := new(domain.Offer)
			offer.ID = *productRequest.SelectedOfferID
			err = pr.offerRepo.GetByID(offer)
			if err != nil {
				return nil, err
			}

			allowedTravelerTransitions := map[types.DeliveryStatus]bool{
				types.Purchased: true,
				types.PickedUp:  true,
				types.Cancel:    true,
			}

			if !allowedTravelerTransitions[req.DeliveryStatus] ||
				!types.AllowedStatusTransitions[productRequest.DeliveryStatus][req.DeliveryStatus] {
				return nil, exception.ErrPermissionDenied
			}

		case false: // Buyer: Allowed to transition only to Cancel
			if *productRequest.UserID != userID {
				return nil, exception.ErrPermissionDenied
			}

			allowedBuyerTransitions := map[types.DeliveryStatus]bool{
				types.Cancel: true,
			}

			if !allowedBuyerTransitions[req.DeliveryStatus] ||
				!types.AllowedStatusTransitions[productRequest.DeliveryStatus][req.DeliveryStatus] {
				return nil, exception.ErrPermissionDenied
			}
		}
	}

	productRequest.DeliveryStatus = req.DeliveryStatus

	err = pr.repo.Update(productRequest)
	if err != nil {
		return nil, err
	}

	return productRequest, nil
}

func (pr *productRequestService) UpdateProductRequest(req *dto.UpdateProductRequestDTO, prID int, userID string) error {
	fmt.Println(prID, userID)
	productRequest, err := pr.repo.FindByID(prID)
	if err != nil {
		return err
	}

	if *productRequest.UserID != userID {
		return exception.ErrPermissionDenied
	}

	if req.SelectedOfferID != 0 {

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
			return exception.ErrOfferNotFound
		}

	}

	productRequest.Name = req.Name
	productRequest.Desc = req.Desc
	productRequest.Category = req.Category
	productRequest.Quantity = req.Quantity
	if req.SelectedOfferID != 0 {
		productRequest.SelectedOfferID = &req.SelectedOfferID
	} else {
		productRequest.SelectedOfferID = nil
	}

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

	chatName := productRequest.Name
	newChat := domain.Chat{
		Name: chatName,
	}

	err := pr.chatRepo.Create(&newChat)
	if err != nil {
		return err
	}

	productRequest.ChatID = newChat.ID

	err = pr.repo.Create(productRequest)
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
	urls, err := util.GetUrls(pr.minioRepo, pr.ctx, pr.cfg, productRequest.Images)
	if err != nil {
		return nil, err
	}

	res := dto.DetailOfProductRequestResponseDTO{
		ID:              productRequest.ID,
		Name:            productRequest.Name,
		Desc:            productRequest.Desc,
		Category:        productRequest.Category,
		Images:          urls,
		Budget:          productRequest.Budget,
		Quantity:        productRequest.Quantity,
		UserID:          productRequest.UserID,
		Offers:          productRequest.Offers,
		To:              productRequest.To,
		From:            productRequest.From,
		CheckService:    productRequest.CheckService,
		Transactions:    productRequest.Transactions,
		SelectedOfferID: productRequest.SelectedOfferID,
		DeliveryStatus:  productRequest.DeliveryStatus,
		CreatedAt:       productRequest.CreatedAt,
		UpdatedAt:       productRequest.UpdatedAt,
		DeletedAt:       &productRequest.DeletedAt.Time,
	}

	return &res, nil
}

func (pr *productRequestService) GetTravelerProductRequestsByUserID(id string) ([]dto.DetailOfProductRequestResponseDTO, error) {
	productRequests, err := pr.repo.FindByOfferUserID(id)
	if err != nil {
		return nil, err
	}

	res := []dto.DetailOfProductRequestResponseDTO{}

	for _, productRequest := range productRequests {
		urls, err := util.GetUrls(pr.minioRepo, pr.ctx, pr.cfg, productRequest.Images)
		if err != nil {
			return nil, err
		}

		productRequestRes := dto.DetailOfProductRequestResponseDTO{
			ID:              productRequest.ID,
			Name:            productRequest.Name,
			Desc:            productRequest.Desc,
			Category:        productRequest.Category,
			Images:          urls,
			Budget:          productRequest.Budget,
			Quantity:        productRequest.Quantity,
			UserID:          productRequest.UserID,
			Offers:          productRequest.Offers,
			To:              productRequest.To,
			From:            productRequest.From,
			CheckService:    productRequest.CheckService,
			SelectedOfferID: productRequest.SelectedOfferID,
			DeliveryStatus:  productRequest.DeliveryStatus,
			CreatedAt:       productRequest.CreatedAt,
			UpdatedAt:       productRequest.UpdatedAt,
			DeletedAt:       &productRequest.DeletedAt.Time,
		}

		res = append(res, productRequestRes)
	}

	return res, nil
}

func (pr *productRequestService) GetBuyerProductRequestsByUserID(id string) ([]dto.DetailOfProductRequestResponseDTO, error) {
	productRequests, err := pr.repo.FindByUserID(id)
	if err != nil {
		return nil, err
	}

	res := []dto.DetailOfProductRequestResponseDTO{}

	for _, productRequest := range productRequests {
		urls, err := util.GetUrls(pr.minioRepo, pr.ctx, pr.cfg, productRequest.Images)
		if err != nil {
			return nil, err
		}

		productRequestRes := dto.DetailOfProductRequestResponseDTO{
			ID:              productRequest.ID,
			Name:            productRequest.Name,
			Desc:            productRequest.Desc,
			Category:        productRequest.Category,
			Images:          urls,
			Budget:          productRequest.Budget,
			Quantity:        productRequest.Quantity,
			UserID:          productRequest.UserID,
			Offers:          productRequest.Offers,
			To:              productRequest.To,
			From:            productRequest.From,
			CheckService:    productRequest.CheckService,
			DeliveryStatus:  productRequest.DeliveryStatus,
			SelectedOfferID: productRequest.SelectedOfferID,
			CreatedAt:       productRequest.CreatedAt,
			UpdatedAt:       productRequest.UpdatedAt,
			DeletedAt:       &productRequest.DeletedAt.Time,
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
		urls, err := util.GetUrls(pr.minioRepo, pr.ctx, pr.cfg, productRequest.Images)
		if err != nil {
			return nil, err
		}
		productRequestRes := dto.DetailOfProductRequestResponseDTO{
			ID:              productRequest.ID,
			Name:            productRequest.Name,
			Desc:            productRequest.Desc,
			Category:        productRequest.Category,
			Images:          urls,
			Budget:          productRequest.Budget,
			Quantity:        productRequest.Quantity,
			UserID:          productRequest.UserID,
			Offers:          productRequest.Offers,
			To:              productRequest.To,
			From:            productRequest.From,
			CheckService:    productRequest.CheckService,
			DeliveryStatus:  productRequest.DeliveryStatus,
			SelectedOfferID: productRequest.SelectedOfferID,
			CreatedAt:       productRequest.CreatedAt,
			UpdatedAt:       productRequest.UpdatedAt,
			DeletedAt:       &productRequest.DeletedAt.Time,
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

func (pr *productRequestService) CancleProductRequest(prID int, userID string) error {
	productRequest, err := pr.repo.FindByID(prID)
	if err != nil {
		return err
	}

	if *productRequest.UserID != userID {
		return errors.New("you are not the owner of this product request")
	}

	err = pr.repo.Delete(productRequest)
	if err != nil {
		return err
	}

	return nil
}

func (pr *productRequestService) UpdateProductRequestStatusAfterPaid(prID int) error {
	productRequest, err := pr.repo.FindByID(prID)
	if err != nil {
		return err
	}

	productRequest.DeliveryStatus = types.Pending
	err = pr.repo.Update(productRequest)
	if err != nil {
		return err
	}

	return nil
}
