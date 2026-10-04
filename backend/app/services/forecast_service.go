package services

import (
	"context"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/goopenmeteo"
	"github.com/h00s/weather/app/models"
)

const (
	// forecastTTL keeps "now" current without asking Open-Meteo twice in a few
	// minutes for the same square kilometre.
	forecastTTL = 10 * time.Minute
	// maxCachedLocations bounds every per-location cache (cells of ~1 km).
	maxCachedLocations = 5000
)

// ForecastService fetches forecasts from Open-Meteo, cached per ~1 km cell.
// APP_OPENMETEO_URL overrides the API root (tests).
type ForecastService struct {
	raptor.Service

	Upstream *UpstreamService

	client *goopenmeteo.OpenMeteo
	cache  *keyedCache[string, *goopenmeteo.Forecast]
}

func (s *ForecastService) Setup() error {
	s.client = goopenmeteo.NewOpenMeteo()
	s.client.BaseURL = s.Config.AppString("openmeteo_url", goopenmeteo.BaseURL)
	s.cache = newKeyedCache[string, *goopenmeteo.Forecast](forecastTTL, maxCachedLocations)
	return nil
}

// Forecast returns the forecast for at's cell and when it was fetched. When
// Open-Meteo fails it returns the last forecast it has for the cell; it errors
// only when it has none.
func (s *ForecastService) Forecast(ctx context.Context, at models.Coordinates) (*goopenmeteo.Forecast, time.Time, error) {
	at = at.Rounded()
	return s.cache.Get(ctx, at.Key(), func(ctx context.Context) (*goopenmeteo.Forecast, error) {
		if err := s.Upstream.Take(); err != nil {
			return nil, err
		}
		forecast, err := s.client.ForecastContext(ctx, goopenmeteo.ForecastOptions{
			Latitude:     at.Latitude,
			Longitude:    at.Longitude,
			Timezone:     "auto", // the location's own zone
			Current:      models.ForecastCurrentVariables,
			Hourly:       models.ForecastHourlyVariables,
			Daily:        models.ForecastDailyVariables,
			ForecastDays: models.ForecastDays,
		})
		if err != nil {
			s.Log.Warn("Open-Meteo forecast request failed", "cell", at.Key(), "error", err)
		}
		return forecast, err
	})
}
