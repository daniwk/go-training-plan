package models

import (
	"time"
)

type StravaActivityResponse struct {
	StravaID          int     `json:"id"`
	Name              string  `json:"name"`
	Distance          float64 `json:"distance"`
	MovingTime        int     `json:"moving_time"`
	StartDate         string  `json:"start_date"`
	SportType         string  `json:"sport_type"`
	PerceivedExertion float64 `json:"perceived_exertion"`
	WorkoutType       int     `json:"workout_type"`
	SufferScore       float64 `json:"suffer_score"`
}

type StravaActivity struct {
	ID                string    `gorm:"primaryKey;default:gen_random_uuid()"`
	StravaID          int       `gorm:"not null;uniqueIndex:idx_stravaid" json:"id"`
	Name              string    `json:"name"`
	Distance          float64   `json:"distance"`
	MovingTime        int       `json:"moving_time"`
	StartDate         time.Time `json:"start_date"`
	SportType         string    `json:"sport_type"`
	WorkoutType       int       `json:"workout_type"`
	PerceivedExertion float64   `json:"perceived_exertion"`
	SufferScore       float64   `json:"suffer_score"`
	PlannedActivityID *uint
}
