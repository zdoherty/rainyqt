package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/zdoherty/rainyqt/pkg/config"
	"github.com/zdoherty/rainyqt/pkg/version"
)

var (
	Build  string
	Config config.Config
)

func parseConfig(c *cli.Context) error {
	// set log level early for config debugging
	if cliLevel := c.GlobalString("log-level"); cliLevel != "" {
		lvl, err := log.ParseLevel(cliLevel)
		if err != nil {
			return err
		}
		log.SetLevel(lvl)
	}

	path := c.GlobalString("config")
	if path == "" {
		return nil
	}

	conf, err := config.NewConfigFromFile(path)
	if err != nil {
		return err
	}
	Config = conf

	// check cli log level again, fall back to config file level if unset
	if cliLevel := c.GlobalString("log-level"); cliLevel == "" {
		lvl, err := log.ParseLevel(conf.LogLevel)
		if err != nil {
			return err
		}
		log.SetLevel(lvl)
	}
	return nil
}

func main() {
	if Build != "" {
		version.RainyqtVersion.Build = Build
	}

	app := cli.NewApp()
	app.Name = "rainyqt"
	app.Version = version.RainyqtVersion.String()
	app.Authors = []cli.Author{
		{
			Name:  "Zack Doherty",
			Email: "me@zack.gg",
		},
	}
	app.Usage = "fun weather finder"
	app.Before = parseConfig

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Load configuration from `PATH`",
			EnvVar: "CONFIG_PATH",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "log-level, l",
			Usage:  "Set logging level (DEBUG, INFO, etc)",
			EnvVar: "LOG_LEVEL",
			Value:  "",
		},
	}

	app.Commands = []cli.Command{
		forecastCmd,
		pointCmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
