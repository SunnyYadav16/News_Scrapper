package main

import (
	"fmt"
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
			fmt.Println("Migration Failed Exiting program.")
			os.Exit(0)
		}
	}()
	services.InitMigrations()
}

func main() {
	var (
		err                error
		driver             selenium.WebDriver
		username, password string
	)
	username, password = services.ValidCredentials()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			fmt.Println("Exiting program.")
			os.Exit(0)
		}
		utils.CloseDatabase()
		services.CloseService(driver)
		fmt.Printf("\n\n-------------------\n\n")
		main()
	}()

	//LOGIN TO THE TWITTER
	driver = services.TwitterLogin(username, password)

	//SLEEP TIMEOUT FOR PAGE LOADING AND AVOIDING CAPTCHA
	err = driver.Wait(conditions.URLContains("https://twitter.com"))
	services.CheckError("Error Loading Twitter Page", err)

	//SCRAPPING DATA EVERY 2 MINUTES
	ticker := time.NewTicker(2 * time.Minute)
	fmt.Println("Running Program!")
	for range ticker.C {

		//ALL SCRAPPING FUNCTIONS WILL RUN HERE
		go services.TimesOfIndia(driver)
	}
}
