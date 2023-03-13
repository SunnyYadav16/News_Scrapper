package models

import (
	"time"
)

type NewsHandler struct {
	TweetId      string `gorm:"primary_key;"`
	ChannelName  string
	TweetContent string
	Timestamp    time.Time
	Media        []Media       `gorm:"foreignKey:TweetId"`
	HashTags     []Hashtags    `gorm:"foreignKey:TweetId"`
	UserHandles  []UserHandles `gorm:"foreignKey:TweetId"`
}
