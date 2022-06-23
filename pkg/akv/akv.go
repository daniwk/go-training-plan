package akv

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/spf13/viper"
)

// Inits and return Azure Key Vault client
func getKeyVaultClient() *azsecrets.Client {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	keyVaultName := viper.Get("KEY_VAULT_NAME").(string)
	keyVaultUrl := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)

	//Create a credential using the NewDefaultAzureCredential type.
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	//Establish a connection to the Key Vault client
	client, err := azsecrets.NewClient(keyVaultUrl, cred, nil)
	if err != nil {
		log.Fatalf("failed to connect to client: %v", err)
	}

	return client
}

type UpdateSecretParams struct {
	SecretName       string
	SecretValue      string
	ExpiresInSeconds int
}

// SetSecretExpiration sets expiration on secret in AKV
func SetSecretExpiration(params UpdateSecretParams) {

	// Get AKV client
	client := getKeyVaultClient()

	// Get AKV secret
	getResp, err := client.GetSecret(context.TODO(), params.SecretName, nil)
	if err != nil {
		log.Fatalf("Failed to fetch secret: %v", err)
	}

	// Fetch secret props
	if getResp.Secret.Properties == nil {
		getResp.Secret.Properties = &azsecrets.Properties{}
	}

	// Set new props
	getResp.Secret.Properties = &azsecrets.Properties{
		Enabled:   to.Ptr(true),
		ExpiresOn: to.Ptr(time.Now().Add(time.Duration(params.ExpiresInSeconds) * time.Second)),
		// Remember to preserve the name and version
		Name: getResp.Secret.Properties.Name,
	}

	// Update secret
	resp, err := client.UpdateSecretProperties(context.TODO(), getResp.Secret, nil)
	if err != nil {
		log.Fatalf("Failed to update secret: %v", err)
	}

	fmt.Printf("Updated secret with ID: %s\n", *resp.Secret.ID)
}

// UpdateSecret creates new secrets and changes the values of existing secrets.
func UpdateSecret(params UpdateSecretParams) string {

	// Get AKV client
	client := getKeyVaultClient()

	//Create a secret
	resp, err := client.SetSecret(context.TODO(), params.SecretName, params.SecretValue, nil)
	if err != nil {
		log.Fatalf("failed to connect to client: %v", err)
	}
	fmt.Printf("Updated/created secret with ID: %s\n", *resp.Secret.ID)

	// Set expiration date
	if params.ExpiresInSeconds > 0 {
		SetSecretExpiration(params)
	}

	return *resp.Secret.ID
}

type AKVSecret struct {
	SecretValue string
	Expired     bool
}

// Returns secret for AKV, based on secret_name
func GetSecret(secret_name string) AKVSecret {

	// Get AKV client
	client := getKeyVaultClient()

	// Retrieve secret
	resp, err := client.GetSecret(context.TODO(), secret_name, nil)
	if err != nil {
		log.Fatalf("Failed to fetch secret: %v\n", err)
	}
	fmt.Printf("Fetched secret %s\n", *resp.Secret.ID)
	secret := AKVSecret{}
	secret.SecretValue = *resp.Secret.Value
	if resp.Secret.Properties.ExpiresOn != nil {
		secret.Expired = resp.Secret.Properties.ExpiresOn.Before(time.Now())
	}

	return secret
}
