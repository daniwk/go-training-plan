package strava

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/daniwk/training-plan/pkg/akv"
	"github.com/daniwk/training-plan/pkg/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefreshAccessTokenResponseBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    int    `json:"expires_at"`
	TokenType    string `json:"token_type"`
}

// RefreshAccessToken renews Strava access and refresh tokens, update AKV secrets and returns the access token
func RefreshAccessToken() string {

	// Fetch new Strava access and refresh tokens
	strava_client_id := viper.Get("STRAVA_CLIENT_ID").(string)
	strava_client_secret := viper.Get("STRAVA_CLIENT_SECRET").(string)
	strava_refresh_token_akv_name := viper.Get("STRAVA_REFRESH_TOKEN_AKV_NAME").(string)
	strava_access_token_akv_name := viper.Get("STRAVA_ACCESS_TOKEN_AKV_NAME").(string)
	strava_refresh_token := akv.GetSecret(strava_refresh_token_akv_name).SecretValue
	strava_refresh_access_token_url := viper.Get("STRAVA_REFRESH_ACCESS_TOKEN_URL").(string)

	data := url.Values{}
	data.Set("client_id", strava_client_id)
	data.Set("client_secret", strava_client_secret)
	data.Set("refresh_token", strava_refresh_token)
	data.Set("grant_type", "refresh_token")

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, strava_refresh_access_token_url, strings.NewReader(data.Encode()))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("HTTP POST request failed with error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("HTTP POST request failed with error: %v", err)
	}

	response_body := RefreshAccessTokenResponseBody{}
	if err := json.Unmarshal(body, &response_body); err != nil {
		log.Fatalf("Cannot unmarshal json: %v", err)
	}

	// Update secrets in AKV
	access_token_params := akv.UpdateSecretParams{SecretName: strava_access_token_akv_name, SecretValue: response_body.AccessToken, ExpiresInSeconds: response_body.ExpiresIn}
	akv.UpdateSecret(access_token_params)
	refresh_token_params := akv.UpdateSecretParams{SecretName: strava_refresh_token_akv_name, SecretValue: response_body.RefreshToken}
	akv.UpdateSecret(refresh_token_params)

	return response_body.AccessToken
}

// List a Strava athletes activities. Returns []StravaActivity
func ListStravaAthleteActivites(StravaAccessToken string, ResultsPerPage, Page int) []models.StravaActivity {
	strava_url := fmt.Sprintf("https://www.strava.com/api/v3/athlete/activities?per_page=%d&page=%d", ResultsPerPage, Page)
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, strava_url, nil)
	access_token := fmt.Sprintf("Bearer %s", StravaAccessToken)
	req.Header.Add("Authorization", access_token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("HTTP GET request failed with error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("HTTP POST request failed with error: %v", err)
	}
	strava_activity_records := []models.StravaActivity{}
	if err := json.Unmarshal(body, &strava_activity_records); err != nil {
		log.Fatalf("Cannot unmarshal json: %v", err)
	}

	return strava_activity_records
}

// Get a Strava Activity by its ID
func GetStravaActivityByID(StravaAccessToken string, StravaActivityID int) models.StravaActivity {
	strava_url := fmt.Sprintf("https://www.strava.com/api/v3/activities/%d", StravaActivityID)
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, strava_url, nil)
	access_token := fmt.Sprintf("Bearer %s", StravaAccessToken)
	req.Header.Add("Authorization", access_token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("HTTP GET request failed with error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("HTTP POST request failed with error: %v", err)
	}
	strava_activity_record := models.StravaActivity{}
	if err := json.Unmarshal(body, &strava_activity_record); err != nil {
		log.Fatalf("Cannot unmarshal json: %v", err)
	}

	return strava_activity_record
}

// GetStravaActivities retrieves last 10 Strava activites and uploads them to DB
func GetStravaActivities() {

	// DB
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()
	dbUrl := viper.Get("DB_URL").(string)
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("HTTP GET request failed with error: %v", err)
	}
	db.AutoMigrate(&models.StravaActivity{})

	// Fetch StravaAccessToken from AKV
	secret_name := viper.Get("STRAVA_ACCESS_TOKEN_AKV_NAME").(string)
	AKVSecret := akv.GetSecret(secret_name)
	if AKVSecret.Expired {
		// Refresh
		AKVSecret.SecretValue = RefreshAccessToken()
	}

	// Get Strava Activities
	ResultsPerPage := 5
	Page := 1
	strava_activity_records := ListStravaAthleteActivites(AKVSecret.SecretValue, ResultsPerPage, Page)

	// Get latest PlannedActivites
	for _, strava_activity := range strava_activity_records {

		// Get Activity to get extra data 😍
		strava_activity := GetStravaActivityByID(AKVSecret.SecretValue, strava_activity.StravaID)

		// Upsert activity
		fmt.Printf("Working with Strava Activity: %v\n", strava_activity)
		existing_strava := models.StravaActivity{}
		if result := db.Where("strava_id = ?", strava_activity.StravaID).Find(&existing_strava); result != nil {
			strava_activity.ID = existing_strava.ID
		}

		// Match Strava activity with planned activity by finding correct activity based on type, date and "Arvo"
		planned_activities := []models.PlannedActivity{}
		activity_type := strava_activity.SportType
		if result := db.Where("day = ? AND month = ? AND year = ? AND activity_type = ?", strava_activity.StartDate.Day(), strava_activity.StartDate.Month(), strava_activity.StartDate.Year(), activity_type).Find(&planned_activities); result.Error != nil {
			log.Fatalf("Cannot unmarshal json: %v", result.Error)
		}
		fmt.Println("Following planned activities where found: \n", planned_activities)

		if len(planned_activities) > 0 {

			// Find planned activity and match
			for _, planned_activity := range planned_activities {

				planned_activity_time := time.Date(planned_activity.Year, time.Month(planned_activity.Month), planned_activity.Day, 12, 0, 0, 0, time.Local)
				if planned_activity.Arvo && strava_activity.StartDate.After(planned_activity_time) || !planned_activity.Arvo && strava_activity.StartDate.Before(planned_activity_time) {

					fmt.Println("Following planned activity where matched: \n", planned_activity)
					planned_activity.StravaActivity = &strava_activity
					fmt.Println("Updating planned activity to: \n", planned_activity)

					if result := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&planned_activity); result.Error != nil {
						log.Fatalf("Couldnt insert planned_activity to db: %v", result.Error)
						continue
					}
				}
			}

		} else {

			// Upload strava activity without binding
			fmt.Print("No planned activities found, adding strava activity without binding: \n")
			db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&strava_activity)
			fmt.Println("Strava activity created in DB: \n", strava_activity)
		}

	}
}
