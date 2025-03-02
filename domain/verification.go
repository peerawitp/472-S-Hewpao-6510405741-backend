package domain

import "gorm.io/gorm"

type Verification struct {
	gorm.Model
	CardImage *string

	IDNumber    string `gorm:"unique;not null"`
	FirstNameTh string
	LastNameTh  string

	FirstNameEn string
	LastNameEn  string

	Gender string

	DOBTh string
	DOBEn string

	ExpireTh string
	ExpireEn string

	IssueTh string
	IssueEn string

	Address     string
	SubDistrict string
	District    string
	Province    string
	PostalCode  string

	UserID string `gorm:"unique;not null"`
}
