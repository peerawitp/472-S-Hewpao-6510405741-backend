package dto

type RegisterUserRequestDTO struct {
	Email      string  `json:"email" validate:"required,email"`
	Password   string  `json:"password" validate:"required,min=8"`
	Name       string  `json:"name" validate:"required"`
	MiddleName *string `json:"middle_name"`
	Surname    string  `json:"surname" validate:"required"`
}

type LoginWithCredentialsRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDTO struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	Name        string  `json:"name"`
	MiddleName  *string `json:"middle_name"`
	Surname     string  `json:"surname"`
	IsVerified  bool    `json:"is_verified"`
	AccessToken string  `json:"access_token"`
}

type LoginWithOAuthRequestDTO struct {
	Provider string `json:"provider" validate:"required"`
	IDToken  string `json:"id_token" validate:"required"`
}
