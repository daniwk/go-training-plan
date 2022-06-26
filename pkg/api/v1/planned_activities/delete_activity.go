package planned_activities

import (
	"net/http"

	"github.com/daniwk/training-plan/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h handler) DeletePlannedActivity(c *gin.Context) {
	id := c.Param("id")

	var planned_activity models.PlannedActivity

	if result := h.DB.First(&planned_activity, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	h.DB.Delete(&planned_activity)

	c.Status(http.StatusOK)
}
