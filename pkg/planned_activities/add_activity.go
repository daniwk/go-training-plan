package planned_activities

import (
	"net/http"

	"github.com/daniwk/training-plan/pkg/common/models"
	"github.com/gin-gonic/gin"
)

type AddPlannedActivityRequestBody struct {
	ActivityType string `json:"activity_type"`
	Trail        bool   `json:"trail,string"`
	Day          int    `json:"day,string"`
	Month        int    `json:"month,string"`
	Year         int    `json:"year,string"`
	Distance     int    `json:"distance,string"`
	Duration     int    `json:"duration,string"`
	Intensity    int    `json:"intensity,string"`
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

	if result := h.DB.Create(&planned_activity); result.Error != nil {
		c.AbortWithError(http.StatusBadRequest, result.Error)
		return
	}

	c.JSON(http.StatusCreated, &planned_activity)
}
