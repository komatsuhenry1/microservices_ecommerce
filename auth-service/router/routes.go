package router

import (
	"auth-service/di"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(container *di.Container) *gin.Engine {
	router := gin.Default()
	router.POST("/register", container.AuthHandler.UserRegister)
	router.POST("/login", container.AuthHandler.UserLogin)
	return router
}