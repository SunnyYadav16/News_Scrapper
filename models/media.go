package models

import (
	"github.com/SunnyYadav16/News_Scrapper/utils"
	"gorm.io/gorm"
)

type Media struct {
	gorm.Model
	Type          string
	URL           string
	NewsHandlerID uint
}

func (media *Media) Insert() error {
	db := utils.NewDatabase()
	return db.Create(&media).Error
}
