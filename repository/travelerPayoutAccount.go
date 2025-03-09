package repository

import (
	"context"

	"github.com/hewpao/hewpao-backend/domain"
)

type TravelerPayoutAccountRepository interface {
	Store(ctx context.Context, travelerPayoutAccount *domain.TravelerPayoutAccount) error
	FindByUserID(ctx context.Context, userID string) ([]domain.TravelerPayoutAccount, error)
}
