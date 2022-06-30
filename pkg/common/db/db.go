package db

import (
	"log"

	"github.com/daniwk/training-plan/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.PlannedActivity{})
	db.AutoMigrate(&models.StravaActivity{})
	db.AutoMigrate(&models.Lap{})
	db.AutoMigrate(&models.Feeling{})

	return db
}
