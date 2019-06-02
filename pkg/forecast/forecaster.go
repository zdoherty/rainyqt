package forecast

import "time"

type LatLong struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewLatLong(lat float64, lng float64) LatLong {
	return LatLong{
		Latitude:  lat,
		Longitude: lng,
	}
}

type RainData struct {
	Type        RainType `json:"type"`
	Inches      float64  `json:"inches"`
	Probability float64  `json:"probability"`
}

type WindData struct {
	GustSpeed float64 `json:"gust_speed"`
	Speed     float64 `json:"speed"`
	Bearing   float64 `json:"bearing"`
}

type RainType int

const (
	Unknown RainType = iota
	Rain
	Sleet
	Snow
)

func (rt RainType) String() string {
	return [...]string{"Unknown", "Rain", "Sleet", "Snow"}[rt]
}

type HourlyForecast struct {
	Time time.Time `json:"time"`
	Temp float64   `json:"temp"`
	Rain RainData  `json:"rain"`
	Wind WindData  `json:"wind"`
}

type DailyForecast struct {
	Time     time.Time `json:"time"`
	HighTemp float64   `json:"high_temp"`
	LowTemp  float64   `json:"low_temp"`
	Rain     RainData  `json:"rain"`
	Wind     WindData  `json:"wind"`
}

type Forecast struct {
	Hourly   []HourlyForecast `json:"hourly"`
	Daily    []DailyForecast  `json:"daily"`
	Location LatLong          `json:"location"`
	Fetched  time.Time        `json:"fetched"`
}

type Forecaster interface {
	Get(location LatLong) (Forecast, error)
}
