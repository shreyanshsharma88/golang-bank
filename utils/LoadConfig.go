package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string `mapstructure:"SYMMETRIC_KEY"`
	ExpiryDuration time.Duration `mapstructure:"EXPIRY_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	// var config Config
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return Config{} , err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{} , err
	}
	return  config, nil
}
