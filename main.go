package main

import (
	"github.com/daniwk/training-plan/pkg/api/v1/planned_activities"
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
	// register more routes here

	r.Run(port)
}
