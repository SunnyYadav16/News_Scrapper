package services

import (
	"github.com/SunnyYadav16/News_Scrapper/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitConnection() {
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=root dbname=TwitterData port=5432 sslmode=disable TimeZone=Asia/Shanghai"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.NewsHandler{})
}
