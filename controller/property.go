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

// Page Property
// @Summary Page Property
// @Tags Property
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.PageProperty true "payload"
// @Success      200  {object}	response.Response{payload=response.Pagination}
// @Failure      500  {object}  response.Response
// @Router /property/page [post]
func (h Property) Page(c echo.Context) error {
	var err error

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	req := new(request.PageProperty)
	if err = c.Bind(req); err != nil {
		errorInternal(c, err)
	}

	if err = c.Validate(req); err != nil {
		return response.Error(http.StatusBadRequest, "error validation", response.ValidationError(err)).SendJSON(c)
	}

	err, cnt, list := getPageProperty(req, loginUser)
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusOK, "success", response.PayloadPagination(req, list, cnt)).SendJSON(c)
}

func getPageProperty(req *request.PageProperty, loginUser UserLogin) (error, int64, response.PageProperty) {
	var err error
	var cnt int64
	var res response.PageProperty

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Model(&res).
		Where("company_id = ? ", loginUser.CompanyID).
		Where("name LIKE ? ", "%"+req.Name+"%").
		Where("delete_dt IS NULL").
		Count(&cnt).Error
	if err != nil {
		return err, cnt, res
	}

	// get data
	if req.GetPage() < 1 {
		req.SetPage(1)
	}

	offsite := 0
	if req.GetPage() > 1 {
		offsite = (req.GetPage() - 1) * req.GetLimit()
	}

	err = conn.Where("company_id = ? ", loginUser.CompanyID).
		Where("name LIKE ? ", "%"+req.Name+"%").
		Where("delete_dt IS NULL").
		Offset(offsite).
		Limit(req.GetLimit()).
		Find(&res).Error
	if err != nil {
		return err, cnt, res
	}

	return err, cnt, res
}

// List Property
// @Summary List Property
// @Tags Property
// @Accept json
// @Produce json
// @Param req body request.ListProperty true "payload"
// @Success      200  {object}	response.Response{payload=[]response.List}
// @Failure      500  {object}  response.Response
// @Router /property/list [post]
func (h Property) List(c echo.Context) error {
	var err error

	req := new(request.ListProperty)
	if err = c.Bind(req); err != nil {
		errorInternal(c, err)
	}

	if err = c.Validate(req); err != nil {
		return response.Error(http.StatusBadRequest, "error validation", response.ValidationError(err)).SendJSON(c)
	}

	err, list := getListProperty(req)
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusOK, "success", list).SendJSON(c)
}

func getListProperty(req *request.ListProperty) (error, []response.List) {
	var err error
	var res []response.List
	var data []model.PropertyView

	conn, closeConn := db.GetConnection()
	defer closeConn()

	if req.CompanyID == "" {
		err = conn.Where("name LIKE ? ", "%"+req.Name+"%").
			Where("description LIKE ? ", "%"+req.Description+"%").
			Where("delete_dt IS NULL").
			Find(&data).Error
	} else {
		err = conn.Where("company_id = ? ", req.CompanyID).
			Where("name LIKE ? ", "%"+req.Name+"%").
			Where("description LIKE ? ", "%"+req.Description+"%").
			Where("delete_dt IS NULL").
			Find(&data).Error
	}
	if err != nil {
		return err, res
	}

	res = make([]response.List, 0)
	for _, value := range data {
		add := response.List{
			Value: value.ID,
			Label: value.Name,
		}
		res = append(res, add)
	}

	return err, res
}

// GetById godoc
// @Tags Property
// @Summary To do get a user
// @Accept json
// @Produce json
// @Param id path string true "Property ID"
// @Success      200  {object}	response.Response{payload=response.Property}
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

	res := response.Property(property)

	return response.Success(http.StatusOK, "success", res).SendJSON(c)
}

// Create godoc
// @Tags Property
// @Summary To do create new property
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
