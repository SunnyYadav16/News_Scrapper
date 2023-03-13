package models

import (
	"time"
)

type NewsHandler struct {
	TweetId      string `gorm:"primary_key;"`
	ChannelName  string
	TweetContent string
	Timestamp    time.Time
	HashTags     []Hashtags    `gorm:"foreignKey:TweetId"`
	UserHandle   []UserHandles `gorm:"foreignKey:TweetId"`
	Media        []Media       `gorm:"foreignKey:TweetId"`
}
type Media struct {
	TweetId string
	Type    string
	URL     string
}
type Hashtags struct {
	TweetId  string
	HashTags string
}
type UserHandles struct {
	TweetId    string
	UserHandle string
}
