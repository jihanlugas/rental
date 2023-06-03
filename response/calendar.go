package response

import "github.com/jihanlugas/rental/model"

type WsCalendar struct {
	ListCalendar []model.CalendarView `json:"calendars"`
	ListProperty []model.PropertyView `json:"properties"`
}

type Timeline struct {
	ListCalendar []model.CalendarView `json:"calendars"`
	ListProperty []model.PropertyView `json:"properties"`
}
