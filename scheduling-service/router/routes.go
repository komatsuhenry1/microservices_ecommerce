package router

import (
	"github.com/gin-gonic/gin"
	"scheduling-service/di"
)

func SetupRoutes(container *di.Container) *gin.Engine {
	router := gin.Default()
	scheduling := router.Group("/scheduling")

	scheduling.POST("/create", container.ScheduleHandler.CreateSchedule)

	// scheduling.POST("/register", container.AuthHandler.UserRegister)
	// scheduling.POST("/login", container.AuthHandler.UserLogin)
	return router
}