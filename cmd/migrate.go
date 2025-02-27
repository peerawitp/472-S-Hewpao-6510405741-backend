package main

import (
	"log"

	"github.com/hewpao/hewpao-backend/bootstrap"
	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
)

func main() {
	cfg := config.NewConfig()
	db := bootstrap.NewDB(&cfg)

	if err := db.AutoMigrate(&domain.User{}, &domain.Account{}, &domain.ProductRequest{}, &domain.Offer{}, &domain.Transaction{}); err != nil {
		log.Fatal(err)
	}

	log.Println("🚀 Migration completed")
}
