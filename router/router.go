package router

import (
	"encoding/json"
	"fmt"
	"github.com/jihanlugas/rental/config"
	"github.com/jihanlugas/rental/constant"
	"github.com/jihanlugas/rental/controller"
	"github.com/jihanlugas/rental/db"
	"github.com/jihanlugas/rental/model"
	"github.com/jihanlugas/rental/response"
	"github.com/labstack/echo/v4"
	"net/http"

	_ "github.com/jihanlugas/rental/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init() *echo.Echo {
	router := websiteRouter()
	checkToken := checkTokenMiddleware()
	//router.Use(middleware.Logger())

	calendarController := controller.CalendarComposer()
	userController := controller.UserComposer()
	companyController := controller.CompanyComposer()
	propertyController := controller.PropertyComposer()
	itemController := controller.ItemComposer()

	router.GET("/swg/*", echoSwagger.WrapHandler)

	router.GET("/", controller.Ping)
	router.POST("/sign-in", userController.SignIn)
	router.GET("/sign-out", userController.SignOut)
	router.GET("/refresh-token", userController.RefreshToken, checkToken)
	router.GET("/init", userController.Init, checkToken)

	user := router.Group("/user")
	user.GET("/:id", userController.GetById)
	user.POST("/change-password", userController.ChangePassword, checkToken)

	company := router.Group("/company")
	company.GET("/:id", companyController.GetById)
	company.POST("", companyController.Create, checkToken)
	company.PUT("/:id", companyController.Update, checkToken)
	company.DELETE("/:id", companyController.Delete, checkToken)

	calendar := router.Group("/calendar")
	calendar.GET("/:id", calendarController.GetById)
	calendar.GET("/detail/:id", calendarController.GetDetailById, checkToken)
	calendar.POST("", calendarController.Create, checkToken)
	calendar.PUT("/:id", calendarController.Update, checkToken)
	calendar.DELETE("/:id", calendarController.Delete, checkToken)
	calendar.GET("/ws", calendarController.WsCalendar)
	calendar.POST("/timeline", calendarController.Timeline, checkToken)

	property := router.Group("/property")
	property.GET("/:id", propertyController.GetById)
	property.POST("/page", propertyController.Page, checkToken)
	property.POST("/list", propertyController.List)
	property.POST("", propertyController.Create, checkToken)
	property.PUT("/:id", propertyController.Update, checkToken)
	property.DELETE("/:id", propertyController.Delete, checkToken)

	item := router.Group("/item")
	item.GET("/:id", itemController.GetById)
	item.POST("/page", itemController.Page, checkToken)
	item.POST("/list", itemController.List)
	item.POST("", itemController.Create, checkToken)
	item.PUT("/:id", itemController.Update, checkToken)
	item.DELETE("/:id", itemController.Delete, checkToken)

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

			userLogin, err := controller.ExtractClaims(c.Request().Header.Get(config.HeaderAuthName))
			if err != nil {
				return response.ErrorForce(http.StatusUnauthorized, err.Error(), response.Payload{}).SendJSON(c)
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
