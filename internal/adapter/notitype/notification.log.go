package notitype

import (
	"log"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/repository"
)

type testNotificationRepo struct {
	cfg *config.Config
}

func NewTestNotificationRepo(cfg *config.Config) (repository.NotificationRepository, error) {
	return &testNotificationRepo{
		cfg: cfg,
	}, nil
}

func (t *testNotificationRepo) PrNotify(user *domain.User, prod *domain.ProductRequest) error {
	log.Println("---------------------------------------------------------------")
	log.Println("send-from: ", t.cfg.EmailUser)
	log.Println("send-to: ", user.Email)
	log.Println("product-request-status: ", prod.DeliveryStatus)
	log.Println("---------------------------------------------------------------")

	return nil
}
