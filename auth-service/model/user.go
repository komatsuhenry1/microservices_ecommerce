package model

import "time"

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
	Role string `json:"role"`
	Cpf string `json:"cpf"`
	Password string `json:"password"`
	AvatarUrl string `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}