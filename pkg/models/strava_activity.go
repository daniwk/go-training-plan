package models

import "gorm.io/gorm"

type StravaActivityResponse struct {
	StravaID          int     `json:"id"`
	Name              string  `json:"name"`
	Distance          float64 `json:"distance"`
	MovingTime        int     `json:"moving_time"`
	StartDate         string  `json:"start_date"`
	SportType         string  `json:"sport_type"`
	PerceivedExertion float64 `json:"perceived_exertion"`
}

type StravaActivity struct {
	gorm.Model
	StravaActivityResponse
}
