package models

import (
	"github.com/SunnyYadav16/News_Scrapper/utils"
	"gorm.io/gorm"
	"time"
)

type NewsHandler struct {
	gorm.Model
	TweetId      string
	ChannelName  string
	TweetContent string
	Media        []Media
	HashTags     []*HashTag    `gorm:"<-:false;many2many:newshandel_hashtags;"`
	UserHandles  []*UserHandle `gorm:"<-:false;many2many:newshandel_userhandle;"`
	Timestamp    time.Time
}

func (newsHandler *NewsHandler) Find() bool {
	db := utils.NewDatabase()
	num := db.Model(&NewsHandler{}).Preload("HashTags", func(db *gorm.DB) *gorm.DB {
		return db.Omit("NewsHandlers")
	}).Preload("UserHandles", func(db *gorm.DB) *gorm.DB {
		return db.Omit("NewsHandlers")
	}).Preload("Media").Where("tweet_id = ?", newsHandler.TweetId).First(&newsHandler)
	if num.RowsAffected > 0 {
		return true
	}
	return false
}

func All() (newsHandles []NewsHandler, err error) {
	db := utils.NewDatabase()
	err = db.Model(&NewsHandler{}).Preload("HashTags", func(db *gorm.DB) *gorm.DB {
		return db.Omit("NewsHandlers")
	}).Preload("UserHandles", func(db *gorm.DB) *gorm.DB {
		return db.Omit("NewsHandlers")
	}).Preload("Media").Find(&newsHandles).Error
	return
}

func (newsHandler *NewsHandler) Insert() (bool, error) {
	db := utils.NewDatabase()
	check := newsHandler.Find()
	if check {
		return false, nil
	}
	return true, db.Create(&newsHandler).Error
}
