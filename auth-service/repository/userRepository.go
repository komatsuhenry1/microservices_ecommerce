package repository

import (
	"database/sql"
	"auth-service/model"
)

type UserRepository interface {
	CreateUser(user *model.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *model.User) error {
	_, err := r.db.Exec("INSERT INTO users (email, name, cpf, password) VALUES ($1, $2, $3, $4)", user.Email, user.Name, user.Cpf, user.Password)
	return err
}