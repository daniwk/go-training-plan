package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivityType string

const (
	RUN      ActivityType = "RUN"
	RIDE     ActivityType = "RIDE"
	STRENGTH ActivityType = "STRENGTH"
)

type WorkoutType string

const (
	NORMAL   WorkoutType = "NORMAL"
	WORKOUT  WorkoutType = "WORKOUT"
	LONG_RUN WorkoutType = "LONG_RUN"
)

type PlannedActivity struct {
	gorm.Model
	ActivityType       ActivityType `json:"activity_type"`
	WorkoutType        WorkoutType  `json:"workout_type"`
	WorkoutDescription string       `json:"workout_description"`
	Trail              bool         `json:"trail"`
	Distance           int          `json:"distance"`
	Duration           int          `json:"duration"`
	Intensity          int          `json:"intensity"`
	Day                int          `json:"day"`
	Month              int          `json:"month"`
	Year               int          `json:"year"`
	Date               time.Time    `json:"date"`
	Arvo               bool         `gorm:"default:false" json:"arvo"`
}
