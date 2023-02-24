package model

import (
	"time"
)

type Picture struct {
	Date        time.Time
	Title       string
	URL         string
	HDURL       *string
	ThumbURL    *string
	LocalURL    string
	MediaType   string
	Copyright   *string
	Explanation string
}

// ChooseURL determines which url should be assigned as a main url.
func (p *Picture) ChooseURL() (url string) {
	switch {
	case p.MediaType == "image" && p.HDURL != nil:
		return *p.HDURL
	case p.MediaType == "image" && p.URL != "":
		return p.URL
	case p.MediaType == "video" && p.ThumbURL != nil:
		return *p.ThumbURL
	default:
		return ""
	}
}
