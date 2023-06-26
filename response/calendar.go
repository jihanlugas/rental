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

type Calendar model.CalendarView

type CalendarDetail struct {
	Calendar     model.CalendarView       `json:"calendar"`
	Calendaritem []model.CalendaritemView `json:"calendaritems"`
}
