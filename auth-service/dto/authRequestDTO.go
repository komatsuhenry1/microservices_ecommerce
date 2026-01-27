package dto


type UserRegisterRequestDTO struct {
	Email string `json:"email"`
	Name     string `json:"name"`
	Cpf      string `json:"cpf"`
	Password string `json:"password"`
	AvatarUrl string `json:"avatar_url"`
}

type UserLoginRequestDTO struct {
	Email string `json:"email"`
	Password string `json:"password"`
}
