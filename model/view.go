package model

import "time"

type UserView struct {
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

func (UserView) TableName() string {
	return "users_view"
}

type CompanyView struct {
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

func (CompanyView) TableName() string {
	return "companies_view"
}

type CompanysettingView struct {
	ID               string `json:"id"`
	DefaultTimeStart int    `json:"defaultTimeStart"`
	DefaultTimeEnd   int    `json:"defaultTimeEnd"`
}

func (CompanysettingView) TableName() string {
	return "companysettings_view"
}

type UsercompanyView struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	CompanyID string     `json:"companyId"`
	CreateBy  string     `json:"createBy"`
	CreateDt  time.Time  `json:"createDt"`
	UpdateBy  string     `json:"updateBy"`
	UpdateDt  time.Time  `json:"updateDt"`
	DeleteBy  string     `json:"deleteBy"`
	DeleteDt  *time.Time `json:"deleteDt"`
}

func (UsercompanyView) TableName() string {
	return "usercompanies_view"
}

type PropertyView struct {
	ID          string     `json:"id"`
	CompanyID   string     `json:"companyId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreateBy    string     `json:"createBy"`
	CreateDt    time.Time  `json:"createDt"`
	UpdateBy    string     `json:"updateBy"`
	UpdateDt    time.Time  `json:"updateDt"`
	DeleteBy    string     `json:"deleteBy"`
	DeleteDt    *time.Time `json:"deleteDt"`
}

func (PropertyView) TableName() string {
	return "properties_view"
}

type CalendarView struct {
	ID         string    `json:"id"`
	CompanyID  string    `json:"companyId"`
	PropertyID string    `json:"propertyId"`
	Name       string    `json:"name"`
	StartDt    time.Time `json:"startDt"`
	EndDt      time.Time `json:"endDt"`
	CreateBy   string    `json:"createBy"`
	CreateDt   time.Time `json:"createDt"`
	UpdateBy   string    `json:"updateBy"`
	UpdateDt   time.Time `json:"updateDt"`
}

func (CalendarView) TableName() string {
	return "calendars_view"
}
