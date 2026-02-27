package dto

type SchedulingResponseDTO struct {
	UserId      int    `json:"user_id"`
	ServiceId   int    `json:"service_id"`
	ServiceName string `json:"service_name"`
	BarberId    int    `json:"barber_id"`
	BarberName  string `json:"barber_name"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Status      string `json:"status"`
}

type DashboardResponseDTO struct {
	TotalSchedules int `json:"total_schedules"`
	PendingSchedules int `json:"pending_schedules"`
	CompletedSchedules int `json:"completed_schedules"`
	CancelledSchedules int `json:"cancelled_schedules"`
}
