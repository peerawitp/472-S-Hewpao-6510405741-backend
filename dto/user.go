package dto

type CreateUserDTO struct {
	Email       string  `json:"email"`
	Password    *string `json:"password"`
	Name        string  `json:"name"`
	MiddleName  *string `json:"middle_name"`
	Surname     string  `json:"surname"`
	PhoneNumber *string `json:"phone_number"`
}
