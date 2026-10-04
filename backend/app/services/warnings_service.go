package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/models"
)

const (
	meteoalarmURL  = "https://feeds.meteoalarm.org"
	meteoalarmFeed = "croatia" // DHMZ's warnings
	warningsTTL    = 15 * time.Minute
)

// WarningsService reads DHMZ's weather warnings from Meteoalarm's Croatian
// feed: one feed for the whole country, cached once. APP_METEOALARM_URL
// overrides the feed host (tests).
type WarningsService struct {
	raptor.Service

	client  *http.Client
	baseURL string
	cache   *staleCache[[]models.MeteoalarmAlert]
}

func (s *WarningsService) Setup() error {
	s.client = &http.Client{Timeout: 10 * time.Second}
	s.baseURL = strings.TrimSuffix(s.Config.AppString("meteoalarm_url", meteoalarmURL), "/")
	s.cache = newStaleCache[[]models.MeteoalarmAlert](warningsTTL)
	return nil
}

// Alerts returns every alert in the feed. When Meteoalarm fails it returns the
// last alerts it has; it errors only when it has none.
func (s *WarningsService) Alerts(ctx context.Context) ([]models.MeteoalarmAlert, error) {
	alerts, _, err := s.cache.Get(ctx, s.fetch)
	return alerts, err
}

func (s *WarningsService) fetch(ctx context.Context) ([]models.MeteoalarmAlert, error) {
	alerts, err := s.request(ctx)
	if err != nil {
		s.Log.Warn("Meteoalarm request failed", "error", err)
	}
	return alerts, err
}

func (s *WarningsService) request(ctx context.Context) ([]models.MeteoalarmAlert, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/api/v1/warnings/feeds-"+meteoalarmFeed, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meteoalarm: unexpected status %d", resp.StatusCode)
	}

	var feed models.MeteoalarmFeed
	if err := json.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("meteoalarm: %w", err)
	}
	alerts := make([]models.MeteoalarmAlert, len(feed.Warnings))
	for i, w := range feed.Warnings {
		alerts[i] = w.Alert
	}
	return alerts, nil
}
