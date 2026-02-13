package repository

import (
	"auth-service/model"
	"database/sql"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *model.User) error {
	_, err := r.db.Exec("INSERT INTO users (email, name, cpf, password, role, avatar_url) VALUES ($1, $2, $3, $4, $5, $6)", user.Email, user.Name, user.Cpf, user.Password, user.Role, user.AvatarUrl)
	return err
}

func (r *userRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	row := r.db.QueryRow("SELECT * FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Cpf, &user.Password, &user.AvatarUrl, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}
