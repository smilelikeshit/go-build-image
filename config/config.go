package config

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	VaultDNS      string `mapstructure:"VAULT_URL"`
	VaultAppToken string `mapstructure:"VAULT_APP_TOKEN"`
	VaultUsername string `mapstructure:"VAULT_USERNAME"`
	VaultPassword string `mapstructure:"VAULT_PASSWORD"`
	VaultPath     string `mapstructure:"VAULT_PATH"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
