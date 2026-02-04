package handler

import (
	"scheduling-service/dto"
	"scheduling-service/service"

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
