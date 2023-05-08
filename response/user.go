package response

import (
	"github.com/jihanlugas/rental/model"
)

type User model.UsersView
type Users []model.UsersView

type Property model.PropertiesView
type Properties []model.PropertiesView

type Propertysetting model.PropertysettingsView
type Propertysettings []model.PropertysettingsView

type LoginUser struct {
	User            User            `json:"user"`
	Property        Property        `json:"property"`
	Propertysetting Propertysetting `json:"propertysetting"`
}
