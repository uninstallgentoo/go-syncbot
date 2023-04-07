package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Server struct {
	Host   string `yaml:"host"`
	Secure bool   `yaml:"secure"`
}

type Channel struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type User struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type Database struct {
	Path string `yaml:"path"`
}

type WeatherAPI struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

type ExternalServices struct {
	Weather WeatherAPI `yaml:"weather"`
}

type Config struct {
	Server           `yaml:"server"`
	Channel          `yaml:"channel"`
	User             `yaml:"user"`
	Database         `yaml:"database"`
	ExternalServices `yaml:"external_services"`
}

const BotConfigFile = "config.yaml"

func NewConfig() (*Config, error) {
	c := &Config{}
	configPath := os.Getenv("CONFIG_PATH")
	configFilePath := filepath.Join(configPath, BotConfigFile)
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return c, err
	}
	defer configFile.Close()

	b, err := ioutil.ReadAll(configFile)
	if err != nil {
		return c, err
	}

	if len(b) != 0 {
		err := yaml.Unmarshal(b, c)
		if err != nil {
			return c, err
		}
	}
	return c, nil
}
