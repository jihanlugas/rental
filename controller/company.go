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

type Company struct{}

func CompanyComposer() Company {
	return Company{}
}

// GetById godoc
// @Tags Company
// @Summary To do get a user
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /company/{id} [get]
func (h Company) GetById(c echo.Context) error {
	var err error
	var company model.CompanyView

	conn, closeConn := db.GetConnection()
	defer closeConn()

	ID := c.Param("id")
	if ID == "" {
		return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
	}

	err = conn.Where("id = ? ", ID).First(&company).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	return response.Success(http.StatusOK, "success", company).SendJSON(c)
}

// Create godoc
// @Tags Company
// @Summary To do create new election
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.CreateCompany true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /company [post]
func (h Company) Create(c echo.Context) error {
	var err error
	var company model.Company

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	req := new(request.CreateCompany)
	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return response.Error(http.StatusBadRequest, "error validation", response.ValidationError(err)).SendJSON(c)
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	company = model.Company{
		Name:     req.Name,
		UserID:   req.UserID,
		CreateBy: loginUser.UserID,
		UpdateBy: loginUser.UserID,
	}

	tx.Save(&company)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusCreated, "success", response.Payload{}).SendJSON(c)
}

// Update godoc
// @Tags Company
// @Summary To do update a company
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param req body request.UpdateCompany true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /company/{id} [put]
func (h Company) Update(c echo.Context) error {
	var err error
	var company model.Company

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	ID := c.Param("id")
	if ID == "" {
		return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
	}

	req := new(request.UpdateCompany)
	if err = c.Bind(req); err != nil {
		return err
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Where("id = ? ", ID).First(&company).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	tx := conn.Begin()

	company.Name = req.Name
	company.UserID = req.UserID
	company.UpdateBy = loginUser.UserID

	tx.Save(&company)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusAccepted, "success", response.Payload{}).SendJSON(c)
}

// Delete Company
// @Summary Delete Company
// @Tags Company
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /company/{id} [delete]
func (h Company) Delete(c echo.Context) error {
	var err error
	var company model.Company

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

	err = conn.Where("id = ? ", ID).First(&company).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	tx := conn.Begin()

	now := time.Now()

	company.DeleteBy = loginUser.UserID
	company.DeleteDt = &now
	tx.Save(&company)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusAccepted, "success", response.Payload{}).SendJSON(c)
}
