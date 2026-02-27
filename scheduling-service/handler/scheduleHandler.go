package handler

import (
	"scheduling-service/dto"
	"scheduling-service/service"
	"scheduling-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	schedulingService service.SchedulingService
}

func NewScheduleHandler(schedulingService service.SchedulingService) *ScheduleHandler {
	return &ScheduleHandler{schedulingService: schedulingService}
}

func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var schedule dto.SchedulingRequestDTO
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.schedulingService.CreateSchedule(&schedule); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Schedule created successfully"})
}

func (h *ScheduleHandler) GetDashboard(c *gin.Context) {
	dashbordData, err := h.schedulingService.GetDashboard()
	if err != nil {
		utils.SendErrorResponse(c, "Failed to get dashboard data", http.StatusInternalServerError)
		return
	}

	c.JSON(200, dashbordData)
}