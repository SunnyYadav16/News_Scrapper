package models

import "gorm.io/gorm"

type UserHandles struct {
	gorm.Model
	TweetId    string
	UserHandle string
}
