package main

import (
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/models"
	"github.com/SunnyYadav16/News_Scrapper/services"
	"github.com/SunnyYadav16/News_Scrapper/utils"
	conditions "github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
	"os"
)

func init() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			fmt.Println("Exiting program.")
			os.Exit(0)
		}
	}()
	db := utils.NewDatabase()
	err := db.AutoMigrate(&models.NewsHandler{})
	services.CheckError("Error Migrating NewsHandler Model", err)
	err = db.AutoMigrate(&models.Media{})
	services.CheckError("Error Migrating Media Model", err)
	err = db.AutoMigrate(&models.UserHandle{})
	services.CheckError("Error Migrating UserHandles Model", err)
	err = db.AutoMigrate(&models.HashTag{})
	services.CheckError("Error Migrating HashTag Model", err)
}

func main() {
	var (
		//jsonResult []byte
		err    error
		driver selenium.WebDriver
		//check      bool
	)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			fmt.Println("Exiting program.")
			os.Exit(0)
		}
		utils.CloseDatabase()
		services.CloseService(driver)
	}()

	//LOGIN TO THE TWITTER
	driver = services.TwitterLogin("SimformGolang", "Golang@Simform@123")

	//SLEEP TIMEOUT FOR PAGE LOADING AND AVOIDING CAPTCHA
	err = driver.Wait(conditions.URLContains("https://twitter.com"))
	services.CheckError("Error Loading Twitter Page", err)

	//SCRAPPING NEWS HANDLE
	scrappedNews := services.NewsScrapper(driver, "timesofindia")

	for _, news := range scrappedNews {
		err = news.Insert()
		utils.PanicError("Error Inserting Scrapped Data", err)
		for _, userHandle := range news.UserHandles {
			err = userHandle.Insert()
			utils.PanicError("Error Inserting Data in User Handle", err)
			err = insertNewsHandleUserHandle(userHandle.ID, news.ID)
			utils.PanicError("Error Inserting In UserHandle And NewsHandler Join Table", err)
		}
		for _, hashTag := range news.HashTags {
			err = hashTag.Insert()
			utils.PanicError("Error Inserting Data in Hash Tag", err)
			err = insertNewsHandleHashTag(hashTag.ID, news.ID)
			utils.PanicError("Error Inserting In UserHandle And NewsHandler Join Table", err)
		}
	}
	//scrappedNews = []models.NewsHandler{}
	//scrappedNews, err = models.All()
	//services.CheckError("Error Getting All Models", err)
	////CONVERTING SCRAPED NEWS TO JSON FORMAT FOR TESTING AND VISIBILITY
	//jsonResult, err = json.MarshalIndent(scrappedNews, " ", "\t")
	//services.CheckError("Cannot Convert to Json Format", err)
	//
	////DISPLAYING SCRAPPED DATA
	//fmt.Println(string(jsonResult))
}

func insertNewsHandleUserHandle(userHandleID, newsID uint) error {
	db := utils.NewDatabase()
	row := db.Table("newshandel_userhandle").Where("user_handle_id = ? AND news_handler_id = ?", userHandleID, newsID)
	if row.RowsAffected < 0 {
		return db.Exec("INSERT INTO  newshandel_userhandle (user_handle_id,news_handler_id) VALUES (?,?)", userHandleID, newsID).Error
	}
	return nil
}

func insertNewsHandleHashTag(hashTagID, newsID uint) error {
	db := utils.NewDatabase()
	row := db.Table("newshandel_hashtags").Where("hash_tag_id = ? AND news_handler_id = ?", hashTagID, newsID)
	if row.RowsAffected < 0 {
		return db.Exec("INSERT INTO  newshandel_hashtags (hash_tag_id,news_handler_id) VALUES (?,?)", hashTagID, newsID).Error
	}
	return nil
}
