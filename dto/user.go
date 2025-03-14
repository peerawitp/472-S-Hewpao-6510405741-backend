package dto

type CreateUserDTO struct {
	Email       string  `json:"email"`
	Password    *string `json:"password"`
	Name        string  `json:"name"`
	MiddleName  *string `json:"middle_name"`
	Surname     string  `json:"surname"`
	PhoneNumber *string `json:"phone_number"`
}

type UserProfileDTO struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	Name        string  `json:"name"`
	MiddleName  *string `json:"middle_name"`
	Surname     string  `json:"surname"`
	PhoneNumber *string `json:"phone_number"`
	IsVerified  bool    `json:"is_verified"`
}

type EditProfileDTO struct {
	Name        string  `json:"name" validate:"required"`
	MiddleName  *string `json:"middle_name"`
	Surname     string  `json:"surname" validate:"required"`
	PhoneNumber *string `json:"phone_number"`
}
