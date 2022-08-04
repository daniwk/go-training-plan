package planned_activities

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	h := &handler{
		DB: db,
	}

	routes := r.Group("/plannedActivites")
	routes.GET("/", h.GetPlannedActivites)
	routes.POST("/", h.AddPlannedActivity)
	routes.GET("/:id", h.GetPlannedActivity)
	routes.DELETE("/:id", h.DeletePlannedActivity)
	routes.PUT("/:id", h.UpdatePlannedActivity)
}
