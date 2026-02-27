package router

import (
	"github.com/gin-gonic/gin"
	"scheduling-service/di"
	"scheduling-service/middleware"
)

func SetupRoutes(container *di.Container) *gin.Engine {
	router := gin.Default()
	scheduling := router.Group("/scheduling")

	// add middlware on this route

	scheduling.POST("/create", middleware.AuthUserOrAdmin(), container.ScheduleHandler.CreateSchedule) // add middlware on this route
	scheduling.GET("/dashboard", middleware.AuthAdmin(), container.ScheduleHandler.GetDashboard) // add middlware on this route

	// scheduling.POST("/register", container.AuthHandler.UserRegister)
	// scheduling.POST("/login", container.AuthHandler.UserLogin)
	return router
}