package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	Webhook string `yaml:"webhook"`
}

// Instance is a global variable that stores the config instance
var Instance Config

func NewConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	config := Config{}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil
	}

	return &config
}
