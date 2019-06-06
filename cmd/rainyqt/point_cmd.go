package main

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli"
	"github.com/zdoherty/rainyqt/pkg/point"
	"github.com/zdoherty/rainyqt/pkg/point/redis"
)

var (
	pointStore point.PointStore
	pointCmd   = cli.Command{
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
		Before: setupPointStore,
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
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "unit, u",
						Usage: "unit of radius (miles, meters, km)",
						Value: "miles",
					},
				},
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

func setupPointStore(c *cli.Context) error {
	switch provider := c.String("provider"); provider {
	case "redis":
		pointStore = redis.NewPointStoreFromConfig(Config.PointStore.Redis)
	default:
		return cli.NewExitError("unsupported point store provider: "+provider, 78)
	}
	return nil
}

func pointPutHandler(c *cli.Context) error {
	if nargs := c.NArg(); nargs != 3 {
		return cli.NewExitError("put needs 3 args, got "+string(nargs), 64)
	}
	args := c.Args()
	p, err := point.ParsePoint(args[0], args[1], args[2])
	if err != nil {
		return cli.NewExitError("error parsing point: "+err.Error(), 1)
	}
	err = pointStore.Put(p)
	if err != nil {
		return cli.NewExitError("error putting point: "+err.Error(), 1)
	}
	fmt.Println("successfully put point: " + p.String())
	return nil
}

func pointGetHandler(c *cli.Context) error {
	if nargs := c.NArg(); nargs != 1 {
		return cli.NewExitError("get needs 1 arg, got "+string(nargs), 64)
	}
	p, err := pointStore.GetByName(c.Args()[0])
	if err != nil {
		return cli.NewExitError("error getting point: "+err.Error(), 1)
	}
	fmt.Println(p.String())
	return nil
}

func pointGetRadiusHandler(c *cli.Context) error {
	unit, err := point.ParseUnit(c.String("unit"))
	if err != nil {
		return cli.NewExitError("couldn't parse radius unit: "+err.Error(), 1)
	}

	var p point.Point
	switch nargs := c.NArg(); nargs {
	case 3:
		p, err = point.ParsePoint("", c.Args()[1], c.Args()[2])
		if err != nil {
			return cli.NewExitError("couldn't parse point from lat lng args: "+err.Error(), 1)
		}
	case 2:
		p, err = pointStore.GetByName(c.Args()[1])
		if err != nil {
			return cli.NewExitError("couldn't get point by name: "+err.Error(), 1)
		}
	default:
		return cli.NewExitError("get-radius needs 2 or 3 args, got "+string(nargs), 64)
	}

	radius, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return cli.NewExitError("couldn't parse radius as int: "+err.Error(), 1)
	}

	ps, err := pointStore.GetByRadius(p, point.Distance{
		N:    radius,
		Unit: unit,
	})

	for idx, p := range ps {
		fmt.Printf("%d: %s\n", idx, p)
	}
	return nil
}

func pointListHandler(c *cli.Context) error {
	ps, err := pointStore.All()
	if err != nil {
		return cli.NewExitError("error listing points: "+err.Error(), 1)
	}
	for idx, p := range ps {
		fmt.Printf("%d: %s\n", idx, p)
	}
	return nil
}

func pointDeleteHandler(c *cli.Context) error {
	if nargs := c.NArg(); nargs != 1 {
		return cli.NewExitError("delete needs 1 arg, got "+string(nargs), 64)
	}

	p, err := pointStore.GetByName(c.Args()[0])
	if err != nil {
		return cli.NewExitError("couldn't get point by name: "+err.Error(), 1)
	}

	err = pointStore.Delete(p)
	if err != nil {
		return cli.NewExitError("couldn't delete point by name: "+err.Error(), 1)
	}

	return nil
}
