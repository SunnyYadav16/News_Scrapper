package main

import (
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/models"
	"github.com/SunnyYadav16/News_Scrapper/services"
	"github.com/SunnyYadav16/News_Scrapper/utils"
	conditions "github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
	"os"
	"time"
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
	)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			fmt.Println("Exiting program.")
			os.Exit(0)
		}
		utils.CloseDatabase()
		services.CloseService(driver)
		main()
	}()

	//LOGIN TO THE TWITTER
	driver = services.TwitterLogin("SimformGolang", "Golang@Simform@123")

	//SLEEP TIMEOUT FOR PAGE LOADING AND AVOIDING CAPTCHA
	err = driver.Wait(conditions.URLContains("https://twitter.com"))
	services.CheckError("Error Loading Twitter Page", err)

	ticker := time.NewTicker(150 * time.Second)
	fmt.Println("Running Program!")
	for range ticker.C {
		go services.TimesOfIndia(driver)
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
