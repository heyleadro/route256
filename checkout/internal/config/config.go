package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const pathToConfig = "config.yaml"

type Config struct {
	Token    string `yaml:"token"`
	Services struct {
		Loms    string `yaml:"loms"`
		Product string `yaml:"product"`
	} `yaml:"services"`
	PostgresDB struct {
		URL string `yaml:"url"`
	} `yaml:"postgres-db"`
}

var AppConfig = Config{}

func Init() error {
	rawYaml, err := os.ReadFile(pathToConfig)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	err = yaml.Unmarshal(rawYaml, &AppConfig)
	if err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	return nil
}
