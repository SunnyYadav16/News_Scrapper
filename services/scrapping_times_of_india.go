package services

import (
	"fmt"
	"github.com/tebeka/selenium"
	"time"
)

func TimesOfIndia(driver selenium.WebDriver) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			panic("Re-Starting the Execution")
		}
	}()
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
		fmt.Println("Records Up-To-Date. ", time.Now())
	} else { //ELSE NUMBER OF RECORDS INSERTED IN DATABASE
		fmt.Println(total, " Record(s) Inserted. ", time.Now())
	}
}
