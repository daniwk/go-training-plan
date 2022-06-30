package models

import "time"

type AddFeelingRequestBody struct {
	Day     int  `json:"day"`
	Month   int  `json:"month"`
	Year    int  `json:"year"`
	Arvo    bool `json:"arvo"`
	Feeling int  `json:"feeling"`
}

type Feeling struct {
	ID      string    `gorm:"primaryKey;default:gen_random_uuid()"`
	Day     int       `json:"day"`
	Month   int       `json:"month"`
	Year    int       `json:"year"`
	Date    time.Time `json:"date"`
	Arvo    bool      `json:"arvo"`
	Feeling int       `json:"feeling"`
}
