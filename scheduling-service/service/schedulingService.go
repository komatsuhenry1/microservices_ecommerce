package service

import (
	"scheduling-service/dto"
	"scheduling-service/repository"
)

type SchedulingService interface {
	CreateSchedule(schedule *dto.SchedulingRequestDTO) error
	GetDashboard() (*dto.DashboardResponseDTO, error)

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

func (s *schedulingService) GetDashboard() (*dto.DashboardResponseDTO, error) {

	dashbboardData := dto.DashboardResponseDTO{
		TotalSchedules: 0,
		PendingSchedules: 0,
		CompletedSchedules: 0,
		CancelledSchedules: 0,
	}

	return &dashbboardData, nil
}

