package models

import (
	"github.com/SunnyYadav16/News_Scrapper/utils"
	"gorm.io/gorm"
)

type HashTag struct {
	gorm.Model
	TagName      string
	NewsHandlers []*NewsHandler `gorm:"<-:false;many2many:newshandel_hashtags;constraint:OnDelete:CASCADE;"`
}

func (hashTag *HashTag) Insert() (bool, error) {
	db := utils.NewDatabase()
	check := hashTag.Find()
	if !check {
		return true, db.Create(&hashTag).Error
	}
	return false, nil
}

func (hashTag *HashTag) Find() bool {
	db := utils.NewDatabase()
	num := db.Where("tag_name", hashTag.TagName).First(&hashTag)
	if num.RowsAffected > 0 {
		return true
	}
	return false
}
