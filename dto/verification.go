package dto

import (
	"github.com/gofiber/fiber/v2"
)

type VerifyWithKYCDTO struct {
	CardImage *fiber.FormFile `form:"card-image" validate:"required"`
}

type UpdateUserVerificationDTO struct {
	Isverified bool `json:"is_verified"`
}

type EKYCResponseDTO struct {
	Address        string  `json:"address"`
	Alley          string  `json:"alley"`
	DetectionScore float64 `json:"detection_score"`
	District       string  `json:"district"`
	EnDOB          string  `json:"en_dob"`
	EnExpire       string  `json:"en_expire"`
	EnFName        string  `json:"en_fname"`
	EnInit         string  `json:"en_init"`
	EnIssue        string  `json:"en_issue"`
	EnLName        string  `json:"en_lname"`
	EnName         string  `json:"en_name"`
	ErrorMessage   string  `json:"error_message"`
	Face           string  `json:"face"` // Base64 image
	Gender         string  `json:"gender"`
	HomeAddress    string  `json:"home_address"`
	HouseNo        string  `json:"house_no"`
	IDNumber       string  `json:"id_number"`
	IDNumberStatus int     `json:"id_number_status"`
	Lane           string  `json:"lane"`
	PostalCode     string  `json:"postal_code"`
	ProcessTime    float64 `json:"process_time"`
	Province       string  `json:"province"`
	Religion       string  `json:"religion"`
	RequestID      string  `json:"request_id"`
	Road           string  `json:"road"`
	SubDistrict    string  `json:"sub_district"`
	ThDOB          string  `json:"th_dob"`
	ThExpire       string  `json:"th_expire"`
	ThFName        string  `json:"th_fname"`
	ThInit         string  `json:"th_init"`
	ThIssue        string  `json:"th_issue"`
	ThLName        string  `json:"th_lname"`
	ThName         string  `json:"th_name"`
	Village        string  `json:"village"`
	VillageNo      string  `json:"village_no"`
}
