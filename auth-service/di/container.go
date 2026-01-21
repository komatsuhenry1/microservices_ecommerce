package di

import (
	"database/sql"

	"auth-service/handler"
	"auth-service/repository"
	"auth-service/service"
)

type Container struct {
	DB             *sql.DB
	UserRepository repository.UserRepository
	AuthService    service.AuthService
	AuthHandler    *handler.AuthHandler
}

func NewContainer(db *sql.DB) *Container {
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authSvc)

	return &Container{
		DB:             db,
		UserRepository: userRepo,
		AuthService:    authSvc,
		AuthHandler:    authHandler,
	}
}