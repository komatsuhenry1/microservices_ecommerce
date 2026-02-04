package repository

import (
	"scheduling-service/dto"
	"database/sql"
)

type ScheduleRepository interface {
	CreateSchedule(schedule *dto.SchedulingRequestDTO) error
}

type scheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository(db *sql.DB) ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) CreateSchedule(schedule *dto.SchedulingRequestDTO) error {
	_, err := r.db.Exec("INSERT INTO schedules (user_id, service_id, service_name, barber_id, barber_name, date, time, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", schedule.UserId, schedule.ServiceId, schedule.ServiceName, schedule.BarberId, schedule.BarberName, schedule.Date, schedule.Time, schedule.Status)
	return err
}
