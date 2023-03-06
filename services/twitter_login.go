package services

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"time"
)

type FieldFinder struct {
	Name          string //NAME OF THE FIELD
	Value         string //VALUE TO BE ADDED TO THE FIELD
	SelectorQuery string //SELECTOR QUERY OF THE FIELD
	SelectorType  string //TYPE OF SELECTOR USED
}

func CheckError(msg string, err error) {
	if err != nil {
		panic(msg)
	}
}

func textBoxFindAndInsert(driver selenium.WebDriver, fieldFinder FieldFinder) {

	//Finding Text-Box
	textBox, err := driver.FindElement(fieldFinder.SelectorType, fieldFinder.SelectorQuery)
	CheckError("Error Finding "+fieldFinder.Name+" Text-Box", err)

	//Inserting TextBox
	err = textBox.SendKeys(fieldFinder.Value)
	CheckError("Error Inserting "+fieldFinder.Name, err)
}

func buttonFindAndClick(driver selenium.WebDriver, fieldFinder FieldFinder) {
	button, err := driver.FindElement(fieldFinder.SelectorType, fieldFinder.SelectorQuery)
	CheckError("Error Finding "+fieldFinder.Name+" Button", err)

	err = button.Click()
	CheckError("Error going to the "+fieldFinder.Name, err)
}

func TwitterLogin(userName, password string) {

	//INITIALISING VARIABLES
	var (
		service *selenium.Service  //SERVICE
		driver  selenium.WebDriver //DRIVER
		err     error              //ERROR
	)

	//STARTING CHROME DRIVER SERVICE
	service, err = selenium.NewChromeDriverService("./chromedriver", 4444)
	CheckError("Something Went Wrong While Creating Chrome Driver Service", err)

	//CLOSING CHROME DRIVER SERVICE
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			fmt.Println("Closing Driver and exiting program.")
			err = service.Stop()
			CheckError("Error Closing Chrome Driver Service", err)
		}
	}()

	//DEFINING DRIVER CAPABILITIES
	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"window-size=1920x1080",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"disable-gpu",
		//"--headless", // COMMENT OUT THIS LINE TO SEE THE BROWSER
	}})

	//CREATING NEW DRIVER
	driver, err = selenium.NewRemote(caps, "")
	CheckError("'Error Creating Remote Service", err)

	//REQUESTING TWITTER LOGIN PAGE
	err = driver.Get("https://twitter.com/i/flow/login?input_flow_data=%7B%22requested_variant%22%3A%22eyJsYW5nIjoiZW4ifQ%3D%3D%22%7D")
	CheckError("Error Redirecting To The Login URL", err)

	//WAITING FOR THE PAGE LOAD
	time.Sleep(5 * time.Second)

	//Find And Insert UserName
	textBoxFindAndInsert(driver, FieldFinder{
		Name:          "Username",
		Value:         userName,
		SelectorQuery: "input[type=text]",
		SelectorType:  selenium.ByCSSSelector,
	})

	//REDIRECTING TO NEXT PAGE
	buttonFindAndClick(driver, FieldFinder{
		Name:          "Next Page",
		SelectorQuery: `//*[@id="layers"]/div/div/div/div/div/div/div[2]/div[2]/div/div/div[2]/div[2]/div/div/div/div[6]`,
		SelectorType:  selenium.ByXPATH,
	})

	//LOADING PAGE WAIT TIME
	time.Sleep(1 * time.Second)

	//FIND AND INSERT PASSWORD
	textBoxFindAndInsert(driver, FieldFinder{
		Name:          "Password",
		Value:         password,
		SelectorQuery: "input[type=password]",
		SelectorType:  selenium.ByCSSSelector,
	})

	//CHECKING USER CREDENTIALS
	buttonFindAndClick(driver, FieldFinder{
		Name:          "Login Page",
		SelectorQuery: "div[role=button]",
		SelectorType:  selenium.ByCSSSelector,
	})
}
