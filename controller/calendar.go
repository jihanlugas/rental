package controller

import (
	"github.com/jihanlugas/rental/db"
	"github.com/jihanlugas/rental/model"
	"github.com/jihanlugas/rental/response"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Calendar struct{}

func CalendarComposer() Calendar {
	return Calendar{}
}

func (h Calendar) Tes(c echo.Context) error {
	var err error
	var ListCalendar []model.Calendar
	var ListUser []model.UsersView

	conn, closeConn := db.GetConnection()
	defer closeConn()
	if err != nil {
		errorInternal(c, err)
	}

	err = conn.Find(&ListCalendar).Error
	if err != nil {
		errorInternal(c, err)
	}

	err = conn.Find(&ListUser).Error
	if err != nil {
		errorInternal(c, err)
	}

	resUser := response.Users(ListUser)

	return response.Success(http.StatusOK, "success", response.Payload{
		"calendars": ListCalendar,
		"users":     resUser,
	}).SendJSON(c)
}
