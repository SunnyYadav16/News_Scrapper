package services

import (
	"github.com/SunnyYadav16/News_Scrapper/models"
	"github.com/SunnyYadav16/News_Scrapper/utils"
)

func InitMigrations() {
	db := utils.NewDatabase()
	err := db.AutoMigrate(&models.NewsHandler{})
	CheckError("Error Migrating NewsHandler Model", err)
	err = db.AutoMigrate(&models.Media{})
	CheckError("Error Migrating Media Model", err)
	err = db.AutoMigrate(&models.UserHandle{})
	CheckError("Error Migrating UserHandles Model", err)
	err = db.AutoMigrate(&models.HashTag{})
	CheckError("Error Migrating HashTag Model", err)
}
