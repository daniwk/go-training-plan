package planned_activities

import (
	"net/http"

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

	planned_activity.ActivityType = models.ActivityType(body.ActivityType)
	planned_activity.Trail = body.Trail

	h.DB.Save(&planned_activity)

	c.JSON(http.StatusOK, &planned_activity)
}
