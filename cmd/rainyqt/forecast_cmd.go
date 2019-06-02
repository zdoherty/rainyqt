package main

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/urfave/cli"
	"github.com/zdoherty/rainyqt/pkg/forecast"
	"github.com/zdoherty/rainyqt/pkg/forecast/darksky"
)

var (
	forecastCmd = cli.Command{
		Name:      "forecast",
		Aliases:   []string{"fc"},
		Usage:     "fetches the forecast for a location",
		ArgsUsage: "<lat> <lng>",
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
	if c.NArg() != 2 {
		return cli.NewExitError("forecast needs two arguments: lat, lng", 64)
	}
	lat, err := strconv.ParseFloat(c.Args().Get(0), 64)
	if err != nil {
		return cli.NewExitError("error parsing lat: "+err.Error(), 65)
	}
	lng, err := strconv.ParseFloat(c.Args().Get(1), 64)
	if err != nil {
		return cli.NewExitError("error parsing lng: "+err.Error(), 65)
	}
	location := forecast.NewLatLong(lat, lng)

	// setup forecast client based on provider
	var forecaster forecast.Forecaster
	switch provider := c.String("provider"); provider {
	case "darksky":
		if key := c.String("darksky-key"); key != "" {
			Config.APIKey = key
		}
		forecaster = darksky.NewClientFromConfig(Config)
	default:
		return cli.NewExitError("unsupported provider: "+provider, 78)
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
