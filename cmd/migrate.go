package main

import (
	"errors"
	"log"

	"github.com/hewpao/hewpao-backend/bootstrap"
	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"gorm.io/gorm"
)

func main() {
	cfg := config.NewConfig()
	db := bootstrap.NewDB(&cfg)

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Account{},
		&domain.ProductRequest{},
		&domain.Offer{},
		&domain.Transaction{},
		&domain.Verification{},
		&domain.Message{},
		&domain.Chat{},
		&domain.Bank{},
		&domain.TravelerPayoutAccount{},
	); err != nil {
		log.Fatal(err)
	}

	log.Println("🚀 Migration completed")

	// Seed Thai banks
	if err := initThaiBanks(db); err != nil {
		log.Fatal(err)
	} else {
		log.Println("🌱 Thai banks seeded")
	}
}

func initThaiBanks(db *gorm.DB) error {
	var count int64
	db.Model(&domain.Bank{}).Count(&count)
	if count > 0 {
		return errors.New("⚠️ Banks are already seeded")
	}

	banks := []domain.Bank{
		{SwiftCode: "BKKBTHBK", NameEN: "Bangkok Bank", NameTH: "ธนาคารกรุงเทพ"},
		{SwiftCode: "AYUDTHBK", NameEN: "Bank of Ayudhya", NameTH: "ธนาคารกรุงศรีอยุธยา"},
		{SwiftCode: "KASITHBK", NameEN: "Kasikorn Bank", NameTH: "ธนาคารกสิกรไทย"},
		{SwiftCode: "KRTHTHBK", NameEN: "Krung Thai Bank", NameTH: "ธนาคารกรุงไทย"},
		{SwiftCode: "SICOTHBK", NameEN: "Siam Commercial Bank", NameTH: "ธนาคารไทยพาณิชย์"},
		{SwiftCode: "TMBKTHBK", NameEN: "TMB Bank", NameTH: "ธนาคารทหารไทย"},
		{SwiftCode: "GSBATHBK", NameEN: "Government Savings Bank", NameTH: "ธนาคารออมสิน"},
		{SwiftCode: "SCBLTHBX", NameEN: "Standard Chartered Bank", NameTH: "ธนาคารสแตนดาร์ดชาร์เตอร์"},
		{SwiftCode: "UOVBTHBK", NameEN: "Union Overseas Bank", NameTH: "ธนาคารยูโอบี"},
		{SwiftCode: "THBKTHBK", NameEN: "Thanachart Bank", NameTH: "ธนาคารธนชาติ"},
		{SwiftCode: "UBOBTHBK", NameEN: "CIMB Thai Bank", NameTH: "ธนาคาร CIMB Thai"},
		{SwiftCode: "CITITHBX", NameEN: "Citibank Thailand", NameTH: "ธนาคาร Citibank Thailand"},
		{SwiftCode: "KIFITHB1", NameEN: "Kiatnakin Bank", NameTH: "ธนาคารเกียรตินาคิน"},
	}

	if err := db.Create(&banks).Error; err != nil {
		return err
	}

	return nil
}
