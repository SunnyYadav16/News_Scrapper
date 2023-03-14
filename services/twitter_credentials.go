package services

import (
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/utils"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

func ValidCredentials() (username, password string) {
	var (
		myEnv map[string]string
		err   error
	)
	myEnv, err = godotenv.Read()
	utils.PanicError("Error Reading .env file", err)
	username = myEnv["TUSER"]
	password = myEnv["TPASSWORD"]

	//CHECKING FOR CREDENTIALS
	if strings.Compare(username, "SimformGolang") != 0 && strings.Compare(password, "Golang@Simform@123") != 0 {
		fmt.Println("Invalid Tweeter Credentials. Please Try Different Credentials.")
		os.Exit(0)
	}
	return username, password
}
