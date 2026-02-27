package router

import (
	"auth-service/di"
	"github.com/gin-gonic/gin"
	"auth-service/middleware"
)

func SetupRoutes(container *di.Container) *gin.Engine {
	router := gin.Default()
	auth := router.Group("/auth")
	auth.POST("/register", container.AuthHandler.UserRegister)
	auth.POST("/login", container.AuthHandler.UserLogin)
	auth.GET("/users", middleware.AuthAdmin(), container.AuthHandler.GetUsers)
	return router
}