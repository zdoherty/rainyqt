package point

import (
	"errors"
	"fmt"
	"strconv"

	geo "github.com/kellydunn/golang-geo"
)

var (
	UnitParseError = errors.New("failed to parse unit")
	aliases        = map[string]DistanceUnit{
		"m":          Meters,
		"meter":      Meters,
		"meters":     Meters,
		"km":         Kilometers,
		"kilometer":  Kilometers,
		"kilometers": Kilometers,
		"mi":         Miles,
		"mile":       Miles,
		"miles":      Miles,
	}
)

type Point struct {
	*geo.Point
	Name string
}

func (p Point) String() string {
	return fmt.Sprintf("%s (%f, %f)", p.Name, p.Lat(), p.Lng())
}

type DistanceUnit int

const (
	Miles DistanceUnit = iota
	Meters
	Kilometers
)

func ParseUnit(name string) (DistanceUnit, error) {
	if du, ok := aliases[name]; ok {
		return du, nil
	}
	return 0, UnitParseError
}

func (du DistanceUnit) String() string {
	return [...]string{"Miles", "Meters", "Kilometers"}[du]
}

type Distance struct {
	N    int
	Unit DistanceUnit
}

func (d Distance) String() string {
	return fmt.Sprintf("%d %s", d.N, d.Unit)
}

func ParsePoint(name string, lat string, lng string) (Point, error) {
	var p Point
	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return p, err
	}
	lngFloat, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		return p, err
	}
	p = NewPoint("", latFloat, lngFloat)
	return p, nil
}

func NewPoint(name string, lat float64, lng float64) Point {
	p := Point{
		Name:  name,
		Point: geo.NewPoint(lat, lng),
	}
	return p
}
