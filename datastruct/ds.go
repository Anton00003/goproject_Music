package datastruct

import (
	"errors"
	"time"
)

var (
	ErrBadNumPage    = errors.New("input bad number of page")
	ErrBadFilter     = errors.New("input bad filter")
	ErrBadField      = errors.New("input bad field of table")
	ErrBadTypeVal    = errors.New("input value with bad type")
	ErrBadId         = errors.New("input bad Id values, no rows in result set")
	ErrBadNumCouplet = errors.New("input bad number of Couplet")
	ErrBadGroup      = errors.New("input bad Group values")
	ErrBadGroupId    = errors.New("input bad Group Id values")
	ErrBadList       = errors.New("no songs in result set")
	ErrBadMusicGroup = errors.New("input bad Music and Group values")
)

type Music struct {
	Id      int       `json:"id"`
	Name    string    `json:"name"`
	GroupId int       `json:"group"`
	Date    time.Time `json:"date"` // перевести в time.Time
	Text    string    `json:"text"`
	Link    string    `json:"link"`
}

type MusicListItem struct { //    MusicListItem
	Id      int       `json:"id"`
	Name    string    `json:"name"`
	GroupId int       `json:"groupId"`
	Date    time.Time `json:"date"` // перевести в time.Time
	Text    string    `json:"text"`
	Link    string    `json:"link"`
	Group   string    `json:"group"`
}

type SongDetail struct {
	ReleaseDate string
	Text        string
	Link        string
}
