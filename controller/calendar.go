package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jihanlugas/rental/db"
	"github.com/jihanlugas/rental/model"
	"github.com/jihanlugas/rental/request"
	"github.com/jihanlugas/rental/response"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
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
		Where("delete_dt IS NULL ").
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

// GetById godoc
// @Tags Calendar
// @Summary To do get a calendar
// @Accept json
// @Produce json
// @Param id path string true "Calendar ID"
// @Success      200  {object}	response.Response{payload=response.Calendar}
// @Failure      500  {object}  response.Response
// @Router /calendar/{id} [get]
func (h Calendar) GetById(c echo.Context) error {
	var err error
	var calendar model.CalendarView

	conn, closeConn := db.GetConnection()
	defer closeConn()

	ID := c.Param("id")
	if ID == "" {
		return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
	}

	err = conn.Where("id = ? ", ID).First(&calendar).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	res := response.Calendar(calendar)

	return response.Success(http.StatusOK, "success", res).SendJSON(c)
}

// GetDetailById godoc
// @Tags Calendar
// @Summary To do get a calendar detail
// @Accept json
// @Produce json
// @Param id path string true "Calendar ID"
// @Success      200  {object}	response.Response{payload=response.CalendarDetail}
// @Failure      500  {object}  response.Response
// @Router /calendar/detail/{id} [get]
func (h Calendar) GetDetailById(c echo.Context) error {
	var err error
	var calendar model.CalendarView
	var calendarItems []model.CalendaritemView

	conn, closeConn := db.GetConnection()
	defer closeConn()

	ID := c.Param("id")
	if ID == "" {
		return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
	}

	err = conn.Where("id = ? ", ID).First(&calendar).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	err = conn.Where("calendar_id = ? ", ID).Find(&calendarItems).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	res := response.CalendarDetail{
		Calendar:     calendar,
		Calendaritem: calendarItems,
	}

	return response.Success(http.StatusOK, "success", res).SendJSON(c)
}

// Create godoc
// @Tags Calendar
// @Summary To do create new calendar event
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.CreateCalendar true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /calendar [post]
func (h Calendar) Create(c echo.Context) error {
	var err error
	var calendar model.Calendar

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	req := new(request.CreateCalendar)
	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return response.Error(http.StatusBadRequest, "error validation", response.ValidationError(err)).SendJSON(c)
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	calendar = model.Calendar{
		CompanyID:  loginUser.CompanyID,
		PropertyID: req.PropertyID,
		Name:       req.Name,
		StartDt:    *req.StartDt,
		EndDt:      *req.EndDt,
		Status:     0,
		CreateBy:   loginUser.UserID,
		UpdateBy:   loginUser.UserID,
	}

	tx.Save(&calendar)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusCreated, "success", response.Payload{}).SendJSON(c)
}

// Update godoc
// @Tags Calendar
// @Summary To do update a calendar
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Calendar ID"
// @Param req body request.UpdateCalendar true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /calendar/{id} [put]
func (h Calendar) Update(c echo.Context) error {
	var err error
	var calendar model.Calendar

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	ID := c.Param("id")
	if ID == "" {
		return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
	}

	req := new(request.UpdateCalendar)
	if err = c.Bind(req); err != nil {
		return err
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Where("id = ? ", ID).First(&calendar).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	tx := conn.Begin()

	calendar.PropertyID = req.PropertyID
	calendar.Name = req.Name
	calendar.StartDt = *req.StartDt
	calendar.EndDt = *req.EndDt
	calendar.UpdateBy = loginUser.UserID

	tx.Save(&calendar)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusAccepted, "success", response.Payload{}).SendJSON(c)
}

// Delete Calendar
// @Summary Delete Calendar
// @Tags Calendar
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Calendar ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /calendar/{id} [delete]
func (h Calendar) Delete(c echo.Context) error {
	var err error
	var calendar model.Calendar

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	ID := c.Param("id")
	if ID == "" {
		return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Where("id = ? ", ID).First(&calendar).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	tx := conn.Begin()

	now := time.Now()

	calendar.DeleteBy = loginUser.UserID
	calendar.DeleteDt = &now
	tx.Save(&calendar)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusAccepted, "success", response.Payload{}).SendJSON(c)
}

// Timeline
// @Tags Calendar
// @Summary To do get data timeline
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.Timeline true "json req body"
// @Success      200  {object}	response.Response{payload=response.Timeline}
// @Failure      500  {object}  response.Response
// @Router /calendar/timeline [post]
func (h Calendar) Timeline(c echo.Context) error {
	var err error
	var ListCalendar []model.CalendarView
	var ListProperty []model.PropertyView

	req := new(request.Timeline)
	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return response.Error(http.StatusBadRequest, "error validation", response.ValidationError(err)).SendJSON(c)
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Where("company_id = ? ", req.CompanyID).
		Where("start_dt <= ? ", req.EndDt).
		Where("end_dt >= ? ", req.StartDt).
		Where("delete_dt IS NULL ").
		Find(&ListCalendar).Error
	if err != nil {
		return response.Error(http.StatusInternalServerError, "error", response.Payload{})
	}

	err = conn.Where("company_id = ? ", req.CompanyID).
		Where("delete_dt IS NULL ").
		Find(&ListProperty).Error
	if err != nil {
		return response.Error(http.StatusInternalServerError, "error", response.Payload{})
	}

	res := response.Timeline{
		ListCalendar: ListCalendar,
		ListProperty: ListProperty,
	}

	return response.Success(http.StatusOK, "success", res).SendJSON(c)
}

//func (h Calendar) Tes(c echo.Context) error {
//	var err error
//	var ListCalendar []model.CalendarView
//	var ListProperty []model.PropertyView
//
//	conn, closeConn := db.GetConnection()
//	defer closeConn()
//	if err != nil {
//		errorInternal(c, err)
//	}
//
//	err = conn.Find(&ListCalendar).Error
//	if err != nil {
//		errorInternal(c, err)
//	}
//
//	err = conn.Find(&ListProperty).Error
//	if err != nil {
//		errorInternal(c, err)
//	}
//
//	return response.Success(http.StatusOK, "success", response.Payload{
//		"calendars":  ListCalendar,
//		"properties": ListProperty,
//	}).SendJSON(c)
//}
