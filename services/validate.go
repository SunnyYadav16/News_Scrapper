package services

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

// Function for validating twitter credentials
func CheckCredentials() (username, password string) {
	err := godotenv.Load(".env")
	CheckError("Error reading .env file", err)
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")
	if strings.Compare(username, "SimformGolang") != 0 && strings.Compare(password, "Golang@Simform@123") != 0 {
		fmt.Println("Invalid twitter credentials. Please try different credentials")
		os.Exit(0)
	}
	return username, password
}
