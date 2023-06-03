package request

import "time"

type WsCalendar struct {
	CompanyID string `query:"companyId" validate:"required"`
	//StartDt time.Time `json:"sndDt" validate:"required"`
	//EndDt   time.Time `json:"endDt" validate:"required"`
}

type WsCalendarMessage struct {
	StartDt time.Time `json:"startDt"`
	EndDt   time.Time `json:"endDt"`
}

type CreateCalendar struct {
	PropertyID string     `json:"propertyId" validate:"required"`
	Name       string     `json:"name" validate:"required"`
	StartDt    *time.Time `json:"startDt" validate:"required"`
	EndDt      *time.Time `json:"endDt" validate:"required"`
}

type UpdateCalendar struct {
	PropertyID string     `json:"propertyId" validate:"required"`
	Name       string     `json:"name" validate:"required"`
	StartDt    *time.Time `json:"startDt" validate:"required"`
	EndDt      *time.Time `json:"endDt" validate:"required"`
}

type Timeline struct {
	CompanyID string     `json:"companyId" validate:"required"`
	StartDt   *time.Time `json:"startDt" validate:"required"`
	EndDt     *time.Time `json:"endDt" validate:"required"`
}
