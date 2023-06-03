package controller

import (
	"errors"
	"github.com/jihanlugas/rental/config"
	"github.com/jihanlugas/rental/cryption"
	"github.com/jihanlugas/rental/db"
	"github.com/jihanlugas/rental/model"
	"github.com/jihanlugas/rental/request"
	"github.com/jihanlugas/rental/response"
	"github.com/jihanlugas/rental/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type User struct{}

func UserComposer() User {
	return User{}
}

func (h User) GetById(c echo.Context) error {
	var err error
	var user model.UserView

	conn, closeConn := db.GetConnection()
	defer closeConn()

	ID := c.Param("id")
	if ID == "" {
		return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
	}

	err = conn.Where("id = ? ", ID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	return response.Success(http.StatusOK, "success", user).SendJSON(c)
}

// SignIn Sign In user
// @Summary Sign in a user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param req body request.Signin true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /sign-in [post]
func (h User) SignIn(c echo.Context) error {
	var err error
	var user model.User
	var usercompany model.Usercompany

	req := new(request.Signin)
	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return response.Error(http.StatusBadRequest, "error validation", response.ValidationError(err)).SendJSON(c)
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	if utils.IsValidEmail(req.Username) {
		user.Email = req.Username
		err = conn.Where("email = ? ", user.Email).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.Error(http.StatusBadRequest, "invalid username or password", response.Payload{}).SendJSON(c)
			}
			errorInternal(c, err)
		}
	} else {
		user.Username = req.Username
		err = conn.Where("username = ? ", user.Username).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.Error(http.StatusBadRequest, "invalid username or password", response.Payload{}).SendJSON(c)
			}
			errorInternal(c, err)
		}
	}

	if !user.Active {
		return response.Error(http.StatusBadRequest, "user not active", response.Payload{}).SendJSON(c)
	}

	err = cryption.CheckAES64(req.Passwd, user.Passwd)
	if err != nil {
		return response.Error(http.StatusBadRequest, "invalid username or password", response.Payload{}).SendJSON(c)
	}

	err = conn.Where("user_id = ? ", user.ID).Where("default_company = ? ", true).First(&usercompany).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	tx := conn.Begin()

	now := time.Now()
	user.LastLoginDt = &now
	user.UpdateDt = now
	tx.Save(&user)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	expiredAt := time.Now().Add(time.Hour * time.Duration(config.AuthTokenExpiredHour))
	token, err := CreateToken(user.ID, user.RoleID, usercompany.CompanyID, user.PassVersion, expiredAt)
	if err != nil {
		return response.Error(http.StatusBadRequest, "Failed generate token", response.Payload{}).SendJSON(c)
	}

	return response.Success(http.StatusOK, "success", response.Payload{
		"token": token,
	}).SendJSON(c)
}

func (h User) SignOut(c echo.Context) error {
	return response.Success(http.StatusOK, "success", response.Payload{}).SendJSON(c)
}

func (h User) Init(c echo.Context) error {
	var err error
	var userview model.UserView
	var companyview model.CompanyView
	var companysettingview model.CompanysettingView

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Where("id = ? ", loginUser.UserID).First(&userview).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "1record not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	err = conn.Where("id = ? ", loginUser.CompanyID).First(&companyview).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "2record not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	err = conn.Where("id = ? ", loginUser.CompanyID).First(&companysettingview).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "3record not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	res := response.LoginUser{
		User:           response.User(userview),
		Company:        response.Company(companyview),
		Companysetting: response.Companysetting(companysettingview),
	}

	return response.Success(http.StatusOK, "success", res).SendJSON(c)
}

func (h User) RefreshToken(c echo.Context) error {
	var err error

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	expiredAt := time.Now().Add(time.Hour * time.Duration(config.AuthTokenExpiredHour))
	token, err := CreateToken(loginUser.UserID, loginUser.RoleID, loginUser.CompanyID, loginUser.PassVersion, expiredAt)
	if err != nil {
		return response.ErrorForce(http.StatusBadRequest, "Failed generate token", response.Payload{}).SendJSON(c)
	}

	return response.Success(http.StatusOK, "success", response.Payload{
		"token": token,
	}).SendJSON(c)
}
