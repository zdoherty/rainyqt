package point

import (
	"strconv"

	geo "github.com/kellydunn/golang-geo"
)

type Point struct {
	*geo.Point
	Name string
}

type DistanceUnit int

const (
	Miles DistanceUnit = iota
	Meters
	Kilometers
)

func (du DistanceUnit) String() string {
	return [...]string{"Miles", "Meters", "Kilometers"}[du]
}

type Distance struct {
	N    int
	Unit DistanceUnit
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
