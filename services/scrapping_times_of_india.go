package services

import (
	"fmt"
	"github.com/tebeka/selenium"
)

func TimesOfIndia(driver selenium.WebDriver) {
	var (
		count  = 0
		length int
	)
	fmt.Println("Scrapping Times OF India...")
	//SCRAPPING NEWS HANDLE
	scrappedNews := NewsScrapper(driver, "timesofindia")

	//TOTAL LENGTH OF SCRAPPED DATA
	length = len(scrappedNews)

	//DATA INSERTION
	for _, news := range scrappedNews {
		count = InsertScrappedData(&news, count)
	}

	//COUNT OF NEW DATA INSERTED IN DATABASE
	total := length - count
	if total == 0 { //IF NO RECORDS WERE INSERTED IN DATABASE
		fmt.Println("Records Up-To-Date.")
	} else { //ELSE NUMBER OF RECORDS INSERTED IN DATABASE
		fmt.Println(total, " Record(s) Inserted.")
	}
}
