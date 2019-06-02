package darksky

import (
	"time"

	"github.com/zdoherty/rainyqt/pkg/forecast"

	"github.com/shawntoffel/darksky"
	"github.com/zdoherty/rainyqt/pkg/config"
)

var (
	RainTypes = map[string]forecast.RainType{
		"rain":  forecast.Rain,
		"sleet": forecast.Sleet,
		"snow":  forecast.Snow,
	}
)

type Client struct {
	Client darksky.DarkSky
}

func NewClientFromConfig(c config.Config) forecast.Forecaster {
	cli := Client{
		Client: darksky.New(c.APIKey),
	}
	return cli
}

func (c Client) Get(location forecast.LatLong) (forecast.Forecast, error) {
	f := forecast.Forecast{
		Location: location,
		Fetched:  time.Now(),
	}

	resp, err := c.Client.Forecast(darksky.ForecastRequest{
		Latitude:  darksky.Measurement(location.Latitude),
		Longitude: darksky.Measurement(location.Longitude),
		Options: darksky.ForecastRequestOptions{
			Exclude: "minutely,weekly",
		},
	})
	if err != nil {
		return f, err
	}

	hours := make([]forecast.HourlyForecast, len(resp.Hourly.Data))
	for i, dp := range resp.Hourly.Data {
		hours[i] = newHourlyFromDataPoint(dp)
	}
	f.Hourly = hours

	days := make([]forecast.DailyForecast, len(resp.Daily.Data))
	for i, dp := range resp.Daily.Data {
		days[i] = newDailyFromDataPoint(dp)
	}
	f.Daily = days

	return f, nil
}

func newRainFromDataPoint(dp darksky.DataPoint) forecast.RainData {
	rd := forecast.RainData{
		Type:        forecast.Unknown,
		Inches:      float64(dp.PrecipIntensity),
		Probability: float64(dp.PrecipProbability),
	}
	if rt, ok := RainTypes[dp.PrecipType]; ok {
		rd.Type = rt
	}
	return rd
}

func newWindFromDataPoint(dp darksky.DataPoint) forecast.WindData {
	return forecast.WindData{
		GustSpeed: float64(dp.WindGust),
		Speed:     float64(dp.WindSpeed),
		Bearing:   float64(dp.WindBearing),
	}
}

func newDailyFromDataPoint(dp darksky.DataPoint) forecast.DailyForecast {
	fc := forecast.DailyForecast{
		Time:     time.Unix(int64(dp.Time), 0),
		HighTemp: float64(dp.TemperatureHigh),
		LowTemp:  float64(dp.TemperatureLow),
		Rain:     newRainFromDataPoint(dp),
		Wind:     newWindFromDataPoint(dp),
	}
	return fc
}

func newHourlyFromDataPoint(dp darksky.DataPoint) forecast.HourlyForecast {
	fc := forecast.HourlyForecast{
		Time: time.Unix(int64(dp.Time), 0),
		Temp: float64(dp.Temperature),
		Rain: newRainFromDataPoint(dp),
		Wind: newWindFromDataPoint(dp),
	}
	return fc
}
