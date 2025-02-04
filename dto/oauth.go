package dto

type OAuthClaims struct {
	Name       string  `json:"name"`
	MiddleName *string `json:"middle_name"`
	Surname    string  `json:"surname"`
	Email      string  `json:"email"`
	Picture    string  `json:"picture"`
}
