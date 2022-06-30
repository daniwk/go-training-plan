package feelings

import (
	"net/http"

	"github.com/daniwk/training-plan/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetFeelings(c *gin.Context) {
	var feelings []models.Feeling

	if result := h.DB.Model(&models.Feeling{}).Order("date desc").Find(&feelings); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, &feelings)
}
