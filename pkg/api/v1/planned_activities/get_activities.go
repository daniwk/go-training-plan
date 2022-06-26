package planned_activities

import (
	"net/http"

	"github.com/daniwk/training-plan/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetPlannedActivites(c *gin.Context) {
	var planned_activities []models.PlannedActivity

	if result := h.DB.Model(&models.PlannedActivity{}).Preload("StravaActivity").Find(&planned_activities); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, &planned_activities)
}
