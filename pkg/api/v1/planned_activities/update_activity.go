package planned_activities

import (
	"net/http"
	"time"

	"github.com/daniwk/training-plan/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h handler) UpdatePlannedActivity(c *gin.Context) {
	id := c.Param("id")
	body := models.AddPlannedActivityRequestBody{}

	// getting request's body
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var planned_activity models.PlannedActivity

	if result := h.DB.First(&planned_activity, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	// planned_activity.WorkoutType = body.WorkoutType
	// planned_activity.WorkoutDescription = body.WorkoutDescription
	// planned_activity.MinutesInQuality = body.MinutesInQuality
	// planned_activity.Distance = body.Distance
	// planned_activity.Intensity = body.Intensity

	planned_activity.Day = body.Day
	planned_activity.Month = body.Month
	planned_activity.Year = body.Year
	planned_activity.Date = time.Date(body.Year, time.Month(body.Month), body.Day, 0, 0, 0, 0, time.Local)

	h.DB.Save(&planned_activity)

	c.JSON(http.StatusOK, &planned_activity)
}
