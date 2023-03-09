package main

import (
	"encoding/json"
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/services"
	conditions "github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
	"os"
)

func main() {
	var (
		jsonResult []byte
		err        error
		driver     selenium.WebDriver
	)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			fmt.Println("Exiting program.")
			os.Exit(0)
		}
		services.CloseService(driver)
	}()

	//LOGIN TO THE TWITTER
	driver = services.TwitterLogin("SimformGolang", "Golang@Simform@123")

	//SLEEP TIMEOUT FOR PAGE LOADING AND AVOIDING CAPTCHA
	err = driver.Wait(conditions.URLContains("https://twitter.com"))
	services.CheckError("Error Loading Twitter Page", err)

	//SCRAPPING NEWS HANDLE
	scrappedNews := services.NewsScrapper(driver, "indiatoday")

	//CONVERTING SCRAPED NEWS TO JSON FORMAT FOR TESTING AND VISIBILITY
	jsonResult, err = json.MarshalIndent(scrappedNews, " ", "\t")
	services.CheckError("Cannot Convert to Json Format", err)

	//DISPLAYING SCRAPPED DATA
	fmt.Println(string(jsonResult))
}
