package statistics

import (
	"net/http"
	"time"

	"github.com/daniwk/training-plan/pkg/models"
	"github.com/gin-gonic/gin"
)

type OverallActivities struct {
	Planned  int `json:"planned"`
	Actual   int `json:"actual"`
	Duration int `json:"duration"`
}

type Distance struct {
	Planned int     `json:"planned"`
	Actual  float64 `json:"actual"`
}

type Workouts struct {
	PlannedNumber             int `json:"planned_number"`
	ActualNumber              int `json:"actual_number"`
	PlannedDuration           int `json:"planned_duration"`
	ActualDuration            int `json:"actual_duration"`
	TotalSufferScore          int `json:"total_suffer_score"`
	AverageHeartRateInWorkout int `json:"avg_heart_rate_in_workout"`
}

type WeeklyRunStatitistics struct {
	StartDate         time.Time
	EndDate           time.Time
	OverallActivities OverallActivities
	Distance          Distance
	Workouts          Workouts
}

func (h handler) GetWeeklyRunStatistics(c *gin.Context) {
	var weekly_stats WeeklyRunStatitistics
	var planned_activities []models.PlannedActivity

	if result := h.DB.Model(&models.PlannedActivity{}).Preload("StravaActivity").Where("activity_type = ?", models.Run).Find(&planned_activities); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	for _, planned_activity := range planned_activities {
		weekly_stats.OverallActivities.Planned += 1
		weekly_stats.Distance.Planned += planned_activity.Distance

		if planned_activity.WorkoutType == models.WORKOUT {
			weekly_stats.Workouts.PlannedNumber += 1
			weekly_stats.Workouts.PlannedDuration += planned_activity.Duration
		}

		if planned_activity.StravaActivity != nil {
			weekly_stats.OverallActivities.Actual += 1
			weekly_stats.OverallActivities.Duration += planned_activity.StravaActivity.MovingTime
			weekly_stats.Distance.Actual += float64(planned_activity.StravaActivity.Distance) / 1000

			if planned_activity.WorkoutType == models.WORKOUT {
				weekly_stats.Workouts.ActualNumber += 1
				weekly_stats.Workouts.ActualDuration += planned_activity.StravaActivity.MovingTime
				weekly_stats.Workouts.TotalSufferScore += int(planned_activity.StravaActivity.SufferScore)
			}
		}
	}

	c.JSON(http.StatusOK, &weekly_stats)
}
