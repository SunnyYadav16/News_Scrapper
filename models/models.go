package models

import (
	"gorm.io/gorm"
	"time"
)

type NewsHandler struct {
	gorm.Model
	TweetId      string `gorm:"primaryKey"`
	ChannelName  string
	TweetContent string
	Media        []Media
	HashTags     []Hashtags
	UserHandle   []UserHandles
	Timestamp    time.Time
}
type Media struct {
	gorm.Model
	TweetId string
	Type    string
	URL     string
}

type Hashtags struct {
	gorm.Model
	TweetId  string
	Hashtags string
}
type UserHandles struct {
	gorm.Model
	TweetId    string
	UserHandle string
}
