package akv

import (
	"context"
	"fmt"
	"log"

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

// UpdateSecret creates new secrets and changes the values of existing secrets.
func UpdateSecret(secret_name, secret_value string) string {

	// Get AKV client
	client := getKeyVaultClient()

	//Create a secret
	resp, err := client.SetSecret(context.TODO(), secret_name, secret_value, nil)
	if err != nil {
		log.Fatalf("failed to connect to client: %v", err)
	}

	return *resp.Secret.ID
}

// Returns secret for AKV, based on secret_name
func GetSecret(secret_name string) string {

	// Get AKV client
	client := getKeyVaultClient()

	// Retrieve secret
	resp, err := client.GetSecret(context.TODO(), secret_name, nil)
	if err != nil {
		log.Fatalf("Failed to fetch secret: %v", err)
	}
	fmt.Printf("Fetched secret %s", *resp.Secret.ID)
	return *resp.Secret.Value
}
