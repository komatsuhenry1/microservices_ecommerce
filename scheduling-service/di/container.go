package di

import (
	"database/sql"

	"scheduling-service/handler"
	"scheduling-service/repository"
	"scheduling-service/service"
)

type Container struct {
	DB             *sql.DB
	ScheduleRepository repository.ScheduleRepository
	SchedulingService  service.SchedulingService
	ScheduleHandler    *handler.ScheduleHandler
}

func NewContainer(db *sql.DB) *Container {
	scheduleRepo := repository.NewScheduleRepository(db)
	schedulingSvc := service.NewSchedulingService(scheduleRepo)
	scheduleHandler := handler.NewScheduleHandler(schedulingSvc)

	return &Container{
		DB:             db,
		ScheduleRepository: scheduleRepo,
		SchedulingService:  schedulingSvc,
		ScheduleHandler:    scheduleHandler,
	}
}