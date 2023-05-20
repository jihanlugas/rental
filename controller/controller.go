package controller

import (
	"fmt"
	"github.com/jihanlugas/rental/constant"
	"github.com/jihanlugas/rental/cryption"
	"github.com/jihanlugas/rental/response"
	"github.com/jihanlugas/rental/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

var Validate *validator.CustomValidator

type UserLogin struct {
	UserID      string
	RoleID      string
	CompanyID   string
	PassVersion int
}

func init() {
	Validate = validator.NewValidator()
}

// Ping godoc
// @Summary      Ping
// @Tags         Ping
// @Accept       json
// @Produce      json
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router       / [get]
func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "ようこそ、美しい世界へ")
}

func errorInternal(c echo.Context, err error) {
	//log.System.Error().Err(err).Str("Host", c.Request().Host).Str("Path", c.Path()).Send()
	panic(err)
}

func getLoginToken(userID string, roleID, companyID string, passVersion int, expiredAt time.Time) (string, error) {
	expiredUnix := expiredAt.Unix()

	token := fmt.Sprintf("%s$$%s$$%s$$%d$$%d", userID, roleID, companyID, passVersion, expiredUnix)

	return cryption.EncryptAES64(token)
}

func getUserLoginInfo(c echo.Context) (UserLogin, error) {
	if u, ok := c.Get(constant.TokenUserContext).(UserLogin); ok {
		return u, nil
	} else {
		return UserLogin{}, response.ErrorForce(http.StatusUnauthorized, "Unauthorized.", response.Payload{})
	}
}
