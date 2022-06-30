package strava_activites

import (
	"net/http"

	"github.com/daniwk/training-plan/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetStravaActivites(c *gin.Context) {
	var strava_activities []models.StravaActivity

	if result := h.DB.Model(&models.StravaActivity{}).Preload("Laps").Order("start_date desc").Find(&strava_activities); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, &strava_activities)
}
