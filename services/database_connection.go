package services

import (
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var Db *gorm.DB
var err error

func InitConnection() {
	err = godotenv.Load(".env")
	Db, err = gorm.Open(postgres.Open(os.Getenv("CONNECTION")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println("Error connecting to database")
	}
	fmt.Println("Successfully connected to database")
	Db.AutoMigrate(&models.NewsHandler{}, &models.Media{}, &models.Hashtags{}, &models.UserHandles{})
}

func Insert(newshandler *[]models.NewsHandler) {
	var nw []models.NewsHandler
	flag := 0
	for _, val := range *newshandler {
		fmt.Println("-----------------------------")
		res := Db.First(&models.NewsHandler{TweetId: val.TweetId})
		if res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				fmt.Println("New Tweet id: ", val.TweetId, " inserted successfully")
				nw = append(nw, val)
				flag++
			} else {
				CheckError("Error inserting record", err)
			}
		} else {
			fmt.Println("Record already present")
		}
		fmt.Println("-----------------------------")
	}
	if flag != 0 {
		i := Db.Create(nw)
		fmt.Println("Total no. of new rows inserted: ", i.RowsAffected)
	} else {
		fmt.Println("No new tweets found")
	}
}
