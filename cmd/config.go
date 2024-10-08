package main

import (
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string   `yaml:"port"`
	LoggingMode string   `yaml:"logLevel"`
	Modes       []string `yaml:"modes"`
}

func NewConfig(path string) Config {
	viper.SetDefault("LoggingLevel", "INFO")

	if path != "" {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(path)
	}
	viper.SetEnvPrefix("DAEMON")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		slog.Error(err.Error())
	}

	config := Config{
		Port:        viper.GetString("Port"),
		LoggingMode: viper.GetString("Log_Level"),
		Modes:       viper.GetStringSlice("Modes"),
	}

	return config
}
