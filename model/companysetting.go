package model

import (
	"github.com/jihanlugas/rental/utils"
	"gorm.io/gorm"
)

func (u *Companysetting) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = utils.GetUniqueID()
	return
}

func (u *Companysetting) BeforeUpdate(tx *gorm.DB) (err error) {
	return
}
