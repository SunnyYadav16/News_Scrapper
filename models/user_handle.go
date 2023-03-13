package models

import (
	"github.com/SunnyYadav16/News_Scrapper/utils"
	"gorm.io/gorm"
)

type UserHandle struct {
	gorm.Model
	Name         string
	NewsHandlers []*NewsHandler `gorm:"<-:false;many2many:newshandel_userhandle;constraint:OnDelete:CASCADE;"`
}

func (userHandle *UserHandle) Insert() (bool, error) {
	db := utils.NewDatabase()
	check := userHandle.Find()
	if !check {
		return true, db.Create(&userHandle).Error
	}
	return false, nil
}

func (userHandle *UserHandle) Find() bool {
	db := utils.NewDatabase()
	num := db.Where("name", userHandle.Name).First(&userHandle)
	if num.RowsAffected > 0 {
		return true
	}
	return false
}
