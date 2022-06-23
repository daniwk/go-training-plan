package main

import (
	"github.com/daniwk/training-plan/pkg/akv"
	"github.com/daniwk/training-plan/pkg/strava"
	"github.com/spf13/viper"
)

// Upload secret to Azure Key Vault. Credentials (AppID, AppSecret, TenantID) for Azure AD App Reg must exist in env.
func main() {

	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	params := akv.UpdateSecretParams{}
	params.SecretName = "SecretName22"
	params.SecretValue = viper.Get("ENV_VAR").(string)
	params.ExpiresInSeconds = 21600 // 6 hours

	// secret := akv.GetSecret(params.SecretName)
	// fmt.Println(secret)
	// secret_id := akv.UpdateSecret(params)
	// fmt.Printf("Updated secret %s", secret_id)
	// secret := strava.RefreshAccessToken()
	// fmt.Printf("Updated secret %s", secret)
	strava.GetStravaActivities()
}
