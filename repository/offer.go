package repository

import (
	"github.com/hewpao/hewpao-backend/domain"
)

type OfferRepository interface {
	Create(req *domain.Offer) error
	GetByID(req *domain.Offer) error
	GetOfferDetailByOfferID(offerID int) (*domain.Offer, error)
}
