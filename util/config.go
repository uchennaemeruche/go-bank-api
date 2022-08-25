package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource            string        `mapstructure:"DB_SOURCE"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(envPath string) (config Config, err error) {

	viper.AddConfigPath(envPath)
	viper.SetConfigName("app") // TODO: set based on test, staging, dev or production
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return

}
