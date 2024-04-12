package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	Webhook string `yaml:"webhook"`
	LastId  int    `yaml:"lastId"`
}

const DefaultPath = "etc/config.yaml"

var instance *Config

func GetInstance() Config {
	return *instance
}

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

	instance = &config
	return &config
}

func SaveConfig(path string, config Config) error {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// check if file exists
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		_, err = os.Create(path)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}

	instance = &config
	return nil
}
