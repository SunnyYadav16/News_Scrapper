package models

import "gorm.io/gorm"

type Media struct {
	gorm.Model
	TweetId string
	Type    string
	URL     string
}
