package model

import "time"

type UsersView struct {
	ID          string     `json:"id"`
	RoleID      string     `json:"role_id"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	NoHp        string     `json:"noHp"`
	Fullname    string     `json:"fullname"`
	Passwd      string     `json:"-"`
	PassVersion int        `json:"passVersion"`
	Active      bool       `json:"active"`
	LastLoginDt *time.Time `json:"lastLoginDt"`
	PhotoID     string     `json:"photoId"`
	PhotoUrl    string     `json:"photoUrl"`
	CreateBy    string     `json:"createBy"`
	CreateDt    time.Time  `json:"createDt"`
	UpdateBy    string     `json:"updateBy"`
	UpdateDt    time.Time  `json:"updateDt"`
	DeleteBy    string     `json:"deleteBy"`
	DeleteDt    *time.Time `json:"deleteDt"`
}

type PropertiesView struct {
	ID       string     `json:"id"`
	UserID   string     `json:"userId"`
	Name     string     `json:"name"`
	CreateBy string     `json:"createBy"`
	CreateDt time.Time  `json:"createDt"`
	UpdateBy string     `json:"updateBy"`
	UpdateDt time.Time  `json:"updateDt"`
	DeleteBy string     `json:"deleteBy"`
	DeleteDt *time.Time `json:"deleteDt"`
}

type PropertysettingsView struct {
	ID               string `json:"id"`
	DefaultTimeStart string `json:"userId"`
	DefaultTimeEnd   string `json:"name"`
}

type UserpropertiesView struct {
	ID         string     `json:"id"`
	UserID     string     `json:"userId"`
	PropertyID string     `json:"propertyId"`
	CreateBy   string     `json:"createBy"`
	CreateDt   time.Time  `json:"createDt"`
	UpdateBy   string     `json:"updateBy"`
	UpdateDt   time.Time  `json:"updateDt"`
	DeleteBy   string     `json:"deleteBy"`
	DeleteDt   *time.Time `json:"deleteDt"`
}

type CalendarView struct {
	ID         string     `json:"id"`
	PropertyID string     `json:"propertyId"`
	GroupID    string     `json:"groupId"`
	Name       string     `json:"name"`
	StartDt    time.Time  `json:"startDt"`
	EndDt      time.Time  `json:"endDt"`
	CreateBy   string     `json:"createBy"`
	CreateDt   time.Time  `json:"createDt"`
	UpdateBy   string     `json:"updateBy"`
	UpdateDt   time.Time  `json:"updateDt"`
	DeleteBy   string     `json:"deleteBy"`
	DeleteDt   *time.Time `json:"deleteDt"`
}

func (UsersView) TableName() string {
	return "users_view"
}

func (PropertiesView) TableName() string {
	return "properties_view"
}

func (PropertysettingsView) TableName() string {
	return "propertysettings_view"
}

func (UserpropertiesView) TableName() string {
	return "userproperties_view"
}

func (CalendarView) TableName() string {
	return "calendar_view"
}
