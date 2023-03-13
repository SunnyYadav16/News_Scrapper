package services

import (
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/models"
	"sort"
)

func insertData(TweetData []models.NewsHandler) {
	filteredData := checkDuplicate(TweetData)

	result := db.Create(&filteredData)
	if result.Error != nil {
		CheckError("error occurred while inserting data:", result.Error)
	} else {
		fmt.Println("data inserted successfully")
	}
}

func checkDuplicate(TweetData []models.NewsHandler) (filteredData []models.NewsHandler) {
	//fetch data
	sort.Slice(TweetData, func(i, j int) bool {
		return TweetData[i].Timestamp.Before(TweetData[j].Timestamp)
	})

	var result []models.NewsHandler
	err := db.Model(&models.NewsHandler{}).Find(&result).Where("timestamp >= ?", TweetData[0].Timestamp)
	if err.Error != nil {
		CheckError("error occured while fetching data:", err.Error)
	} else {
		for _, singleTweetData := range TweetData {

			isPresent := false
			for _, object := range result {

				if singleTweetData.TweetId == object.TweetId {
					fmt.Println("Tweet Already Present, skipping...", object.TweetId)
					isPresent = true
				}
			}
			if !isPresent {
				filteredData = append(filteredData, singleTweetData)
			}
		}
	}
	return
}
