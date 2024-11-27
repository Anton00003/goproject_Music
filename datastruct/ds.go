package datastruct

import (
	"errors"
)

var (
	ErrBadNumPage    = errors.New("input bad number of page")
	ErrBadFilter     = errors.New("input bad filter")
	ErrBadField      = errors.New("input bad field of table")
	ErrBadTypeVal    = errors.New("input value with bad type")
	ErrBadNameGroup  = errors.New("input bad Name and Group values")
	ErrBadNumCouplet = errors.New("input bad number of Couplet")

	TableField = map[string]bool{
		"musicName":        true,
		"musicGroup":       true,
		"musicDate":        true,
		"musicText":        true,
		"musicLink":        true,
		"musicTextCouplet": true,
	}
)

type Music struct {
	MusicName        string   `json:"name"`
	MusicGroup       string   `json:"group"`
	MusicDate        string   `json:"date"`
	MusicText        string   `json:"text"`
	MusicLink        string   `json:"link"`
	MusicTextCouplet []string `json:"couplet"`
}

type SongDetail struct {
	ReleaseDate string
	Text        string
	Link        string
}
