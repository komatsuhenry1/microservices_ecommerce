package dto


type UserRegisterRequestDTO struct {
	Email string `json:"email"`
	Name     string `json:"name"`
	Cpf      string `json:"cpf"`
	Password string `json:"password"`
}
