package main

import "github.com/urfave/cli"

var (
	pointCmd = cli.Command{
		Name:     "point",
		Aliases:  []string{"p"},
		Usage:    "get, put, and delete points of interest",
		Category: "forecasting",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "provider, p",
				Usage:  "which forecast provider to use",
				EnvVar: "POINT_PROVIDER",
				Value:  "redis",
			},
		},
		Subcommands: []cli.Command{
			{
				Name:      "put",
				Usage:     "add or update a point",
				ArgsUsage: "<name> <lat> <lng>",
				Action:    pointPutHandler,
			},
			{
				Name:      "get",
				Usage:     "find a point by name",
				ArgsUsage: "<name>",
				Action:    pointGetHandler,
			},
			{
				Name:      "get-radius",
				Usage:     "find points within a radius of a point",
				ArgsUsage: "<radius> ( <name> | <lat> <lng> )",
				Action:    pointGetRadiusHandler,
			},
			{
				Name:   "ls",
				Usage:  "list all points",
				Action: pointListHandler,
			},
			{
				Name:      "rm",
				Usage:     "delete a named point",
				ArgsUsage: "<name>",
				Action:    pointDeleteHandler,
			},
		},
	}
)

func pointPutHandler(c *cli.Context) error {
	return nil
}

func pointGetHandler(c *cli.Context) error {
	return nil
}

func pointGetRadiusHandler(c *cli.Context) error {
	return nil
}

func pointListHandler(c *cli.Context) error {
	return nil
}

func pointDeleteHandler(c *cli.Context) error {
	return nil
}
