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
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func TestProductRequestService_GetDetailByID(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	mockS3 := new(mock_repos.MockS3Repository)
	ctx := context.Background()
	cfg := &config.Config{S3Expiration: "2h30m"}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, mockS3, ctx, nil, nil, nil, cfg, message)

	productRequest := &domain.ProductRequest{Model: gorm.Model{ID: 1}, Name: "Test Product", Images: []string{"hewpao-s3/test_image.png"}}

	mockRepo.On("FindByID", 1).Return(productRequest, nil)
	mockS3.On("GetSignedURL", ctx, mock.Anything, mock.Anything, mock.Anything).Return("https://example.com/hewpao-s3/test_image.png", nil)

	result, err := service.GetDetailByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, int(result.ID))

	arrImage := []string{result.Images[0]}

	assert.Equal(t, []string{"https://example.com/hewpao-s3/test_image.png"}, arrImage)
	mockRepo.AssertExpectations(t)
	mockS3.AssertExpectations(t)
}

func TestProductRequestService_GetBuyerProductRequestsByUserID(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	mockS3 := new(mock_repos.MockS3Repository)
	ctx := context.Background()
	cfg := &config.Config{S3BucketName: "hewpao-s3", S3Expiration: "2h30m"}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, mockS3, ctx, nil, nil, nil, cfg, message)

	productRequests := []domain.ProductRequest{
		{Name: "Product 1", Images: []string{"hewpao-s3/key1"}},
		{Name: "Product 2", Images: []string{"hewpao-s3/key2"}},
	}

	mockRepo.On("FindByUserID", "user1").Return(productRequests, nil)
	mockS3.On("GetSignedURL", ctx, "hewpao-s3", "key1", mock.Anything).Return("https://example.com/hewpao-s3/key1", nil)
	mockS3.On("GetSignedURL", ctx, "hewpao-s3", "key2", mock.Anything).Return("https://example.com/hewpao-s3/key2", nil)

	result, err := service.GetBuyerProductRequestsByUserID("user1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "https://example.com/hewpao-s3/key1", result[0].Images[0])
	assert.Equal(t, "https://example.com/hewpao-s3/key2", result[1].Images[0])
	mockRepo.AssertExpectations(t)
	mockS3.AssertExpectations(t)
}

func TestProductRequestService_GetTravelerProductRequestsByUserID(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	mockS3 := new(mock_repos.MockS3Repository)
	ctx := context.Background()
	cfg := &config.Config{S3BucketName: "hewpao-s3", S3Expiration: "2h30m"}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, mockS3, ctx, nil, nil, nil, cfg, message)

	productRequests := []domain.ProductRequest{
		{Model: gorm.Model{ID: 1}, Name: "Product 1", Images: []string{"hewpao-s3/key1"}},
		{Model: gorm.Model{ID: 2}, Name: "Product 2", Images: []string{"hewpao-s3/key2"}},
	}

	mockRepo.On("FindByOfferUserID", "user1").Return(productRequests, nil)
	mockS3.On("GetSignedURL", ctx, "hewpao-s3", "key1", mock.Anything).Return("url1", nil)
	mockS3.On("GetSignedURL", ctx, "hewpao-s3", "key2", mock.Anything).Return("url2", nil)

	result, err := service.GetTravelerProductRequestsByUserID("user1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "url1", result[0].Images[0])
	assert.Equal(t, "url2", result[1].Images[0])
	mockRepo.AssertExpectations(t)
	mockS3.AssertExpectations(t)
}

func TestProductRequestService_GetPaginatedProductRequests(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	mockS3 := new(mock_repos.MockS3Repository)
	ctx := context.Background()
	cfg := &config.Config{S3BucketName: "hewpao-s3", S3Expiration: "2h30m"}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, mockS3, ctx, nil, nil, nil, cfg, message)

	productRequests := []domain.ProductRequest{
		{Model: gorm.Model{ID: 1}, Name: "Product 1", Images: []string{"hewpao-s3/key1"}},
		{Model: gorm.Model{ID: 2}, Name: "Product 2", Images: []string{"hewpao-s3/key2"}},
	}
	totalRows := int64(2)
	page := 1
	limit := 1

	mockRepo.On("FindPaginatedProductRequests", page, limit).Return(productRequests, totalRows, nil)
	mockS3.On("GetSignedURL", ctx, "hewpao-s3", "key1", mock.Anything).Return("url1", nil)
	mockS3.On("GetSignedURL", ctx, "hewpao-s3", "key2", mock.Anything).Return("url2", nil)

	result, err := service.GetPaginatedProductRequests(page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, page, result.Page)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, totalRows, result.TotalRows)
	assert.Equal(t, 2, result.TotalPages)
	assert.Equal(t, "url1", result.Data[0].Images[0])
	assert.Equal(t, "url2", result.Data[1].Images[0])
	mockRepo.AssertExpectations(t)
	mockS3.AssertExpectations(t)
}

func TestProductRequestService_UpdateProductRequestStatus(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	mockUserRepo := new(mock_repos.MockUserRepo)
	mockOfferRepo := new(mock_repos.MockOfferRepo)
	ctx := context.Background()
	cfg := &config.Config{}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, nil, ctx, mockOfferRepo, mockUserRepo, nil, cfg, message)

	req := &dto.UpdateProductRequestStatusDTO{DeliveryStatus: types.Purchased}
	prID := 1
	userID := "user1"
	offerID := uint(1)

	productRequest := &domain.ProductRequest{
		Model:           gorm.Model{ID: uint(prID)},
		UserID:          &userID,
		DeliveryStatus:  types.Pending,
		SelectedOfferID: &offerID,
	}
	user := &domain.User{
		ID:         userID,
		IsVerified: true,
	}
	offer := &domain.Offer{Model: gorm.Model{ID: 1}}

	mockRepo.On("FindByID", prID).Return(productRequest, nil)
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockOfferRepo.On("GetByID", offer).Return(nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	result, err := service.UpdateProductRequestStatus(req, prID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, types.Purchased, result.DeliveryStatus)
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockOfferRepo.AssertExpectations(t)
}

func TestProductRequestService_UpdateProductRequestStatus_Admin(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	mockUserRepo := new(mock_repos.MockUserRepo)
	mockOfferRepo := new(mock_repos.MockOfferRepo)
	ctx := context.Background()
	cfg := &config.Config{}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, nil, ctx, mockOfferRepo, mockUserRepo, nil, cfg, message)

	req := &dto.UpdateProductRequestStatusDTO{DeliveryStatus: types.Purchased}
	prID := 1
	userID := "admin1"
	offerID := uint(1)

	productRequest := &domain.ProductRequest{
		Model:           gorm.Model{ID: uint(prID)},
		UserID:          &userID,
		DeliveryStatus:  types.Pending,
		SelectedOfferID: &offerID,
	}

	user := &domain.User{
		ID:         userID,
		IsVerified: true,
		Role:       types.Admin,
	}

	mockRepo.On("FindByID", prID).Return(productRequest, nil)
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	result, err := service.UpdateProductRequestStatus(req, prID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, types.Purchased, result.DeliveryStatus)
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockOfferRepo.AssertNotCalled(t, "GetByID", mock.Anything)
}

func TestProductRequestService_UpdateProductRequestStatus_BuyerCancel(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	mockUserRepo := new(mock_repos.MockUserRepo)
	mockOfferRepo := new(mock_repos.MockOfferRepo)
	ctx := context.Background()
	cfg := &config.Config{}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, nil, ctx, mockOfferRepo, mockUserRepo, nil, cfg, message)

	req := &dto.UpdateProductRequestStatusDTO{DeliveryStatus: types.Cancel}
	prID := 1
	userID := "buyer1"

	productRequest := &domain.ProductRequest{
		Model:          gorm.Model{ID: uint(prID)},
		UserID:         &userID,
		DeliveryStatus: types.Pending,
	}

	user := &domain.User{
		ID:         userID,
		IsVerified: false,
	}

	mockRepo.On("FindByID", prID).Return(productRequest, nil)
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	result, err := service.UpdateProductRequestStatus(req, prID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, types.Cancel, result.DeliveryStatus)
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockOfferRepo.AssertNotCalled(t, "GetByID", mock.Anything)
}

func TestProductRequestService_UpdateProductRequestStatusAfterPaid(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	ctx := context.Background()
	cfg := &config.Config{}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, nil, ctx, nil, nil, nil, cfg, message)

	prID := 1
	productRequest := &domain.ProductRequest{
		Model:          gorm.Model{ID: uint(prID)},
		DeliveryStatus: types.Purchased,
	}

	mockRepo.On("FindByID", prID).Return(productRequest, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	err := service.UpdateProductRequestStatusAfterPaid(prID)

	assert.NoError(t, err)
	assert.Equal(t, types.Pending, productRequest.DeliveryStatus)
	mockRepo.AssertExpectations(t)
}

func TestProductRequestService_UpdateProductRequestStatusAfterPaid_FindByIDError(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	ctx := context.Background()
	cfg := &config.Config{}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, nil, ctx, nil, nil, nil, cfg, message)

	mockRepo.On("FindByID", 1).Return(nil, assert.AnError)

	err := service.UpdateProductRequestStatusAfterPaid(1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductRequestService_UpdateProductRequestStatusAfterPaid_UpdateError(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	ctx := context.Background()
	cfg := &config.Config{}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, nil, ctx, nil, nil, nil, cfg, message)

	prID := 1
	productRequest := &domain.ProductRequest{
		Model:          gorm.Model{ID: uint(prID)},
		DeliveryStatus: types.Purchased,
	}

	mockRepo.On("FindByID", prID).Return(productRequest, nil)
	mockRepo.On("Update", mock.Anything).Return(assert.AnError)

	err := service.UpdateProductRequestStatusAfterPaid(prID)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductRequestService_CreateProductRequest(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	mockS3 := new(mock_repos.MockS3Repository)
	mockChat := new(mock_repos.MockChatRepository)
	ctx := context.Background()
	cfg := &config.Config{}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, mockS3, ctx, nil, nil, mockChat, cfg, message)

	productRequest := &domain.ProductRequest{Name: "Test Product"}
	files := []*multipart.FileHeader{{Filename: "test.jpg", Size: 1024, Header: map[string][]string{"Content-Type": {"image/jpeg"}}}}
	readers := []io.Reader{nil}

	uploadInfo := minio.UploadInfo{Bucket: "test-bucket", Key: "test-key"}
	mockS3.On("UploadFile", ctx, "test.jpg", nil, int64(1024), "image/jpeg", "product-request-images").Return(uploadInfo, nil)
	mockChat.On("Create", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		chat := args.Get(0).(*domain.Chat)
		chat.ID = 1 // Assign a mock ID
	})
	mockRepo.On("Create", mock.Anything).Return(nil)

	err := service.CreateProductRequest(productRequest, files, readers)

	assert.NoError(t, err)
	prImagesArr := []string{productRequest.Images[0]}
	assert.Equal(t, []string{"test-bucket/test-key"}, prImagesArr)
	assert.Equal(t, 1, int(productRequest.ChatID))
	mockRepo.AssertExpectations(t)
	mockS3.AssertExpectations(t)
	mockChat.AssertExpectations(t)
}

func TestProductRequestService_CreateProductRequest_UploadFileError(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	mockS3 := new(mock_repos.MockS3Repository)
	mockChat := new(mock_repos.MockChatRepository)
	ctx := context.Background()
	cfg := &config.Config{}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, mockS3, ctx, nil, nil, mockChat, cfg, message)

	productRequest := &domain.ProductRequest{Name: "Test Product"}
	files := []*multipart.FileHeader{{Filename: "test.jpg", Size: 1024, Header: map[string][]string{"Content-Type": {"image/jpeg"}}}}
	readers := []io.Reader{nil}

	mockS3.On("UploadFile", ctx, "test.jpg", nil, int64(1024), "image/jpeg", "product-request-images").Return(minio.UploadInfo{}, assert.AnError)

	err := service.CreateProductRequest(productRequest, files, readers)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	mockS3.AssertExpectations(t)
	mockChat.AssertNotCalled(t, "Create", mock.Anything)
}

func TestProductRequestService_CreateProductRequest_ChatCreateError(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	mockS3 := new(mock_repos.MockS3Repository)
	mockChat := new(mock_repos.MockChatRepository)
	ctx := context.Background()
	cfg := &config.Config{}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, mockS3, ctx, nil, nil, mockChat, cfg, message)

	productRequest := &domain.ProductRequest{Name: "Test Product"}
	files := []*multipart.FileHeader{{Filename: "test.jpg", Size: 1024, Header: map[string][]string{"Content-Type": {"image/jpeg"}}}}
	readers := []io.Reader{nil}

	uploadInfo := minio.UploadInfo{Bucket: "test-bucket", Key: "test-key"}
	mockS3.On("UploadFile", ctx, "test.jpg", nil, int64(1024), "image/jpeg", "product-request-images").Return(uploadInfo, nil)
	mockChat.On("Create", mock.Anything).Return(assert.AnError)

	err := service.CreateProductRequest(productRequest, files, readers)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	mockS3.AssertExpectations(t)
	mockChat.AssertExpectations(t)
}

func TestProductRequestService_CreateProductRequest_RepoCreateError(t *testing.T) {
	mockRepo := new(mock_repos.MockProductRequestRepo)
	mockS3 := new(mock_repos.MockS3Repository)
	mockChat := new(mock_repos.MockChatRepository)
	ctx := context.Background()
	cfg := &config.Config{}
	message := gomail.NewMessage()

	service := usecase.NewProductRequestService(mockRepo, mockS3, ctx, nil, nil, mockChat, cfg, message)

	productRequest := &domain.ProductRequest{Name: "Test Product"}
	files := []*multipart.FileHeader{{Filename: "test.jpg", Size: 1024, Header: map[string][]string{"Content-Type": {"image/jpeg"}}}}
	readers := []io.Reader{nil}

	uploadInfo := minio.UploadInfo{Bucket: "test-bucket", Key: "test-key"}
	mockS3.On("UploadFile", ctx, "test.jpg", nil, int64(1024), "image/jpeg", "product-request-images").Return(uploadInfo, nil)
	mockChat.On("Create", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		chat := args.Get(0).(*domain.Chat)
		chat.ID = 1
	})
	mockRepo.On("Create", mock.Anything).Return(assert.AnError)

	err := service.CreateProductRequest(productRequest, files, readers)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
	mockS3.AssertExpectations(t)
	mockChat.AssertExpectations(t)
}
