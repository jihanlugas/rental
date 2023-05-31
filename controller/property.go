package controller

import (
	"errors"
	"github.com/jihanlugas/rental/db"
	"github.com/jihanlugas/rental/model"
	"github.com/jihanlugas/rental/request"
	"github.com/jihanlugas/rental/response"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Property struct{}

func PropertyComposer() Property {
	return Property{}
}

// GetById godoc
// @Tags Property
// @Summary To do get a user
// @Accept json
// @Produce json
// @Param id path string true "Property ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /property/{id} [get]
func (h Property) GetById(c echo.Context) error {
	var err error
	var property model.PropertyView

	conn, closeConn := db.GetConnection()
	defer closeConn()

	ID := c.Param("id")
	if ID == "" {
		return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
	}

	err = conn.Where("id = ? ", ID).First(&property).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	return response.Success(http.StatusOK, "success", property).SendJSON(c)
}

// Create godoc
// @Tags Property
// @Summary To do create new election
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.CreateProperty true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /property [post]
func (h Property) Create(c echo.Context) error {
	var err error
	var property model.Property

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	req := new(request.CreateProperty)
	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return response.Error(http.StatusBadRequest, "error validation", response.ValidationError(err)).SendJSON(c)
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	property = model.Property{
		CompanyID:   loginUser.CompanyID,
		Name:        req.Name,
		Description: req.Description,
		CreateBy:    loginUser.UserID,
		UpdateBy:    loginUser.UserID,
	}

	tx.Save(&property)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusCreated, "success", response.Payload{}).SendJSON(c)
}

// Update godoc
// @Tags Property
// @Summary To do update a property
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Property ID"
// @Param req body request.UpdateProperty true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /property/{id} [put]
func (h Property) Update(c echo.Context) error {
	var err error
	var property model.Property

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	ID := c.Param("id")
	if ID == "" {
		return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
	}

	req := new(request.UpdateProperty)
	if err = c.Bind(req); err != nil {
		return err
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Where("id = ? ", ID).First(&property).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	tx := conn.Begin()

	property.Name = req.Name
	property.Description = req.Description
	property.UpdateBy = loginUser.UserID

	tx.Save(&property)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusAccepted, "success", response.Payload{}).SendJSON(c)
}

// Delete Property
// @Summary Delete Property
// @Tags Property
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Property ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /property/{id} [delete]
func (h Property) Delete(c echo.Context) error {
	var err error
	var property model.Property

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

	err = conn.Where("id = ? ", ID).First(&property).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	tx := conn.Begin()

	now := time.Now()

	property.DeleteBy = loginUser.UserID
	property.DeleteDt = &now
	tx.Save(&property)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusAccepted, "success", response.Payload{}).SendJSON(c)
}
