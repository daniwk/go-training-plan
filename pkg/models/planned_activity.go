package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivityType string

const (
	Run            ActivityType = "Run"
	Ride           ActivityType = "Ride"
	WeightTraining ActivityType = "WeightTraining"
)

type WorkoutType string

const (
	NORMAL   WorkoutType = "NORMAL"
	WORKOUT  WorkoutType = "WORKOUT"
	LONG_RUN WorkoutType = "LONG_RUN"
)

type WorkoutDetails struct {
	gorm.Model
	WorkoutDescription string `json:"workout_description"`
	MinutesInQuality   int    `json:"minutes_in_quality"`
}

type PlannedActivity struct {
	gorm.Model
	ActivityType       ActivityType `json:"activity_type"`
	WorkoutType        WorkoutType  `json:"workout_type"`
	WorkoutDescription string       `json:"workout_description"`
	MinutesInQuality   int          `json:"minutes_in_quality"`
	Trail              bool         `json:"trail"`
	Distance           int          `json:"distance"`
	Duration           int          `json:"duration"`
	Intensity          int          `json:"intensity"`
	Day                int          `json:"day"`
	Month              int          `json:"month"`
	Year               int          `json:"year"`
	Date               time.Time    `json:"date"`
	Arvo               bool         `gorm:"default:false" json:"arvo"`
	StravaActivity     *StravaActivity
}

type AddPlannedActivityRequestBody struct {
	ActivityType       string      `json:"activity_type"`
	Trail              bool        `json:"trail"`
	Day                int         `json:"day"`
	Month              int         `json:"month"`
	Year               int         `json:"year"`
	Distance           int         `json:"distance"`
	Duration           int         `json:"duration"`
	Intensity          int         `json:"intensity"`
	Arvo               bool        `json:"arvo"`
	WorkoutType        WorkoutType `json:"workout_type"`
	WorkoutDescription string      `json:"workout_description"`
	MinutesInQuality   int         `json:"minutes_in_quality"`
}
