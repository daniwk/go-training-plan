package feelings

import (
	"fmt"
	"net/http"
	"time"

	"github.com/daniwk/training-plan/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (h handler) AddFeeling(c *gin.Context) {
	body := models.AddFeelingRequestBody{}

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var feeling models.Feeling

	feeling.Day = body.Day
	feeling.Month = body.Month
	feeling.Year = body.Year
	feeling.Date = time.Date(body.Year, time.Month(body.Month), body.Day, 9, 0, 0, 0, time.Local)
	feeling.Arvo = body.Arvo
	feeling.Feeling = body.Feeling

	fmt.Printf("Upserting: %v", feeling)

	if result := h.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&feeling); result.Error != nil {
		c.AbortWithError(http.StatusBadRequest, result.Error)
		return
	}

	c.JSON(http.StatusCreated, &feeling)
}
