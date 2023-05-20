package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jihanlugas/rental/db"
	"github.com/jihanlugas/rental/model"
	"github.com/jihanlugas/rental/request"
	"github.com/jihanlugas/rental/response"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Calendar struct{}

func CalendarComposer() Calendar {
	return Calendar{}
}

var (
	upgrader = websocket.Upgrader{}
)

func getCalendarData(CompanyID string, CalendarMessage request.WsCalendarMessage) *response.Response {

	var err error
	var ListCalendar []model.CalendarView
	var ListProperty []model.PropertyView

	conn, closeConn := db.GetConnection()
	defer closeConn()
	if err != nil {
		return response.Error(http.StatusInternalServerError, "error", response.Payload{})
	}

	err = conn.Where("company_id = ? ", CompanyID).
		Where("start_dt <= ? ", CalendarMessage.EndDt).
		Where("end_dt >= ? ", CalendarMessage.StartDt).
		Find(&ListCalendar).Error
	if err != nil {
		return response.Error(http.StatusInternalServerError, "error", response.Payload{})
	}

	err = conn.Where("company_id = ? ", CompanyID).Find(&ListProperty).Error
	if err != nil {
		return response.Error(http.StatusInternalServerError, "error", response.Payload{})
	}

	data := response.WsCalendar{
		ListCalendar: ListCalendar,
		ListProperty: ListProperty,
	}

	res := response.Success(http.StatusOK, "success", data)

	return res
}

func writer(conn *websocket.Conn, done chan struct{}, data chan interface{}) {
	defer conn.Close()

	for {
		select {
		case <-done:
			// the reader is done, so return
			return
		case message := <-data: // get data from channel
			byte, err := json.Marshal(&message)
			err = conn.WriteMessage(1, byte)
			if err != nil {
				fmt.Printf("there is errors %s \n", err)
				return
			}
		}
	}
}

func reader(conn *websocket.Conn, done chan struct{}, data chan interface{}, req *request.WsCalendar) {
	var reqMessage request.WsCalendarMessage
	defer conn.Close()
	defer close(done)
	for {
		_, msg, err := conn.ReadMessage() //what is message type?
		if err != nil {
			fmt.Println("there is errors%s", err)
			return
		}

		reqMessage = request.WsCalendarMessage{}
		err = json.Unmarshal(msg, &reqMessage)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		data <- getCalendarData(req.CompanyID, reqMessage)
	}
}
func (h Calendar) WsCalendar(c echo.Context) error {
	var err error
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	req := new(request.WsCalendar)
	if err = c.Bind(req); err != nil {
		errorInternal(c, err)
	}

	if err = c.Validate(req); err != nil {
		return response.Error(http.StatusBadRequest, "error validation", response.ValidationError(err)).SendJSON(c)
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	done := make(chan struct{})
	data := make(chan interface{})

	go writer(ws, done, data)
	go reader(ws, done, data, req)

	return nil
}

func (h Calendar) Tes(c echo.Context) error {
	var err error
	var ListCalendar []model.CalendarView
	var ListProperty []model.PropertyView

	conn, closeConn := db.GetConnection()
	defer closeConn()
	if err != nil {
		errorInternal(c, err)
	}

	err = conn.Find(&ListCalendar).Error
	if err != nil {
		errorInternal(c, err)
	}

	err = conn.Find(&ListProperty).Error
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusOK, "success", response.Payload{
		"calendars":  ListCalendar,
		"properties": ListProperty,
	}).SendJSON(c)
}
