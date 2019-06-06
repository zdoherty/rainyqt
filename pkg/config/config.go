package config

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type RedisPointStoreConfig struct {
	GeoSetName string `yaml:"geo-set-name"`
	Addr       string `yaml:"addr"`
	Password   string `yaml:"password"`
	DB         int    `yaml:"db"`
}

type PointStoreConfig struct {
	Redis RedisPointStoreConfig `yaml:"redis"`
}

type DarkSkyForecasterConfig struct {
	APIKey string `yaml:"api-key"`
}

type ForecasterConfig struct {
	DarkSky DarkSkyForecasterConfig `yaml:"darksky"`
}

type Config struct {
	LogLevel   string           `yaml:"log-level"`
	Forecaster ForecasterConfig `yaml:"forecaster"`
	PointStore PointStoreConfig `yaml:"point-store"`
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
		LogLevel: "INFO",
		Forecaster: ForecasterConfig{
			DarkSky: DarkSkyForecasterConfig{
				APIKey: "",
			},
		},
		PointStore: PointStoreConfig{
			Redis: RedisPointStoreConfig{
				GeoSetName: "RAINYQT_POINTS",
				Addr:       "localhost:6379",
				Password:   "",
				DB:         0,
			},
		},
	}
	return c
}
