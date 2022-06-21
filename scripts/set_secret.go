package main

import (
	"fmt"

	"github.com/daniwk/training-plan/pkg/akv"
	"github.com/spf13/viper"
)

// Upload secret to Azure Key Vault. Credentials (AppID, AppSecret, TenantID) for Azure AD App Reg must exist in env.
func main() {

	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	mySecretName := "SecretName"
	mySecretvalue := viper.Get("ENV_VAR").(string)

	secret_id := akv.UpdateSecret(mySecretName, mySecretvalue)
	fmt.Printf("Set secret %s", secret_id)
}
