package request

import "time"

type WsCalendar struct {
	CompanyID string `query:"companyId"`
	//StartDt time.Time `json:"sndDt" validate:"required"`
	//EndDt   time.Time `json:"endDt" validate:"required"`
}

type WsCalendarMessage struct {
	StartDt time.Time `json:"startDt"`
	EndDt   time.Time `json:"endDt"`
}
