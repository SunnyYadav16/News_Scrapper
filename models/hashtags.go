package models

import "gorm.io/gorm"

type Hashtags struct {
	gorm.Model
	TweetId  string
	Hashtags string
}
