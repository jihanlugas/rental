package response

import (
	"github.com/jihanlugas/rental/model"
)

type User model.UserView
type Users []model.UserView

type Company model.CompanyView
type Companies []model.CompanyView

type Companysetting model.CompanysettingView
type Companysettings []model.CompanysettingView

type LoginUser struct {
	User           User           `json:"user"`
	Company        Company        `json:"company"`
	Companysetting Companysetting `json:"companysetting"`
}
