package services

import (
	"github.com/SunnyYadav16/News_Scrapper/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitConnection() {

	db, err = gorm.Open(postgres.Open("host=localhost user=postgres password=root dbname=TwitterData port=5432 sslmode=disable TimeZone=Asia/Shanghai"), &gorm.Config{})
	db.Debug()
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.NewsHandler{}, &models.Media{}, &models.Hashtags{}, &models.UserHandles{})
	//db.Model(&models.NewsHandler{}).Preload("media")

	CheckError("error occurred while migrating data:", err)
}
