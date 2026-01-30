package router

import (
	"auth-service/di"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(container *di.Container) *gin.Engine {
	router := gin.Default()
	auth := router.Group("/auth")
	auth.POST("/register", container.AuthHandler.UserRegister)
	auth.POST("/login", container.AuthHandler.UserLogin)
	return router
}