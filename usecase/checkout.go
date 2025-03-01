package usecase

import (
	"context"
	"time"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/hewpao/hewpao-backend/types"
)

type CheckoutUsecase interface {
	CheckoutWithPaymentGateway(ctx context.Context, userID string, req *dto.CheckoutRequestDTO) (*dto.CheckoutResponseDTO, error)
	UpdateTransactionStatus(ctx context.Context, transactionID string, status types.PaymentStatus) error
}

type checkoutService struct {
	userRepo           repository.UserRepository
	productRequestRepo repository.ProductRequestRepository
	transactionRepo    repository.TransactionRepository
	paymentRepoFactory *repository.PaymentRepositoryFactory
	cfg                *config.Config
	ctx                context.Context
}

func NewCheckoutUsecase(userRepo repository.UserRepository, productRequestRepo repository.ProductRequestRepository, transactionRepo repository.TransactionRepository, paymentRepoFactory *repository.PaymentRepositoryFactory, cfg *config.Config, minioRepo repository.S3Repository, ctx context.Context) CheckoutUsecase {
	return &checkoutService{
		userRepo:           userRepo,
		productRequestRepo: productRequestRepo,
		transactionRepo:    transactionRepo,
		paymentRepoFactory: paymentRepoFactory,
		cfg:                cfg,
		ctx:                ctx,
	}
}

func (c *checkoutService) CheckoutWithPaymentGateway(ctx context.Context, userID string, req *dto.CheckoutRequestDTO) (*dto.CheckoutResponseDTO, error) {
	_, err := c.productRequestRepo.IsOwnedByUser(int(req.ProductRequestID), userID)
	if err != nil {
		return nil, err
	}

	provider, err := c.paymentRepoFactory.GetRepository(req.PaymentGateway)
	if err != nil {
		return nil, err
	}

	productRequest, err := c.productRequestRepo.FindByID(int(req.ProductRequestID))
	if err != nil {
		return nil, err
	}

	payment, err := provider.CreatePayment(ctx, productRequest)
	if err != nil {
		return nil, err
	}

	// Create transaction
	transaction := &domain.Transaction{
		UserID:              userID,
		Amount:              productRequest.Budget,
		Currency:            "THB",
		Type:                "Income",
		ThirdPartyPaymentID: &payment.PaymentID,
		ThirdPartyGateway:   req.PaymentGateway,
		ProductRequestID:    &productRequest.ID,
		Status:              types.PaymentPending,
		CreatedAt:           time.Now(),
	}

	txErr := c.transactionRepo.Store(ctx, transaction)
	if txErr != nil {
		return nil, txErr
	}

	res := &dto.CheckoutResponseDTO{
		Payment: payment,
	}

	return res, nil
}

func (c *checkoutService) UpdateTransactionStatus(ctx context.Context, thirdPartyPaymentID string, status types.PaymentStatus) error {
	return c.transactionRepo.UpdateStatusByThirdPartyPaymentID(ctx, thirdPartyPaymentID, status)
}
