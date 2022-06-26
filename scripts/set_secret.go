package main

import (
	"fmt"

	"github.com/daniwk/training-plan/pkg/akv"
	"github.com/spf13/viper"
)

// Upload secret to Azure Key Vault. Credentials (AppID, AppSecret, TenantID) for Azure AD App Reg must exist in env.
func main2() {

	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	params := akv.UpdateSecretParams{}
	params.SecretName = "StravaAccessToken"
	params.SecretValue = viper.Get("ENV_VAR").(string)
	params.ExpiresInSeconds = 21600 // 6 hours

	secret_id := akv.GetSecret(params.SecretName)
	fmt.Print("Secret value: ", secret_id)
}
