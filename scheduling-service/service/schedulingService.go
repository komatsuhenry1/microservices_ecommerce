package service

import (
	"scheduling-service/dto"
	"scheduling-service/repository"
)

type SchedulingService interface {
	CreateSchedule(schedule *dto.SchedulingRequestDTO) error
}

type schedulingService struct {
	scheduleRepo repository.ScheduleRepository
}

func NewSchedulingService(scheduleRepo repository.ScheduleRepository) SchedulingService {
	return &schedulingService{scheduleRepo: scheduleRepo}
}

func (s *schedulingService) CreateSchedule(schedule *dto.SchedulingRequestDTO) error {
	return s.scheduleRepo.CreateSchedule(schedule)
}
