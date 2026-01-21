package model

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
	Cpf string `json:"cpf"`
	Password string `json:"password"`
}