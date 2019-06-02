package config

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	APIKey   string `yaml:"darksky-api-key"`
	LogLevel string `yaml:"log-level"`
}

func NewConfigFromFile(filename string) (Config, error) {
	c := NewConfigFromDefaults()

	log.WithField("filename", filename).Debug("loading config from file")
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func NewConfigFromDefaults() Config {
	c := Config{
		APIKey:   "",
		LogLevel: "INFO",
	}
	return c
}
