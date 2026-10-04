package services

import (
	"context"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/goopenmeteo"
	"github.com/h00s/weather/app/models"
)

// airQualityTTL: the CAMS models behind it update a few times a day.
const airQualityTTL = 30 * time.Minute

// AirQualityService fetches today's air quality and pollen from Open-Meteo,
// cached per ~1 km cell. APP_AIRQUALITY_URL overrides the API root (tests).
type AirQualityService struct {
	raptor.Service

	Upstream *UpstreamService

	client *goopenmeteo.OpenMeteo
	cache  *keyedCache[string, *goopenmeteo.AirQuality]
}

func (s *AirQualityService) Setup() error {
	s.client = goopenmeteo.NewOpenMeteo()
	s.client.AirQualityURL = s.Config.AppString("airquality_url", goopenmeteo.AirQualityBaseURL)
	s.cache = newKeyedCache[string, *goopenmeteo.AirQuality](airQualityTTL, maxCachedLocations)
	return nil
}

// AirQuality returns today's air quality at at's cell and when it was fetched.
// When Open-Meteo fails it returns the last one it has; it errors only when it
// has none.
func (s *AirQualityService) AirQuality(ctx context.Context, at models.Coordinates) (*goopenmeteo.AirQuality, time.Time, error) {
	at = at.Rounded()
	return s.cache.Get(ctx, at.Key(), func(ctx context.Context) (*goopenmeteo.AirQuality, error) {
		if err := s.Upstream.Take(); err != nil {
			return nil, err
		}
		vars := models.AirQualityVariables()
		aq, err := s.client.AirQualityContext(ctx, goopenmeteo.AirQualityOptions{
			Latitude:     at.Latitude,
			Longitude:    at.Longitude,
			Timezone:     "auto",
			Current:      vars,
			Hourly:       vars,
			ForecastDays: 1,
		})
		if err != nil {
			s.Log.Warn("Open-Meteo air quality request failed", "cell", at.Key(), "error", err)
		}
		return aq, err
	})
}
