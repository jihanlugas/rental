package model

import (
	"github.com/jihanlugas/rental/utils"
	"gorm.io/gorm"
	"time"
)

func (u *Item) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()

	u.ID = utils.GetUniqueID()
	u.CreateDt = now
	u.UpdateDt = now
	return
}

func (u *Item) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.DeleteDt == nil {
		now := time.Now()
		u.UpdateDt = now
	}
	return
}
