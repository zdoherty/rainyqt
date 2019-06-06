package redis

import (
	"strconv"

	"github.com/go-redis/redis"
	"github.com/mmcloughlin/geohash"
	"github.com/zdoherty/rainyqt/pkg/config"
	"github.com/zdoherty/rainyqt/pkg/point"
)

var (
	DistanceUnits = map[point.DistanceUnit]string{
		point.Miles:      "mi",
		point.Meters:     "m",
		point.Kilometers: "km",
	}
)

type PointStore struct {
	Client     *redis.Client
	GeoSetName string
}

func NewPointStoreFromConfig(c config.RedisPointStoreConfig) PointStore {
	ps := PointStore{
		Client: redis.NewClient(&redis.Options{
			Addr:     c.Addr,
			Password: c.Password,
			DB:       c.DB,
		}),
		GeoSetName: c.GeoSetName,
	}
	return ps
}

func (ps PointStore) All() ([]point.Point, error) {
	set, err := ps.Client.ZRangeWithScores(ps.GeoSetName, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	points := make([]point.Point, len(set))
	for idx, item := range set {
		points[idx] = pointFromZ(item)
	}
	return points, nil
}

func (ps PointStore) GetByName(name string) (point.Point, error) {
	var p point.Point
	score, err := ps.Client.ZScore(ps.GeoSetName, name).Result()
	if err != nil {
		return p, err
	}
	lat, lng := latLngFromGeoHashFloat(score)
	return point.NewPoint(name, lat, lng), nil
}

func (ps PointStore) GetByRadius(p point.Point, radius point.Distance) ([]point.Point, error) {
	q := redis.GeoRadiusQuery{
		Radius:    float64(radius.N),
		Unit:      DistanceUnits[radius.Unit],
		WithCoord: true,
	}
	set, err := ps.Client.GeoRadius(ps.GeoSetName, p.Point.Lng(), p.Point.Lat(), q).Result()
	if err != nil {
		return nil, err
	}
	points := make([]point.Point, len(set))
	for idx, item := range set {
		points[idx] = point.NewPoint(item.Name, item.Latitude, item.Longitude)
	}
	return points, nil
}

func (ps PointStore) Put(p point.Point) error {
	_, err := ps.Client.GeoAdd(ps.GeoSetName, &redis.GeoLocation{
		Name:      p.Name,
		Longitude: p.Point.Lng(),
		Latitude:  p.Point.Lat(),
	}).Result()
	return err
}

func (ps PointStore) Delete(p point.Point) error {
	_, err := ps.Client.ZRem(ps.GeoSetName, p.Name).Result()
	return err
}

func latLngFromGeoHashFloat(f float64) (lat, lng float64) {
	return geohash.Decode(strconv.FormatFloat(f, 'f', 0, 64))
}

func pointFromZ(z redis.Z) point.Point {
	lat, lng := latLngFromGeoHashFloat(z.Score)
	return point.NewPoint(z.Member.(string), lat, lng)
}
