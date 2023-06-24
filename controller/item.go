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

type Item struct{}

func ItemComposer() Item {
	return Item{}
}

// Page Item
// @Summary Page Item
// @Tags Item
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.PageItem true "payload"
// @Success      200  {object}	response.Response{payload=response.Pagination}
// @Failure      500  {object}  response.Response
// @Router /item/page [post]
func (h Item) Page(c echo.Context) error {
	var err error

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	req := new(request.PageItem)
	if err = c.Bind(req); err != nil {
		errorInternal(c, err)
	}

	if err = c.Validate(req); err != nil {
		return response.Error(http.StatusBadRequest, "error validation", response.ValidationError(err)).SendJSON(c)
	}

	err, cnt, list := getPageItem(req, loginUser)
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusOK, "success", response.PayloadPagination(req, list, cnt)).SendJSON(c)
}

func getPageItem(req *request.PageItem, loginUser UserLogin) (error, int64, response.PageItem) {
	var err error
	var cnt int64
	var res response.PageItem

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

// List Item
// @Summary List Item
// @Tags Item
// @Accept json
// @Produce json
// @Param req body request.ListItem true "payload"
// @Success      200  {object}	response.Response{payload=[]response.List}
// @Failure      500  {object}  response.Response
// @Router /item/list [post]
func (h Item) List(c echo.Context) error {
	var err error

	req := new(request.ListItem)
	if err = c.Bind(req); err != nil {
		errorInternal(c, err)
	}

	if err = c.Validate(req); err != nil {
		return response.Error(http.StatusBadRequest, "error validation", response.ValidationError(err)).SendJSON(c)
	}

	err, list := getListItem(req)
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusOK, "success", list).SendJSON(c)
}

func getListItem(req *request.ListItem) (error, []response.List) {
	var err error
	var res []response.List
	var data []model.ItemView

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
// @Tags Item
// @Summary To do get a user
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success      200  {object}	response.Response{payload=response.Item}
// @Failure      500  {object}  response.Response
// @Router /item/{id} [get]
func (h Item) GetById(c echo.Context) error {
	var err error
	var item model.ItemView

	conn, closeConn := db.GetConnection()
	defer closeConn()

	ID := c.Param("id")
	if ID == "" {
		return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
	}

	err = conn.Where("id = ? ", ID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	res := response.Item(item)

	return response.Success(http.StatusOK, "success", res).SendJSON(c)
}

// Create godoc
// @Tags Item
// @Summary To do create new item
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.CreateItem true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /item [post]
func (h Item) Create(c echo.Context) error {
	var err error
	var item model.Item

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	req := new(request.CreateItem)
	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return response.Error(http.StatusBadRequest, "error validation", response.ValidationError(err)).SendJSON(c)
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	item = model.Item{
		CompanyID:   loginUser.CompanyID,
		Name:        req.Name,
		Description: req.Description,
		CreateBy:    loginUser.UserID,
		UpdateBy:    loginUser.UserID,
	}

	tx.Save(&item)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusCreated, "success", response.Payload{}).SendJSON(c)
}

// Update godoc
// @Tags Item
// @Summary To do update a item
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param req body request.UpdateItem true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /item/{id} [put]
func (h Item) Update(c echo.Context) error {
	var err error
	var item model.Item

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	ID := c.Param("id")
	if ID == "" {
		return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
	}

	req := new(request.UpdateItem)
	if err = c.Bind(req); err != nil {
		return err
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Where("id = ? ", ID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	tx := conn.Begin()

	item.Name = req.Name
	item.Description = req.Description
	item.UpdateBy = loginUser.UserID

	tx.Save(&item)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusAccepted, "success", response.Payload{}).SendJSON(c)
}

// Delete Item
// @Summary Delete Item
// @Tags Item
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /item/{id} [delete]
func (h Item) Delete(c echo.Context) error {
	var err error
	var item model.Item

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

	err = conn.Where("id = ? ", ID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	tx := conn.Begin()

	now := time.Now()

	item.DeleteBy = loginUser.UserID
	item.DeleteDt = &now
	tx.Save(&item)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	return response.Success(http.StatusAccepted, "success", response.Payload{}).SendJSON(c)
}
