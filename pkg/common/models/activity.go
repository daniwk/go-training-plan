package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivityType string

const (
	RUN      ActivityType = "RUN"
	BICYCLE  ActivityType = "BICYCLE"
	STRENGTH ActivityType = "STRENGTH"
)

type PlannedActivity struct {
	gorm.Model
	ActivityType ActivityType `json:"activity_type"`
	Trail        bool         `json:"trail"`
	Day          int          `json:"day"`
	Month        int          `json:"month"`
	Year         int          `json:"year"`
	Date         time.Time    `json:"date"`
	Distance     int          `json:"distance"`
	Duration     int          `json:"duration"`
	Intensity    int          `json:"intensity"`
}
