package router

import (
	"encoding/json"
	"fmt"
	"github.com/jihanlugas/rental/config"
	"github.com/jihanlugas/rental/constant"
	"github.com/jihanlugas/rental/controller"
	"github.com/jihanlugas/rental/cryption"
	"github.com/jihanlugas/rental/db"
	"github.com/jihanlugas/rental/model"
	"github.com/jihanlugas/rental/response"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Init() *echo.Echo {
	router := websiteRouter()
	checkToken := checkTokenMiddleware()
	//router.Use(middleware.Logger())

	calendaraController := controller.CalendarComposer()
	userController := controller.UserComposer()

	router.GET("/", controller.Ping)

	router.POST("/sign-in", userController.SignIn)
	router.GET("/sign-out", userController.SignOut)
	router.GET("/refresh-token", userController.RefreshToken, checkToken)

	router.GET("/init", userController.Init, checkToken)

	calendar := router.Group("/calendar")
	calendar.GET("", calendaraController.Tes, checkToken)
	calendar.GET("/ws", calendaraController.WsCalendar)

	return router
}

func httpErrorHandler(err error, c echo.Context) {
	var errorResponse *response.Response
	code := http.StatusInternalServerError
	switch e := err.(type) {
	case *echo.HTTPError:
		// Handle pada saat URL yang di request tidak ada. atau ada kesalahan server.
		code = e.Code
		errorResponse = &response.Response{
			Status:  false,
			Message: fmt.Sprintf("%v", e.Message),
			Payload: map[string]interface{}{},
			Code:    code,
		}
	case *response.Response:
		errorResponse = e
	default:
		// Handle error dari panic
		code = http.StatusInternalServerError
		if config.Debug {
			errorResponse = &response.Response{
				Status:  false,
				Message: err.Error(),
				Payload: map[string]interface{}{},
				Code:    http.StatusInternalServerError,
			}
		} else {
			errorResponse = &response.Response{
				Status:  false,
				Message: "Internal server error",
				Payload: map[string]interface{}{},
				Code:    http.StatusInternalServerError,
			}
		}
	}

	js, err := json.Marshal(errorResponse)
	if err == nil {
		_ = c.Blob(code, echo.MIMEApplicationJSONCharsetUTF8, js)
	} else {
		b := []byte("{error: true, message: \"unresolved error\"}")
		_ = c.Blob(code, echo.MIMEApplicationJSONCharsetUTF8, b)
	}
}

func checkTokenMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			token := c.Request().Header.Get(config.HeaderAuthName)

			tokenPayload, err := cryption.DecryptAES64(token)
			if err != nil {
				return response.ErrorForce(http.StatusUnauthorized, "Unauthorized", response.Payload{}).SendJSON(c)
			}

			data := strings.Split(tokenPayload, "$$")
			if len(data) != 5 {
				return response.ErrorForce(http.StatusUnauthorized, "Unauthorized.", response.Payload{}).SendJSON(c)
			}

			expiredUnix, err := strconv.ParseInt(data[4], 10, 64)
			if err != nil {
				return response.ErrorForce(http.StatusUnauthorized, "Token Expired.", response.Payload{}).SendJSON(c)
			}

			expiredAt := time.Unix(expiredUnix, 0)
			now := time.Now()
			if now.After(expiredAt) {
				return response.ErrorForce(http.StatusUnauthorized, "Token Expired..", response.Payload{}).SendJSON(c)
			}

			intVar, err := strconv.Atoi(data[3])
			if err != nil {
				return response.ErrorForce(http.StatusUnauthorized, "Token Expired...", response.Payload{}).SendJSON(c)
			}

			userLogin := controller.UserLogin{
				UserID:      data[0],
				RoleID:      data[1],
				CompanyID:   data[2],
				PassVersion: intVar,
			}

			conn, closeConn := db.GetConnection()
			defer closeConn()

			var user model.User
			err = conn.Where("id = ? ", userLogin.UserID).First(&user).Error
			if err != nil {
				return response.ErrorForce(http.StatusUnauthorized, "Token Expired!", response.Payload{}).SendJSON(c)
			}

			if user.PassVersion != userLogin.PassVersion {
				return response.ErrorForce(http.StatusUnauthorized, "Token Expired~", response.Payload{}).SendJSON(c)
			}

			c.Set(constant.TokenUserContext, userLogin)
			return next(c)
		}
	}
}
