package main

import (
	"encoding/json"
	"os"

	"github.com/urfave/cli"
	"github.com/zdoherty/rainyqt/pkg/forecast"
	"github.com/zdoherty/rainyqt/pkg/forecast/darksky"
	"github.com/zdoherty/rainyqt/pkg/point"
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
	p, err := point.ParsePoint("", c.Args()[0], c.Args()[1])

	// setup forecast client based on provider
	var forecaster forecast.Forecaster
	switch provider := c.String("provider"); provider {
	case "darksky":
		if key := c.String("darksky-key"); key != "" {
			Config.Forecaster.DarkSky.APIKey = key
		}
		forecaster = darksky.NewClientFromConfig(Config)
	default:
		return cli.NewExitError("unsupported provider: "+provider, 78)
	}

	// dump forecast to stdout
	fc, err := forecaster.Get(p)
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
