package planned_activities

import (
	"net/http"

	"github.com/daniwk/training-plan/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetPlannedActivity(c *gin.Context) {
	id := c.Param("id")

	var planned_activity models.PlannedActivity
	if result := h.DB.Model(&models.PlannedActivity{}).Preload("StravaActivity").Preload("StravaActivity.Laps").Order("date desc, arvo").First(&planned_activity, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, &planned_activity)
}
