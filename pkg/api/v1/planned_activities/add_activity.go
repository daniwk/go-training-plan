package planned_activities

import (
	"fmt"
	"net/http"
	"time"

	"github.com/daniwk/training-plan/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (h handler) AddPlannedActivity(c *gin.Context) {
	body := models.AddPlannedActivityRequestBody{}

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
	planned_activity.Date = time.Date(body.Year, time.Month(body.Month), body.Day, 0, 0, 0, 0, time.Local)
	planned_activity.Arvo = body.Arvo
	planned_activity.WorkoutDescription = body.WorkoutDescription
	planned_activity.MinutesInQuality = body.MinutesInQuality
	planned_activity.WorkoutType = body.WorkoutType

	fmt.Printf("Upserting: %v", planned_activity)

	if result := h.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&planned_activity); result.Error != nil {
		c.AbortWithError(http.StatusBadRequest, result.Error)
		return
	}

	c.JSON(http.StatusCreated, &planned_activity)
}
