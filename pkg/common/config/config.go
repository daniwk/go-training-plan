package config

import "github.com/spf13/viper"

type Config struct {
	Port               string `mapstructure:"PORT"`
	DBUrl              string `mapstructure:"DB_URL"`
	KeyVaultName       string `mapstructure:"KEY_VAULT_NAME"`
	StravaAccessToken  string `mapstructure:"STRAVA_ACCESS_TOKEN"`
	StravaRefreshToken string `mapstructure:"STRAVA_REFRESH_TOKEN"`
	AzureAppID         string `mapstructure:"AZURE_APP_ID"`
	AzureAppSecret     string `mapstructure:"AZURE_APP_SECRET"`
	AzureTenantID      string `mapstructure:"AZURE_TENANT_ID"`
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath("./pkg/common/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
