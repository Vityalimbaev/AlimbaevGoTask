package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type AppConfig struct {
	DataPath string `mapstructure:"DIR_PATH"`
}

func GetConfig() AppConfig {

	viper.AddConfigPath("./app")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("Error: Can't read config file, because: %v", err)
		os.Exit(1)
	}

	viper.AutomaticEnv()

	appConfig := AppConfig{}

	if err := viper.Unmarshal(&appConfig); err != nil {
		fmt.Printf("Can't unmarshal app config, because %v", err)
		os.Exit(1)
	}

	return appConfig
}
