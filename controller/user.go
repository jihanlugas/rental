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

func (h User) SignIn(c echo.Context) error {
	var err error
	var user model.User
	var userproperty model.Userproperties

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

	err = conn.Where("user_id = ? ", user.ID).Where("default_property = ? ", true).First(&userproperty).Error
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
	token, err := getLoginToken(user.ID, user.RoleID, userproperty.PropertyID, user.PassVersion, expiredAt)
	if err != nil {
		return response.Error(http.StatusBadRequest, "generate token failed", response.Payload{}).SendJSON(c)
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
	var userview model.UsersView
	var propertyview model.PropertiesView
	var propertysettingview model.PropertysettingsView

	loginUser, err := getUserLoginInfo(c)
	if err != nil {
		errorInternal(c, err)
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Where("id = ? ", loginUser.UserID).First(&userview).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "record not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	err = conn.Where("id = ? ", loginUser.PropertyID).First(&propertyview).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "record not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	err = conn.Where("id = ? ", loginUser.PropertyID).First(&propertysettingview).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "record not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	res := response.LoginUser{
		User:            response.User(userview),
		Property:        response.Property(propertyview),
		Propertysetting: response.Propertysetting(propertysettingview),
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
	token, err := getLoginToken(loginUser.UserID, loginUser.RoleID, loginUser.PropertyID, loginUser.PassVersion, expiredAt)
	if err != nil {
		return response.ErrorForce(http.StatusBadRequest, "generate token failed", response.Payload{}).SendJSON(c)
	}

	return response.Success(http.StatusOK, "success", response.Payload{
		"token": token,
	}).SendJSON(c)
}
