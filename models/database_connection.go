package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB
var err error

const DSN = "host=localhost user=postgres password=Simform@123 dbname=TwitterData port=5432 sslmode=disable TimeZone=Asia/Shanghai"

func InitConnection() {
	Db, err = gorm.Open(postgres.Open(DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println("Error connecting to database")
	}
	fmt.Println("Successfully connected to database")
	Db.AutoMigrate(&NewsHandler{}, &Media{}, &Hashtags{}, &UserHandles{})
}

func Insert(newshandler *[]NewsHandler) {
	var nw []NewsHandler
	flag := 0
	for _, val := range *newshandler {
		fmt.Println("-----------------------------")
		res := Db.First(&NewsHandler{TweetId: val.TweetId})
		if res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				fmt.Println("New Tweet id: ", val.TweetId, " inserted successfully")
				nw = append(nw, val)
				flag++
			}
		} else {
			fmt.Println("Record already present")
		}
		fmt.Println("-----------------------------")
	}
	if flag != 0 {
		i := Db.Create(nw)
		fmt.Println("Rows affected", i.RowsAffected)
	} else {
		fmt.Println("No new tweets found")
	}
}
