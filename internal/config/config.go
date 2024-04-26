package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	Webhooks       []string        `yaml:"webhooks"`
	LastId         int             `yaml:"lastId"`
	MonitoredPages []MonitoredPage `yaml:"monitored_pages"`
}

type MonitoredPage struct {
	Id   int    `yaml:"id"`
	Hash string `yaml:"hash"`
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
