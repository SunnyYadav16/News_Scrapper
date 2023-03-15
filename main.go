package main

import (
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/scrapper_utils"
	"github.com/SunnyYadav16/News_Scrapper/services"
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
		services.InitConnection()
	}()
}
func main() {
	var (
		username, password string
		driver             selenium.WebDriver
	)
	username, password = services.CheckCredentials()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			fmt.Println("Exiting program.")
		}
		services.CloseService(driver)
		main()
	}()

	driver = services.TwitterLogin(username, password)
	time.Sleep(10 * time.Second)
	url, err := driver.CurrentURL()
	services.CheckError("Error Getting Current URL", err)
	fmt.Println(url)

	//SCRAPPING DATA EVERY 2 MINUTES
	/*ticker := time.NewTicker(2 * time.Minute)
	fmt.Println("Running program")
	for range ticker.C {
		go func() {

		}()
	}*/
	scrapper_utils.TwitterLandingPage(driver)
	newshandler := scrapper_utils.NewsScrapper(driver)
	services.Insert(&newshandler)
	scrapper_utils.WriteIntoJSONFILE(&newshandler)
}
