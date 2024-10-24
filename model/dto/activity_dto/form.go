package activity_dto

type SaveRequestDTO struct {
	Id           uint    `json:"id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Cron         string  `json:"cron"`
	Advance      float64 `json:"advance"`
	Interval     float64 `json:"interval"`
	UrlPattern   string  `json:"urlPattern"`
	QueryPattern string  `json:"queryPattern"`
	Type         int     `json:"type"`
	Field        string  `json:"field"`
	AlertAhead   int     `json:"alertAhead"`
	Status       int     `json:"status"`
}
