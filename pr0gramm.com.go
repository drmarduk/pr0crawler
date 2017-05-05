package main

import (
	"fmt"
	"strings"
)

// ProgrammItem is a single Upload item with metadata stuff
type ProgrammItem struct {
	ID      int    `json:"id"`
	Created int    `json:"created"`
	User    string `json:"user"`

	Down  int `json:"down"`
	Up    int `json:"up"`
	Flags int `json:"flags"`

	Audio bool `json:"audio"`

	Image    string `json:"image"`
	Fullsize string `json:"fullsize"`
	Thumb    string `json:"thumb"`
	Source   string `json:"source"`

	Height int `json:"height"`
	Width  int `json:"width"`

	Mark     int `json:"mark"`
	Promoted int `json:"promoted"`
}

// URL returns a formatted url based on the file extension
func (i *ProgrammItem) URL() string {
	if strings.HasSuffix(i.Image, "mp4") {
		return fmt.Sprintf("https://vid.pr0gramm.com/%s", i.Image)
	}
	return fmt.Sprintf("https://img.pr0gramm.com/%s", i.Image)

}

// ProgrammItemsGet is a request for a user with his uploads
type ProgrammItemsGet struct {
	AtEnd   bool           `json:"atEnd"`
	AtStart bool           `json:"atStart"`
	Cache   string         `json:"cache"`
	Error   string         `json:"error"`
	Items   []ProgrammItem `json:"items"`
	Qc      int            `json:"qc"`
	Rt      int            `json:"rt"`
	Ts      int            `json:"ts"`
}

// NewProgrammItemsGet returns an good object
func NewProgrammItemsGet() *ProgrammItemsGet {
	return &ProgrammItemsGet{
		Items: make([]ProgrammItem, 120),
	}
}
