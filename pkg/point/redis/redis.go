package redis

import (
	"github.com/go-redis/redis"
	"github.com/zdoherty/rainyqt/pkg/config"
	"github.com/zdoherty/rainyqt/pkg/point"
)

var (
	distanceUnits = map[point.DistanceUnit]string{
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
	set, err := ps.Client.ZRange(ps.GeoSetName, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	points := make([]point.Point, len(set))
	for idx, name := range set {
		coords, err := ps.Client.GeoPos(ps.GeoSetName, name).Result()
		if err != nil {
			return nil, err
		}
		points[idx] = point.NewPoint(name, coords[0].Latitude, coords[0].Longitude)
	}
	return points, nil
}

func (ps PointStore) GetByName(name string) (point.Point, error) {
	var p point.Point
	coords, err := ps.Client.GeoPos(ps.GeoSetName, name).Result()
	if err != nil {
		return p, err
	}
	return point.NewPoint(name, coords[0].Latitude, coords[0].Longitude), nil
}

func (ps PointStore) GetByRadius(p point.Point, radius point.Distance) ([]point.Point, error) {
	q := &redis.GeoRadiusQuery{
		Radius:    float64(radius.N),
		Unit:      distanceUnits[radius.Unit],
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
