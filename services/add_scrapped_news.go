package services

import (
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/models"
	"github.com/SunnyYadav16/News_Scrapper/utils"
)

// InsertScrappedData : INSERTING SCRAPPED DATA
func InsertScrappedData(news *models.NewsHandler, count int) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			panic("Failed To Insert Data")
		}
	}()

	var (
		check bool
		err   error
	)
	//TRY TO INSERT SCRAPPED DATA
	check, err = news.Insert()
	if check { //IF DATA NOT EXISTS IN DATABASE
		utils.PanicError("Error Inserting Scrapped Data", err)

		//TRY TO INSERT USER-HANDLE DATA
		for _, userHandle := range news.UserHandles {
			check, err = userHandle.Insert()
			if check { //IF DATA NOT EXISTS IN DATABASE
				utils.PanicError("Error Inserting Data in User Handle", err)
				check, err = insertNewsHandleUserHandle(userHandle.ID, news.ID)
				if check { //IF DATA NOT EXISTS IN DATABASE
					utils.PanicError("Error Inserting In UserHandle And NewsHandler Join Table", err)
				}
			}
		}

		//TRY TO INSERT HASH-TAG DATA
		for _, hashTag := range news.HashTags {
			check, err = hashTag.Insert()
			if check { //IF DATA NOT EXISTS IN DATABASE
				utils.PanicError("Error Inserting Data in Hash Tag", err)
			}
			check, err = insertNewsHandleHashTag(hashTag.ID, news.ID)
			if check { //IF DATA NOT EXISTS IN DATABASE
				utils.PanicError("Error Inserting In UserHandle And NewsHandler Join Table", err)
			}
		}
	} else {

		//COUNT OF SCRAPPED DATA ALREADY PRESENT IN DATABASE
		count++
	}
	return count
}

// INSERTING DATA IN NEWS-HANDLE AND USER-HANDLE JOIN TABLE
func insertNewsHandleUserHandle(userHandleID, newsID uint) (bool, error) {
	db := utils.NewDatabase()
	//CHECK IF DATA ALREADY EXISTS
	row := db.Table("newshandel_userhandle").Where("user_handle_id = ? AND news_handler_id = ?", userHandleID, newsID)
	if row.RowsAffected < 1 { //IF DATA NOT EXISTS
		return true, db.Exec("INSERT INTO  newshandel_userhandle (user_handle_id,news_handler_id) VALUES (?,?)", userHandleID, newsID).Error
	}
	return false, nil
}

// INSERTING DATA IN NEWS-HANDLE AND HASH-TAG JOIN TABLE
func insertNewsHandleHashTag(hashTagID, newsID uint) (bool, error) {
	db := utils.NewDatabase()
	//CHECK IF DATA ALREADY EXISTS
	row := db.Table("newshandel_hashtags").Where("hash_tag_id = ? AND news_handler_id = ?", hashTagID, newsID)
	if row.RowsAffected < 1 { //IF DATA NOT EXISTS
		return true, db.Exec("INSERT INTO  newshandel_hashtags (hash_tag_id,news_handler_id) VALUES (?,?)", hashTagID, newsID).Error
	}
	return false, nil
}
