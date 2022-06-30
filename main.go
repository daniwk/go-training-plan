package main

import (
	"github.com/daniwk/training-plan/pkg/api/v1/feelings"
	"github.com/daniwk/training-plan/pkg/api/v1/planned_activities"
	"github.com/daniwk/training-plan/pkg/api/v1/statistics"
	strava_activites "github.com/daniwk/training-plan/pkg/api/v1/strava_activities"
	"github.com/daniwk/training-plan/pkg/common/db"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	port := viper.Get("PORT").(string)
	dbUrl := viper.Get("DB_URL").(string)

	r := gin.Default()
	h := db.Init(dbUrl)

	planned_activities.RegisterRoutes(r, h)
	strava_activites.RegisterRoutes(r, h)
	statistics.RegisterRoutes(r, h)
	feelings.RegisterRoutes(r, h)

	r.Run(port)
}
