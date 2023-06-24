package model

import (
	"time"
)

type User struct {
	ID          string     `gorm:"primaryKey"`
	RoleID      string     `gorm:"not null"`
	Email       string     `gorm:"not null"`
	Username    string     `gorm:"not null"`
	NoHp        string     `gorm:"not null"`
	Fullname    string     `gorm:"not null"`
	Passwd      string     `gorm:"not null"`
	PassVersion int        `gorm:"not null"`
	Active      bool       `gorm:"not null"`
	LastLoginDt *time.Time `gorm:"null"`
	PhotoID     string     `gorm:"not null"`
	CreateBy    string     `gorm:"not null"`
	CreateDt    time.Time  `gorm:"not null"`
	UpdateBy    string     `gorm:"not null"`
	UpdateDt    time.Time  `gorm:"not null"`
	DeleteBy    string     `gorm:"not null"`
	DeleteDt    *time.Time `gorm:"null"`
}

type Company struct {
	ID       string     `gorm:"primaryKey"`
	UserID   string     `gorm:"not null"`
	Name     string     `gorm:"not null"`
	CreateBy string     `gorm:"not null"`
	CreateDt time.Time  `gorm:"not null"`
	UpdateBy string     `gorm:"not null"`
	UpdateDt time.Time  `gorm:"not null"`
	DeleteBy string     `gorm:"not null"`
	DeleteDt *time.Time `gorm:"null"`
}

type Companysetting struct {
	ID               string `gorm:"primaryKey"`
	DefaultTimeStart int    `gorm:"not null"`
	DefaultTimeEnd   int    `gorm:"not null"`
}

type Usercompany struct {
	ID             string     `gorm:"primaryKey"`
	UserID         string     `gorm:"not null"`
	CompanyID      string     `gorm:"not null"`
	DefaultCompany bool       `gorm:"not null"`
	CreateBy       string     `gorm:"not null"`
	CreateDt       time.Time  `gorm:"not null"`
	UpdateBy       string     `gorm:"not null"`
	UpdateDt       time.Time  `gorm:"not null"`
	DeleteBy       string     `gorm:"not null"`
	DeleteDt       *time.Time `gorm:"null"`
}

type Property struct {
	ID          string     `gorm:"primaryKey"`
	CompanyID   string     `gorm:"not null"`
	Name        string     `gorm:"not null"`
	Description string     `gorm:"not null"`
	CreateBy    string     `gorm:"not null"`
	CreateDt    time.Time  `gorm:"not null"`
	UpdateBy    string     `gorm:"not null"`
	UpdateDt    time.Time  `gorm:"not null"`
	DeleteBy    string     `gorm:"not null"`
	DeleteDt    *time.Time `gorm:"null"`
}

type Item struct {
	ID          string     `gorm:"primaryKey"`
	CompanyID   string     `gorm:"not null"`
	Name        string     `gorm:"not null"`
	Description string     `gorm:"not null"`
	CreateBy    string     `gorm:"not null"`
	CreateDt    time.Time  `gorm:"not null"`
	UpdateBy    string     `gorm:"not null"`
	UpdateDt    time.Time  `gorm:"not null"`
	DeleteBy    string     `gorm:"not null"`
	DeleteDt    *time.Time `gorm:"null"`
}

type Calendar struct {
	ID         string     `gorm:"primaryKey"`
	CompanyID  string     `gorm:"not null"`
	PropertyID string     `gorm:"not null"`
	Name       string     `gorm:"not null"`
	StartDt    time.Time  `gorm:"not null"`
	EndDt      time.Time  `gorm:"not null"`
	Status     int        `gorm:"not null"`
	CreateBy   string     `gorm:"not null"`
	CreateDt   time.Time  `gorm:"not null"`
	UpdateBy   string     `gorm:"not null"`
	UpdateDt   time.Time  `gorm:"not null"`
	DeleteBy   string     `gorm:"not null"`
	DeleteDt   *time.Time `gorm:"null"`
}
