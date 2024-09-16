package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port  string   `yaml:"port"`
	Modes []string `yaml:"modes"`
}

func NewConfig(path string) Config {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return Config{}
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
		return Config{}
	}

	return config
}
