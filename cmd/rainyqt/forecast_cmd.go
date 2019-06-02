package main

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/zdoherty/rainyqt/pkg/forecast/darksky"

	"github.com/urfave/cli"
	"github.com/zdoherty/rainyqt/pkg/forecast"
)

var (
	forecastCmd = cli.Command{
		Name:      "forecast",
		Aliases:   []string{"fc"},
		Usage:     "fetches the forecast for a location",
		ArgsUsage: "latitude longitude",
		Category:  "forecasting",
		Action:    forecastHandler,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "provider, p",
				Usage:  "which forecast provider to use",
				EnvVar: "FORECAST_PROVIDER",
				Value:  "darksky",
			},
			cli.StringFlag{
				Name:   "darksky-key",
				Usage:  "darksky api key",
				EnvVar: "DARKSKY_API_KEY",
				Value:  "",
			},
		},
	}
)

func forecastHandler(c *cli.Context) error {
	// parse latlong from args
	lat, err := strconv.ParseFloat(c.Args().Get(0), 64)
	if err != nil {
		return err
	}
	lng, err := strconv.ParseFloat(c.Args().Get(1), 64)
	if err != nil {
		return err
	}
	location := forecast.NewLatLong(lat, lng)

	// setup forecast client based on provider
	var forecaster forecast.Forecaster
	switch c.String("provider") {
	case "darksky":
		if c.String("darksky-key") != "" {
			Config.APIKey = c.String("darksky-key")
		}
		forecaster = darksky.NewClientFromConfig(Config)
	default:
		return errors.New("unknown forecast provider: " + c.String("provider"))
	}

	// dump forecast to stdout
	fc, err := forecaster.Get(location)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(fc)
	if err != nil {
		return err
	}
	return nil
}
