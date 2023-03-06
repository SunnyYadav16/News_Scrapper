package main

import (
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/services"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			fmt.Println("Exiting program.")
		}
	}()
	services.TwitterLogin("SimformGolang", "Golang@Simform@123")
}
