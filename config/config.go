package config

import (
	"fmt"
	"os"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database Database `yaml:"database" validate:"required"`
}

type Database struct {
	MigrationsSource string `yaml:"migrations_source" validate:"required"`
	DbSource         string `yaml:"db_source" validate:"required"`
}

func NewConfig(configPath string) (*Config, error) {
	var cfg Config
	cfgFile, openErr := os.Open(configPath)
	if openErr != nil {
		return nil, fmt.Errorf("cannot open config path(%q): %w", configPath, openErr)
	}
	cfgDecoder := yaml.NewDecoder(cfgFile)
	decodeErr := cfgDecoder.Decode(&cfg)
	if decodeErr != nil {
		return nil, fmt.Errorf("cannot parse config: %w", decodeErr)
	}

	yamlValidator := validator.New()
	validateErr := yamlValidator.Struct(&cfg)
	if validateErr != nil {
		return nil, fmt.Errorf("error validating config file: %w", validateErr)
	}

	return &cfg, nil
}
