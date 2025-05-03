package services

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/goopenmeteo"
)

type ForecastService struct {
	raptor.Service

	OpenMeteo *goopenmeteo.OpenMeteo
}

func NewForecastService() *ForecastService {
	return &ForecastService{
		OpenMeteo: goopenmeteo.NewOpenMeteo(),
	}
}

func (fs *ForecastService) GetForecast(latitude, longitude float64) (*goopenmeteo.Forecast, error) {
	return fs.OpenMeteo.Forecast(goopenmeteo.ForecastOptions{
		Latitude:  latitude,
		Longitude: longitude,
		Current: goopenmeteo.WeatherVariables{
			goopenmeteo.Temperature2M,
			goopenmeteo.RelativeHumidity2M,
			goopenmeteo.Precipitation,
			goopenmeteo.WeatherCode,
		},
		ForecastDays: 7,
		Daily: goopenmeteo.WeatherVariables{
			goopenmeteo.Temperature2MMax,
			goopenmeteo.Temperature2MMin,
			goopenmeteo.PrecipitationSum,
			goopenmeteo.WeatherCode,
		},
		ForecastHours: 24,
		Hourly: goopenmeteo.WeatherVariables{
			goopenmeteo.Temperature2M,
			goopenmeteo.RelativeHumidity2M,
			goopenmeteo.Precipitation,
			goopenmeteo.WeatherCode,
		},
	})
}
