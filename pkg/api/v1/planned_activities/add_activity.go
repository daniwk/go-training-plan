package planned_activities

import (
	"net/http"
	"time"

	"github.com/daniwk/training-plan/pkg/models"
	"github.com/gin-gonic/gin"
)

type AddPlannedActivityRequestBody struct {
	ActivityType       string             `json:"activity_type"`
	Trail              bool               `json:"trail"`
	Day                int                `json:"day"`
	Month              int                `json:"month"`
	Year               int                `json:"year"`
	Distance           int                `json:"distance"`
	Duration           int                `json:"duration"`
	Intensity          int                `json:"intensity"`
	Arvo               bool               `json:"arvo"`
	WorkoutType        models.WorkoutType `json:"workout_type"`
	WorkoutDescription string             `json:"workout_description"`
}

func (h handler) AddPlannedActivity(c *gin.Context) {
	body := AddPlannedActivityRequestBody{}

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var planned_activity models.PlannedActivity

	planned_activity.ActivityType = models.ActivityType(body.ActivityType)
	planned_activity.Trail = body.Trail
	planned_activity.Day = body.Day
	planned_activity.Month = body.Month
	planned_activity.Year = body.Year
	planned_activity.Distance = body.Distance
	planned_activity.Duration = body.Duration
	planned_activity.Intensity = body.Intensity
	planned_activity.Date = time.Date(body.Year, time.Month(body.Month), body.Day, 9, 0, 0, 0, time.Local)
	planned_activity.Arvo = body.Arvo
	planned_activity.WorkoutDescription = body.WorkoutDescription
	planned_activity.WorkoutType = body.WorkoutType

	if result := h.DB.Create(&planned_activity); result.Error != nil {
		c.AbortWithError(http.StatusBadRequest, result.Error)
		return
	}

	c.JSON(http.StatusCreated, &planned_activity)
}
