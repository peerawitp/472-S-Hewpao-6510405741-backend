package domain

type Bank struct {
	SwiftCode string `gorm:"size:20;primaryKey"`
	NameEN    string `gorm:"size:100;not null"`
	NameTH    string `gorm:"size:100;not null"`
}
