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

func NewConfig() *Config {
	file, err := os.Open(DefaultPath)
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

func (c Config) Save() error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// check if file exists
	_, err = os.Stat(DefaultPath)
	if os.IsNotExist(err) {
		_, err = os.Create(DefaultPath)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(DefaultPath, bytes, 0644)
	if err != nil {
		return err
	}

	instance = &c
	return nil
}
