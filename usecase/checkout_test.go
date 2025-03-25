package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/types"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/usecase/test/mock_repos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCheckoutWithPaymentGateway(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mock_repos.MockUserRepo)
	productRequestRepo := new(mock_repos.MockProductRequestRepo)
	transactionRepo := new(mock_repos.MockTransactionRepository)
	paymentRepo := new(mock_repos.MockPaymentRepository)
	paymentRepoFactory := new(mock_repos.MockPaymentRepositoryFactory)
	minioRepo := new(mock_repos.MockS3Repository)

	checkoutService := usecase.NewCheckoutUsecase(userRepo, productRequestRepo, transactionRepo, paymentRepoFactory, &config.Config{}, minioRepo, ctx)
	t.Run("Success", func(t *testing.T) {
		userRepo.ExpectedCalls = nil
		productRequestRepo.ExpectedCalls = nil
		transactionRepo.ExpectedCalls = nil
		paymentRepo.ExpectedCalls = nil
		paymentRepoFactory.ExpectedCalls = nil
		minioRepo.ExpectedCalls = nil

		// TODO: repo flow

		productRequestID := 1
		userID := "1"
		req := &dto.CheckoutRequestDTO{
			ProductRequestID: uint(productRequestID),
			PaymentGateway:   "provider",
		}

		productRequest := &domain.ProductRequest{
			Model: gorm.Model{},
		}

		payment := &dto.CreatePaymentResponseDTO{
			PaymentID: "payment123",
		}

		productRequestRepo.On("FindByID", productRequestID).Return(productRequest, nil)
		productRequestRepo.On("IsOwnedByUser", productRequestID, userID).Return(true, nil)
		paymentRepoFactory.On("GetRepository", req.PaymentGateway).Return(paymentRepo, nil)
		paymentRepo.On("CreatePayment", ctx, &domain.ProductRequest{}).Return(payment, nil)
		transactionRepo.On("Store", ctx, mock.MatchedBy(func(transaction *domain.Transaction) bool {
			// Check all fields except CreatedAt
			return transaction.UserID == userID &&
				transaction.Amount == productRequest.Budget &&
				transaction.Currency == "THB" &&
				*transaction.ThirdPartyPaymentID == payment.PaymentID &&
				transaction.ThirdPartyGateway == req.PaymentGateway &&
				*transaction.ProductRequestID == productRequest.ID &&
				transaction.Status == types.PaymentPending &&
				time.Since(transaction.CreatedAt) < 2*time.Second
		})).Return(nil)

		_, err := checkoutService.CheckoutWithPaymentGateway(ctx, userID, req)

		assert.NoError(t, err)
		productRequestRepo.AssertExpectations(t)
		paymentRepoFactory.AssertExpectations(t)
		paymentRepo.AssertExpectations(t)
		transactionRepo.AssertExpectations(t)
	})
	t.Run("Error_ProductRequestNotFound", func(t *testing.T) {
		userRepo.ExpectedCalls = nil
		productRequestRepo.ExpectedCalls = nil
		transactionRepo.ExpectedCalls = nil
		paymentRepo.ExpectedCalls = nil
		paymentRepoFactory.ExpectedCalls = nil
		minioRepo.ExpectedCalls = nil

		// TODO: repo flow

		productRequestID := 1
		userID := "1"
		req := &dto.CheckoutRequestDTO{
			ProductRequestID: uint(productRequestID),
			PaymentGateway:   "provider",
		}

		productRequest := &domain.ProductRequest{
			Model: gorm.Model{},
		}

		productRequestRepo.On("FindByID", productRequestID).Return(productRequest, assert.AnError)
		_, err := checkoutService.CheckoutWithPaymentGateway(ctx, userID, req)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		productRequestRepo.AssertExpectations(t)
	})
	t.Run("Error_ProductRequestNoPermission", func(t *testing.T) {
		userRepo.ExpectedCalls = nil
		productRequestRepo.ExpectedCalls = nil
		transactionRepo.ExpectedCalls = nil
		paymentRepo.ExpectedCalls = nil
		paymentRepoFactory.ExpectedCalls = nil
		minioRepo.ExpectedCalls = nil

		// TODO: repo flow

		productRequestID := 1
		userID := "1"
		req := &dto.CheckoutRequestDTO{
			ProductRequestID: uint(productRequestID),
			PaymentGateway:   "provider",
		}

		productRequest := &domain.ProductRequest{
			Model: gorm.Model{},
		}

		productRequestRepo.On("FindByID", productRequestID).Return(productRequest, nil)
		productRequestRepo.On("IsOwnedByUser", productRequestID, userID).Return(false, nil)

		_, err := checkoutService.CheckoutWithPaymentGateway(ctx, userID, req)

		assert.Error(t, err)
		assert.Equal(t, exception.ErrPermissionDenied, err)
		productRequestRepo.AssertExpectations(t)
	})

	t.Run("Error_CouldNotGetProductRequestPermission", func(t *testing.T) {
		userRepo.ExpectedCalls = nil
		productRequestRepo.ExpectedCalls = nil
		transactionRepo.ExpectedCalls = nil
		paymentRepo.ExpectedCalls = nil
		paymentRepoFactory.ExpectedCalls = nil
		minioRepo.ExpectedCalls = nil

		// TODO: repo flow

		productRequestID := 1
		userID := "1"
		req := &dto.CheckoutRequestDTO{
			ProductRequestID: uint(productRequestID),
			PaymentGateway:   "provider",
		}

		productRequest := &domain.ProductRequest{
			Model: gorm.Model{},
		}

		productRequestRepo.On("FindByID", productRequestID).Return(productRequest, nil)
		productRequestRepo.On("IsOwnedByUser", productRequestID, userID).Return(true, assert.AnError)

		_, err := checkoutService.CheckoutWithPaymentGateway(ctx, userID, req)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		productRequestRepo.AssertExpectations(t)
	})
	t.Run("Error_CouldNotGetPaymentRepository", func(t *testing.T) {
		userRepo.ExpectedCalls = nil
		productRequestRepo.ExpectedCalls = nil
		transactionRepo.ExpectedCalls = nil
		paymentRepo.ExpectedCalls = nil
		paymentRepoFactory.ExpectedCalls = nil
		minioRepo.ExpectedCalls = nil

		// TODO: repo flow

		productRequestID := 1
		userID := "1"
		req := &dto.CheckoutRequestDTO{
			ProductRequestID: uint(productRequestID),
			PaymentGateway:   "provider",
		}

		productRequest := &domain.ProductRequest{
			Model: gorm.Model{},
		}

		productRequestRepo.On("FindByID", productRequestID).Return(productRequest, nil)
		productRequestRepo.On("IsOwnedByUser", productRequestID, userID).Return(true, nil)
		paymentRepoFactory.On("GetRepository", req.PaymentGateway).Return(paymentRepo, assert.AnError)

		_, err := checkoutService.CheckoutWithPaymentGateway(ctx, userID, req)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		productRequestRepo.AssertExpectations(t)
		paymentRepoFactory.AssertExpectations(t)
	})
	t.Run("Error_CouldNotCreatePayment", func(t *testing.T) {
		userRepo.ExpectedCalls = nil
		productRequestRepo.ExpectedCalls = nil
		transactionRepo.ExpectedCalls = nil
		paymentRepo.ExpectedCalls = nil
		paymentRepoFactory.ExpectedCalls = nil
		minioRepo.ExpectedCalls = nil

		// TODO: repo flow

		productRequestID := 1
		userID := "1"
		req := &dto.CheckoutRequestDTO{
			ProductRequestID: uint(productRequestID),
			PaymentGateway:   "provider",
		}

		productRequest := &domain.ProductRequest{
			Model: gorm.Model{},
		}

		payment := &dto.CreatePaymentResponseDTO{
			PaymentID: "payment123",
		}

		productRequestRepo.On("FindByID", productRequestID).Return(productRequest, nil)
		productRequestRepo.On("IsOwnedByUser", productRequestID, userID).Return(true, nil)
		paymentRepoFactory.On("GetRepository", req.PaymentGateway).Return(paymentRepo, nil)
		paymentRepo.On("CreatePayment", ctx, &domain.ProductRequest{}).Return(payment, assert.AnError)

		_, err := checkoutService.CheckoutWithPaymentGateway(ctx, userID, req)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		productRequestRepo.AssertExpectations(t)
		paymentRepoFactory.AssertExpectations(t)
		paymentRepo.AssertExpectations(t)
	})
	t.Run("Error_CouldNotStoreTransaction", func(t *testing.T) {
		userRepo.ExpectedCalls = nil
		productRequestRepo.ExpectedCalls = nil
		transactionRepo.ExpectedCalls = nil
		paymentRepo.ExpectedCalls = nil
		paymentRepoFactory.ExpectedCalls = nil
		minioRepo.ExpectedCalls = nil

		// TODO: repo flow

		productRequestID := 1
		userID := "1"
		req := &dto.CheckoutRequestDTO{
			ProductRequestID: uint(productRequestID),
			PaymentGateway:   "provider",
		}

		productRequest := &domain.ProductRequest{
			Model: gorm.Model{},
		}

		payment := &dto.CreatePaymentResponseDTO{
			PaymentID: "payment123",
		}

		productRequestRepo.On("FindByID", productRequestID).Return(productRequest, nil)
		productRequestRepo.On("IsOwnedByUser", productRequestID, userID).Return(true, nil)
		paymentRepoFactory.On("GetRepository", req.PaymentGateway).Return(paymentRepo, nil)
		paymentRepo.On("CreatePayment", ctx, &domain.ProductRequest{}).Return(payment, nil)
		transactionRepo.On("Store", ctx, mock.MatchedBy(func(transaction *domain.Transaction) bool {
			// Check all fields except CreatedAt
			return transaction.UserID == userID &&
				transaction.Amount == productRequest.Budget &&
				transaction.Currency == "THB" &&
				*transaction.ThirdPartyPaymentID == payment.PaymentID &&
				transaction.ThirdPartyGateway == req.PaymentGateway &&
				*transaction.ProductRequestID == productRequest.ID &&
				transaction.Status == types.PaymentPending &&
				time.Since(transaction.CreatedAt) < 2*time.Second
		})).Return(assert.AnError)

		_, err := checkoutService.CheckoutWithPaymentGateway(ctx, userID, req)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		productRequestRepo.AssertExpectations(t)
		paymentRepoFactory.AssertExpectations(t)
		paymentRepo.AssertExpectations(t)
		transactionRepo.AssertExpectations(t)
	})
}

func TestUpdateTransactionStatus(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mock_repos.MockUserRepo)
	productRequestRepo := new(mock_repos.MockProductRequestRepo)
	transactionRepo := new(mock_repos.MockTransactionRepository)
	paymentRepo := new(mock_repos.MockPaymentRepository)
	paymentRepoFactory := new(mock_repos.MockPaymentRepositoryFactory)
	minioRepo := new(mock_repos.MockS3Repository)

	checkoutService := usecase.NewCheckoutUsecase(userRepo, productRequestRepo, transactionRepo, paymentRepoFactory, &config.Config{}, minioRepo, ctx)
	t.Run("Success", func(t *testing.T) {
		userRepo.ExpectedCalls = nil
		productRequestRepo.ExpectedCalls = nil
		transactionRepo.ExpectedCalls = nil
		paymentRepo.ExpectedCalls = nil
		paymentRepoFactory.ExpectedCalls = nil
		minioRepo.ExpectedCalls = nil

		// TODO: repo flow

		transactionRepo.On("UpdateStatusByThirdPartyPaymentID", ctx, "3rdpid", types.PaymentSuccess).Return(nil)

		err := checkoutService.UpdateTransactionStatus(ctx, "3rdpid", types.PaymentSuccess)

		assert.NoError(t, err)
		transactionRepo.AssertExpectations(t)
	})
	t.Run("Error_CouldNotUpdateStatusByThirdParty", func(t *testing.T) {
		userRepo.ExpectedCalls = nil
		productRequestRepo.ExpectedCalls = nil
		transactionRepo.ExpectedCalls = nil
		paymentRepo.ExpectedCalls = nil
		paymentRepoFactory.ExpectedCalls = nil
		minioRepo.ExpectedCalls = nil

		// TODO: repo flow

		transactionRepo.On("UpdateStatusByThirdPartyPaymentID", ctx, "3rdpid", types.PaymentSuccess).Return(assert.AnError)

		err := checkoutService.UpdateTransactionStatus(ctx, "3rdpid", types.PaymentSuccess)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		transactionRepo.AssertExpectations(t)
	})
}
