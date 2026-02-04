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
