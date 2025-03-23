package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/usecase/test/mock_repos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateOffer(t *testing.T) {
	mockOfferRepo := new(mock_repos.MockOfferRepo)
	mockProductRequestRepo := new(mock_repos.MockProductRequestRepo)
	mockUserRepo := new(mock_repos.MockUserRepo)
	ctx := context.Background()

	offerService := usecase.NewOfferService(mockOfferRepo, mockProductRequestRepo, mockUserRepo, ctx)

	t.Run("Success", func(t *testing.T) {
		mockOfferRepo.ExpectedCalls = nil
		mockProductRequestRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil

		travelerID := "traveler123"
		buyerID := "buyer123"
		productRequestID := uint(1)
		offerDate := time.Now()

		createOfferDTO := &dto.CreateOfferDTO{
			ProductRequestID: productRequestID,
			OfferDate:        offerDate,
		}

		productRequest := &domain.ProductRequest{
			Model: gorm.Model{
				ID: 1,
			},
			UserID: &buyerID,
		}

		user := &domain.User{
			ID: travelerID,
		}

		mockProductRequestRepo.On("FindByID", int(productRequestID)).Return(productRequest, nil)
		mockUserRepo.On("FindByID", ctx, travelerID).Return(user, nil)
		mockOfferRepo.On("Create", mock.AnythingOfType("*domain.Offer")).Return(nil)

		err := offerService.CreateOffer(createOfferDTO, travelerID)

		assert.NoError(t, err)
		mockProductRequestRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockOfferRepo.AssertExpectations(t)
	})

	t.Run("Error_SelfOffer", func(t *testing.T) {
		mockOfferRepo.ExpectedCalls = nil
		mockProductRequestRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil

		travelerID := "traveler123"
		productRequestID := uint(1)
		offerDate := time.Now()

		createOfferDTO := &dto.CreateOfferDTO{
			ProductRequestID: productRequestID,
			OfferDate:        offerDate,
		}

		productRequest := &domain.ProductRequest{
			Model: gorm.Model{
				ID: 1,
			},
			UserID: &travelerID,
		}

		user := &domain.User{
			ID: travelerID,
		}

		mockProductRequestRepo.On("FindByID", int(productRequestID)).Return(productRequest, nil)
		mockUserRepo.On("FindByID", ctx, travelerID).Return(user, nil)
		mockOfferRepo.On("Create", mock.AnythingOfType("*domain.Offer")).Return(nil)

		err := offerService.CreateOffer(createOfferDTO, travelerID)

		assert.Error(t, err)
		assert.Equal(t, exception.ErrCouldNotSelfOffer, err)
		mockProductRequestRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockOfferRepo.AssertExpectations(t)
	})

	t.Run("Error_ProductRequestNotFound", func(t *testing.T) {
		mockOfferRepo.ExpectedCalls = nil
		mockProductRequestRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil

		travelerID := "traveler123"
		productRequestID := uint(1)
		offerDate := time.Now()

		createOfferDTO := &dto.CreateOfferDTO{
			ProductRequestID: productRequestID,
			OfferDate:        offerDate,
		}

		mockProductRequestRepo.On("FindByID", int(productRequestID)).Return(&domain.ProductRequest{}, assert.AnError)

		err := offerService.CreateOffer(createOfferDTO, travelerID)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockProductRequestRepo.AssertExpectations(t)
	})

	t.Run("Error_UserNotFound", func(t *testing.T) {
		mockOfferRepo.ExpectedCalls = nil
		mockProductRequestRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil

		travelerID := "traveler123"
		buyerID := "buyer123"
		productRequestID := uint(1)
		offerDate := time.Now()

		createOfferDTO := &dto.CreateOfferDTO{
			ProductRequestID: productRequestID,
			OfferDate:        offerDate,
		}

		productRequest := &domain.ProductRequest{
			Model: gorm.Model{
				ID: 1,
			},
			UserID: &buyerID,
		}

		mockProductRequestRepo.On("FindByID", int(productRequestID)).Return(productRequest, nil)
		mockUserRepo.On("FindByID", ctx, travelerID).Return(&domain.User{}, assert.AnError)

		err := offerService.CreateOffer(createOfferDTO, travelerID)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockProductRequestRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error_CouldNotCreateOffer", func(t *testing.T) {
		mockOfferRepo.ExpectedCalls = nil
		mockProductRequestRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil

		travelerID := "traveler123"
		buyerID := "buyer123"
		productRequestID := uint(1)
		offerDate := time.Now()

		createOfferDTO := &dto.CreateOfferDTO{
			ProductRequestID: productRequestID,
			OfferDate:        offerDate,
		}

		productRequest := &domain.ProductRequest{
			Model: gorm.Model{
				ID: 1,
			},
			UserID: &buyerID,
		}

		user := &domain.User{
			ID: travelerID,
		}

		mockProductRequestRepo.On("FindByID", int(productRequestID)).Return(productRequest, nil)
		mockUserRepo.On("FindByID", ctx, travelerID).Return(user, nil)
		mockOfferRepo.On("Create", mock.AnythingOfType("*domain.Offer")).Return(assert.AnError)

		err := offerService.CreateOffer(createOfferDTO, travelerID)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockProductRequestRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockOfferRepo.AssertExpectations(t)
	})
}

func TestGetOfferDetailByOfferID(t *testing.T) {
	mockOfferRepo := new(mock_repos.MockOfferRepo)
	mockProductRequestRepo := new(mock_repos.MockProductRequestRepo)
	mockUserRepo := new(mock_repos.MockUserRepo)
	ctx := context.Background()

	offerService := usecase.NewOfferService(mockOfferRepo, mockProductRequestRepo, mockUserRepo, ctx)

	t.Run("Success", func(t *testing.T) {
		mockOfferRepo.ExpectedCalls = nil
		mockProductRequestRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil

		buyerID := "buyer123"
		offerID := 1
		productRequestID := 1

		mockOfferRepo.On("GetOfferDetailByOfferID", offerID).Return(&domain.Offer{
			ProductRequest: &domain.ProductRequest{
				Model: gorm.Model{
					ID: uint(productRequestID),
				},
			},
		}, nil)
		mockProductRequestRepo.On("IsOwnedByUser", productRequestID, buyerID).Return(true, nil)

		_, err := offerService.GetOfferDetailByOfferID(offerID, buyerID)

		assert.NoError(t, err)
		mockOfferRepo.AssertExpectations(t)
		mockProductRequestRepo.AssertExpectations(t)
	})
	t.Run("Err_OfferNotFound", func(t *testing.T) {
		mockOfferRepo.ExpectedCalls = nil
		mockProductRequestRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil

		buyerID := "buyer123"
		offerID := 1

		mockOfferRepo.On("GetOfferDetailByOfferID", offerID).Return(&domain.Offer{}, assert.AnError)

		_, err := offerService.GetOfferDetailByOfferID(offerID, buyerID)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockOfferRepo.AssertExpectations(t)
	})
	t.Run("Err_CouldNotVerifyOwner", func(t *testing.T) {
		mockOfferRepo.ExpectedCalls = nil
		mockProductRequestRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil

		buyerID := "buyer123"
		offerID := 1
		productRequestID := 1

		mockOfferRepo.On("GetOfferDetailByOfferID", offerID).Return(&domain.Offer{
			ProductRequest: &domain.ProductRequest{
				Model: gorm.Model{
					ID: uint(productRequestID),
				},
			},
		}, nil)
		mockProductRequestRepo.On("IsOwnedByUser", productRequestID, buyerID).Return(false, assert.AnError)

		_, err := offerService.GetOfferDetailByOfferID(offerID, buyerID)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockOfferRepo.AssertExpectations(t)
		mockProductRequestRepo.AssertExpectations(t)
	})
	t.Run("Err_NotOwnedProductRequest", func(t *testing.T) {
		mockOfferRepo.ExpectedCalls = nil
		mockProductRequestRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil

		buyerID := "buyer123"
		offerID := 1
		productRequestID := 1

		mockOfferRepo.On("GetOfferDetailByOfferID", offerID).Return(&domain.Offer{
			ProductRequest: &domain.ProductRequest{
				Model: gorm.Model{
					ID: uint(productRequestID),
				},
			},
		}, nil)
		mockProductRequestRepo.On("IsOwnedByUser", productRequestID, buyerID).Return(false, nil)

		_, err := offerService.GetOfferDetailByOfferID(offerID, buyerID)

		assert.Error(t, err)
		assert.Equal(t, exception.ErrPermissionDenied, err)
		mockOfferRepo.AssertExpectations(t)
		mockProductRequestRepo.AssertExpectations(t)
	})
}
