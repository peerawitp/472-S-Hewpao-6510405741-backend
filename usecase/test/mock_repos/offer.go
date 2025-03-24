package mock_repos

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/stretchr/testify/mock"
)

type MockOfferRepo struct {
	mock.Mock
}

func (mo *MockOfferRepo) Create(req *domain.Offer) error {
	args := mo.Called(req)
	return args.Error(0)
}

func (mo *MockOfferRepo) GetByID(req *domain.Offer) error {
	args := mo.Called(req)
	req.UserID = "traveler123"
	return args.Error(0)
}

func (mo *MockOfferRepo) GetOfferDetailByOfferID(offerID int) (*domain.Offer, error) {
	args := mo.Called(offerID)
	return args.Get(0).(*domain.Offer), args.Error(1)
}
