package controller

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/jihanlugas/rental/config"
	"github.com/jihanlugas/rental/constant"
	"github.com/jihanlugas/rental/response"
	"github.com/jihanlugas/rental/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
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
// @Success      200  {string} interface{} "ようこそ、美しい世界へ"
// @Router       / [get]
func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "ようこそ、美しい世界へ")
}

func errorInternal(c echo.Context, err error) {
	//log.System.Error().Err(err).Str("Host", c.Request().Host).Str("Path", c.Path()).Send()
	panic(err)
}

func CreateToken(userID string, roleID, companyID string, passVersion int, expiredAt time.Time) (string, error) {
	var err error

	now := time.Now()
	expiredUnix := expiredAt.Unix()
	subject := fmt.Sprintf("%s$$%s$$%s$$%d$$%d", userID, roleID, companyID, passVersion, expiredUnix)
	jwtKey := []byte(config.JwtSecretKey)
	claims := jwt.MapClaims{
		"iss": "Services",
		"sub": subject,
		"iat": now.Unix(),
		"exp": expiredUnix,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ExtractClaims(header string) (UserLogin, error) {
	var err error
	var userlogin UserLogin

	if header == "" {
		err = errors.New("Unauthorized")
		return userlogin, err
	}

	token := header[(len(constant.BearerSchema) + 1):]
	claims, err := parseToken(token)
	if err != nil {
		return userlogin, err
	}

	content := claims["sub"].(string)
	contentData := strings.Split(content, "$$")
	if len(contentData) != 5 {
		err = errors.New("Unauthorized")
		return userlogin, err
	}

	expiredUnix, err := strconv.ParseInt(contentData[4], 10, 64)
	if err != nil {
		return userlogin, err
	}

	expiredAt := time.Unix(expiredUnix, 0)
	now := time.Now()
	if now.After(expiredAt) {
		err = errors.New("Token Expired")
		return userlogin, err
	}

	passVersion, err := strconv.Atoi(contentData[3])
	if err != nil {
		return userlogin, err
	}

	userlogin = UserLogin{
		UserID:      contentData[0],
		RoleID:      contentData[1],
		CompanyID:   contentData[2],
		PassVersion: passVersion,
	}

	return userlogin, err
}

func parseToken(token string) (jwt.MapClaims, error) {

	jwtKey := []byte(config.JwtSecretKey)
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return jwtKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims, nil
}

//func getLoginToken(userID string, roleID, companyID string, passVersion int, expiredAt time.Time) (string, error) {
//	expiredUnix := expiredAt.Unix()
//
//	token := fmt.Sprintf("%s$$%s$$%s$$%d$$%d", userID, roleID, companyID, passVersion, expiredUnix)
//
//	return cryption.EncryptAES64(token)
//}

func getUserLoginInfo(c echo.Context) (UserLogin, error) {
	if u, ok := c.Get(constant.TokenUserContext).(UserLogin); ok {
		return u, nil
	} else {
		return UserLogin{}, response.ErrorForce(http.StatusUnauthorized, "Unauthorized.", response.Payload{})
	}
}
