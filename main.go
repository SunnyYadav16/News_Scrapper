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
		count  = 0
		length int
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

	//TOTAL LENGTH OF SCRAPPED DATA
	length = len(scrappedNews)

	//DATA INSERTION
	for _, news := range scrappedNews {
		count = insertScrappedData(&news, count)
	}

	//COUNT OF NEW DATA INSERTED IN DATABASE
	total := length - count
	if total == 0 { //IF NO RECORDS WERE INSERTED IN DATABASE
		fmt.Println("Records Up-To-Date.")
	} else { //ELSE NUMBER OF RECORDS INSERTED IN DATABASE
		fmt.Println(total, " Record(s) Inserted.")
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

// INSERTING SCRAPPED DATA
func insertScrappedData(news *models.NewsHandler, count int) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			panic("Failed To Insert Data")
		}
	}()

	var (
		check bool
		err   error
	)
	//TRY TO INSERT SCRAPPED DATA
	check, err = news.Insert()
	if check { //IF DATA NOT EXISTS IN DATABASE
		utils.PanicError("Error Inserting Scrapped Data", err)

		//TRY TO INSERT USER-HANDLE DATA
		for _, userHandle := range news.UserHandles {
			check, err = userHandle.Insert()
			if check { //IF DATA NOT EXISTS IN DATABASE
				utils.PanicError("Error Inserting Data in User Handle", err)
			}
			check, err = insertNewsHandleUserHandle(userHandle.ID, news.ID)
			if check { //IF DATA NOT EXISTS IN DATABASE
				utils.PanicError("Error Inserting In UserHandle And NewsHandler Join Table", err)
			}
		}

		//TRY TO INSERT HASH-TAG DATA
		for _, hashTag := range news.HashTags {
			check, err = hashTag.Insert()
			if check { //IF DATA NOT EXISTS IN DATABASE
				utils.PanicError("Error Inserting Data in Hash Tag", err)
			}
			check, err = insertNewsHandleHashTag(hashTag.ID, news.ID)
			if check { //IF DATA NOT EXISTS IN DATABASE
				utils.PanicError("Error Inserting In UserHandle And NewsHandler Join Table", err)
			}
		}
	} else {

		//COUNT OF SCRAPPED DATA ALREADY PRESENT IN DATABASE
		count++
	}
	return count
}

// INSERTING DATA IN NEWS-HANDLE AND USER-HANDLE JOIN TABLE
func insertNewsHandleUserHandle(userHandleID, newsID uint) (bool, error) {
	db := utils.NewDatabase()

	//CHECK IF DATA ALREADY EXISTS
	row := db.Table("newshandel_userhandle").Where("user_handle_id = ? AND news_handler_id = ?", userHandleID, newsID)
	if row.RowsAffected < 0 { //IF DATA NOT EXISTS
		return true, db.Exec("INSERT INTO  newshandel_userhandle (user_handle_id,news_handler_id) VALUES (?,?)", userHandleID, newsID).Error
	}
	return false, nil
}

// INSERTING DATA IN NEWS-HANDLE AND HASH-TAG JOIN TABLE
func insertNewsHandleHashTag(hashTagID, newsID uint) (bool, error) {
	db := utils.NewDatabase()

	//CHECK IF DATA ALREADY EXISTS
	row := db.Table("newshandel_hashtags").Where("hash_tag_id = ? AND news_handler_id = ?", hashTagID, newsID)
	if row.RowsAffected < 0 { //IF DATA NOT EXISTS
		return true, db.Exec("INSERT INTO  newshandel_hashtags (hash_tag_id,news_handler_id) VALUES (?,?)", hashTagID, newsID).Error
	}
	return false, nil
}
