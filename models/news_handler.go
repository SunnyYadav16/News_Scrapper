package models

import "time"

type NewsHandler struct {
	TweetId      string
	ChannelName  string
	TweetContent string
	Media        []Media
	HashTags     []string
	UserHandle   []string
	Timestamp    time.Time
}
