package main

import (
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/models"
	"github.com/SunnyYadav16/News_Scrapper/services"
	"github.com/tebeka/selenium"
	"time"
)

func main() {
	var driver selenium.WebDriver
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			fmt.Println("Exiting program.")
		}
		services.CloseService(driver)
	}()
	models.InitConnection()
	driver = services.TwitterLogin("SimformGolang", "Golang@Simform@123")
	time.Sleep(10 * time.Second)
	url, err := driver.CurrentURL()
	services.CheckError("Error Getting Current URL", err)
	fmt.Println(url)
	services.TwitterLandingPage(driver)
	newshandler := services.NewsScrapper(driver)
	models.Insert(&newshandler)
	services.WriteIntoJSONFILE(&newshandler)
}
