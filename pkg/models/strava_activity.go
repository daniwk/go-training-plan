package models

type StravaActivityResponse struct {
	StravaActivityID  int     `json:"id"`
	Name              string  `json:"name"`
	Distance          float64 `json:"distance"`
	MovingTime        int     `json:"moving_time"`
	StartDate         string  `json:"start_date"`
	SportType         string  `json:"sport_type"`
	PerceivedExertion float64 `json:"perceived_exertion"`
}

type StravaActivity struct {
	ID                string  `gorm:"primaryKey;default:gen_random_uuid()"`
	StravaActivityID  int     `gorm:"not null;uniqueIndex:idx_activityid" json:"id"`
	Name              string  `json:"name"`
	Distance          float64 `json:"distance"`
	MovingTime        int     `json:"moving_time"`
	StartDate         string  `json:"start_date"`
	SportType         string  `json:"sport_type"`
	PerceivedExertion float64 `json:"perceived_exertion"`
}
