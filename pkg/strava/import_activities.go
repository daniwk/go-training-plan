package strava

import (
	"github.com/google/uuid"
)

type StravaActivity struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Distance          int       `json:"distance,string"`
	MovingTime        int       `json:"moving_time,string"`
	StartTime         string    `json:"start_time"`
	SportType         string    `json:"sport_type"`
	PerceivedExertion int       `json:"perceived_exertion,string"`
}

func GetStravaActivities() {

}
