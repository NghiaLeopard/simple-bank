package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDrive       string `mapstructure:"DB_DRIVE"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ACTION"`
	SymmetricKey string `mapstructure:"SYMMETRICKEY"`
	Duration time.Duration `mapstructure:"DURATION"`
}

func LoadConfig(path string) (config Config,err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// allow override 
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return 
}