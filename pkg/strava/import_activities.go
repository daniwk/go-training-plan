package strava

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/daniwk/training-plan/pkg/akv"
	"github.com/daniwk/training-plan/pkg/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	strava_url := "https://www.strava.com/api/v3/athlete/activities?per_page=10"
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, strava_url, nil)
	access_token := fmt.Sprintf("Bearer %s", AKVSecret.SecretValue)
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

	response_body := []models.StravaActivity{}
	if err := json.Unmarshal(body, &response_body); err != nil {
		log.Fatalf("Cannot unmarshal json: %v", err)
	}
	if result := db.Create(&response_body); result.Error != nil {
		log.Fatalf("Couldnt insert data to db: %v", result.Error)
	}

}
