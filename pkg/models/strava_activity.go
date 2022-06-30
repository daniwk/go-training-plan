package models

import (
	"time"
)

type Lap struct {
	ID               string  `gorm:"primaryKey;default:gen_random_uuid()"`
	StravaID         int     `gorm:"not null;uniqueIndex" json:"id"`
	Split            int     `json:"split"`
	MovingTime       int     `json:"moving_time"`
	Distance         float64 `json:"distance"`
	AverageSpeed     float64 `json:"average_speed"`
	MaxSpeed         float64 `json:"max_speed"`
	AverageHeartRate float64 `json:"average_heartrate"`
	MaxHeartRate     float64 `json:"max_heartrate"`
	StravaActivityID string
}

// type StravaActivityAPIResponse struct {
// 	StravaID          int     `json:"id"`
// 	Name              string  `json:"name"`
// 	Distance          float64 `json:"distance"`
// 	MovingTime        int     `json:"moving_time"`
// 	StartDate         string  `json:"start_date"`
// 	SportType         string  `json:"sport_type"`
// 	PerceivedExertion float64 `json:"perceived_exertion"`
// 	WorkoutType       int     `json:"workout_type"`
// 	SufferScore       float64 `json:"suffer_score"`
// 	Laps              []Lap   `json:"laps"`
// }

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
	Laps              []Lap     `json:"laps"`
	PlannedActivityID *uint
}
