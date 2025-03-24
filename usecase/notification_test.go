package usecase_test

import (
	"context"
	"testing"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/types"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/usecase/test/mock_repos"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func TestPrNotify(t *testing.T) {
	mockNotificationRepoFactory := new(mock_repos.MockNotificationRepositoryFactory)
	mockNotificationRepo := new(mock_repos.MockNotificationRepository)
	mockUserRepo := new(mock_repos.MockUserRepo)
	mockOfferRepo := new(mock_repos.MockOfferRepo)
	mockMessage := gomail.NewMessage()
	mockConfig := &config.Config{}
	ctx := context.Background()

	notificationUsecase := usecase.NewNotificationUsecase(mockNotificationRepoFactory, mockUserRepo, ctx, mockMessage, mockConfig, mockOfferRepo)

	t.Run("Success_NotifyBuyer", func(t *testing.T) {
		mockNotificationRepoFactory.ExpectedCalls = nil
		mockNotificationRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil
		mockOfferRepo.ExpectedCalls = nil

		provider := "email"
		buyerID := "buyer123"

		productRequest := &domain.ProductRequest{
			UserID:         &buyerID,
			DeliveryStatus: types.Pending,
		}

		buyer := &domain.User{
			ID: buyerID,
		}

		mockNotificationRepoFactory.On("GetRepository", provider).Return(mockNotificationRepo, nil)
		mockUserRepo.On("FindByID", ctx, buyerID).Return(buyer, nil)
		mockNotificationRepo.On("PrNotify", buyer, productRequest).Return(nil)

		err := notificationUsecase.PrNotify(productRequest, provider)

		assert.NoError(t, err)
		mockNotificationRepoFactory.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockNotificationRepo.AssertExpectations(t)
	})

	t.Run("Success_NotifyTraveler", func(t *testing.T) {
		mockNotificationRepoFactory.ExpectedCalls = nil
		mockNotificationRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil
		mockOfferRepo.ExpectedCalls = nil

		provider := "email"
		buyerID := "buyer123"
		travelerID := "traveler123"
		selectedOfferID := uint(1)

		productRequest := &domain.ProductRequest{
			UserID:          &buyerID,
			DeliveryStatus:  types.Opening,
			SelectedOfferID: &selectedOfferID,
		}

		traveler := &domain.User{
			ID: travelerID,
		}

		offer := &domain.Offer{
			Model: gorm.Model{ID: 1},
		}

		mockNotificationRepoFactory.On("GetRepository", provider).Return(mockNotificationRepo, nil)
		mockOfferRepo.On("GetByID", offer).Return(nil)
		mockUserRepo.On("FindByID", ctx, travelerID).Return(traveler, nil)
		mockNotificationRepo.On("PrNotify", traveler, productRequest).Return(nil)

		err := notificationUsecase.PrNotify(productRequest, provider)

		assert.NoError(t, err)
		mockNotificationRepoFactory.AssertExpectations(t)
		mockOfferRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockNotificationRepo.AssertExpectations(t)
	})
	t.Run("Success_NotifyBuyerAndTraveler", func(t *testing.T) {
		mockNotificationRepoFactory.ExpectedCalls = nil
		mockNotificationRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil
		mockOfferRepo.ExpectedCalls = nil

		provider := "email"
		buyerID := "buyer123"
		travelerID := "traveler123"
		selectedOfferID := uint(1)

		productRequest := &domain.ProductRequest{
			UserID:          &buyerID,
			DeliveryStatus:  types.PickedUp,
			SelectedOfferID: &selectedOfferID,
		}

		traveler := &domain.User{
			ID: travelerID,
		}

		buyer := &domain.User{
			ID: buyerID,
		}

		offer := &domain.Offer{
			Model: gorm.Model{ID: 1},
		}

		mockNotificationRepoFactory.On("GetRepository", provider).Return(mockNotificationRepo, nil)
		mockOfferRepo.On("GetByID", offer).Return(nil)
		mockUserRepo.On("FindByID", ctx, buyerID).Return(buyer, nil)
		mockNotificationRepo.On("PrNotify", buyer, productRequest).Return(nil)
		mockUserRepo.On("FindByID", ctx, travelerID).Return(traveler, nil)
		mockNotificationRepo.On("PrNotify", traveler, productRequest).Return(nil)

		err := notificationUsecase.PrNotify(productRequest, provider)

		assert.NoError(t, err)
		mockNotificationRepoFactory.AssertExpectations(t)
		mockOfferRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockNotificationRepo.AssertExpectations(t)
	})
	t.Run("Success_NotifyBuyerAndTravelerCaseNoTraveler", func(t *testing.T) {
		mockNotificationRepoFactory.ExpectedCalls = nil
		mockNotificationRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil
		mockOfferRepo.ExpectedCalls = nil

		provider := "email"
		buyerID := "buyer123"

		productRequest := &domain.ProductRequest{
			UserID:         &buyerID,
			DeliveryStatus: types.PickedUp,
		}

		buyer := &domain.User{
			ID: buyerID,
		}

		mockNotificationRepoFactory.On("GetRepository", provider).Return(mockNotificationRepo, nil)
		mockUserRepo.On("FindByID", ctx, buyerID).Return(buyer, nil)
		mockNotificationRepo.On("PrNotify", buyer, productRequest).Return(nil)

		err := notificationUsecase.PrNotify(productRequest, provider)

		assert.NoError(t, err)
		mockNotificationRepoFactory.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockNotificationRepo.AssertExpectations(t)
	})

	t.Run("Success_NotifyTravelerCaseNoTraveler", func(t *testing.T) {
		mockNotificationRepoFactory.ExpectedCalls = nil
		mockNotificationRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil
		mockOfferRepo.ExpectedCalls = nil

		provider := "email"
		buyerID := "buyer123"

		productRequest := &domain.ProductRequest{
			UserID:         &buyerID,
			DeliveryStatus: types.Opening,
		}

		mockNotificationRepoFactory.On("GetRepository", provider).Return(mockNotificationRepo, nil)

		err := notificationUsecase.PrNotify(productRequest, provider)

		assert.NoError(t, err)
		mockNotificationRepoFactory.AssertExpectations(t)
	})
	t.Run("Error_RepositoryNotFound", func(t *testing.T) {
		mockNotificationRepoFactory.ExpectedCalls = nil
		mockNotificationRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil
		mockOfferRepo.ExpectedCalls = nil

		provider := "email"
		buyerID := "buyer123"

		productRequest := &domain.ProductRequest{
			UserID:         &buyerID,
			DeliveryStatus: types.Pending,
		}

		mockNotificationRepoFactory.On("GetRepository", provider).Return(&mock_repos.MockBankRepository{}, assert.AnError)
		err := notificationUsecase.PrNotify(productRequest, provider)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockNotificationRepoFactory.AssertExpectations(t)
	})
	t.Run("Error_OfferNotFound", func(t *testing.T) {
		mockNotificationRepoFactory.ExpectedCalls = nil
		mockNotificationRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil
		mockOfferRepo.ExpectedCalls = nil

		provider := "email"
		buyerID := "buyer123"
		selectedOfferID := uint(1)

		productRequest := &domain.ProductRequest{
			UserID:          &buyerID,
			DeliveryStatus:  types.PickedUp,
			SelectedOfferID: &selectedOfferID,
		}

		offer := &domain.Offer{
			Model: gorm.Model{ID: 1},
		}

		mockNotificationRepoFactory.On("GetRepository", provider).Return(mockNotificationRepo, nil)
		mockOfferRepo.On("GetByID", offer).Return(assert.AnError)

		err := notificationUsecase.PrNotify(productRequest, provider)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockNotificationRepoFactory.AssertExpectations(t)
		mockOfferRepo.AssertExpectations(t)
	})
	t.Run("Error_UserNotFound", func(t *testing.T) {
		mockNotificationRepoFactory.ExpectedCalls = nil
		mockNotificationRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil
		mockOfferRepo.ExpectedCalls = nil

		provider := "email"
		buyerID := "buyer123"
		selectedOfferID := uint(1)

		productRequest := &domain.ProductRequest{
			UserID:          &buyerID,
			DeliveryStatus:  types.PickedUp,
			SelectedOfferID: &selectedOfferID,
		}

		offer := &domain.Offer{
			Model: gorm.Model{ID: 1},
		}

		mockNotificationRepoFactory.On("GetRepository", provider).Return(mockNotificationRepo, nil)
		mockOfferRepo.On("GetByID", offer).Return(nil)
		mockUserRepo.On("FindByID", ctx, buyerID).Return(&domain.User{}, assert.AnError)

		err := notificationUsecase.PrNotify(productRequest, provider)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockNotificationRepoFactory.AssertExpectations(t)
		mockOfferRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
	t.Run("Error_CouldNotNotify", func(t *testing.T) {
		mockNotificationRepoFactory.ExpectedCalls = nil
		mockNotificationRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil
		mockOfferRepo.ExpectedCalls = nil

		provider := "email"
		buyerID := "buyer123"
		selectedOfferID := uint(1)

		productRequest := &domain.ProductRequest{
			UserID:          &buyerID,
			DeliveryStatus:  types.PickedUp,
			SelectedOfferID: &selectedOfferID,
		}

		buyer := &domain.User{
			ID: buyerID,
		}

		offer := &domain.Offer{
			Model: gorm.Model{ID: 1},
		}

		mockNotificationRepoFactory.On("GetRepository", provider).Return(mockNotificationRepo, nil)
		mockOfferRepo.On("GetByID", offer).Return(nil)
		mockUserRepo.On("FindByID", ctx, buyerID).Return(buyer, nil)
		mockNotificationRepo.On("PrNotify", buyer, productRequest).Return(assert.AnError)

		err := notificationUsecase.PrNotify(productRequest, provider)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockNotificationRepoFactory.AssertExpectations(t)
		mockOfferRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockNotificationRepo.AssertExpectations(t)
	})
}
