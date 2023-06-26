package model

import (
	"github.com/jihanlugas/rental/utils"
	"gorm.io/gorm"
	"time"
)

func (u *Calendaritem) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()

	u.ID = utils.GetUniqueID()
	u.CreateDt = now
	u.UpdateDt = now
	return
}

func (u *Calendaritem) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.DeleteDt == nil {
		now := time.Now()
		u.UpdateDt = now
	}
	return
}
