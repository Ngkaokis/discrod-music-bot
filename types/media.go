package types

import (
	"fmt"
	"time"
)

type Media struct {
	Title     string
	FilePath  string
	Uploader  string
	URL       string
	Thumbnail string
	Duration  time.Duration
}

func NewMedia(title, filePath, uploader, url, thumbnail string, durationInSeconds int) *Media {
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", durationInSeconds))
	return &Media{
		Title:     title,
		FilePath:  filePath,
		Uploader:  uploader,
		URL:       url,
		Thumbnail: thumbnail,
		Duration:  duration,
	}
}

type UserActions struct {
	SkipChan chan bool
	StopChan chan bool

	Stopped bool
}

func NewActions() *UserActions {
	return &UserActions{
		SkipChan: make(chan bool, 1),
		StopChan: make(chan bool, 1),
	}
}

