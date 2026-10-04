package services

import (
	"errors"
	"fmt"

	"github.com/go-raptor/raptor/v4"
	"golang.org/x/time/rate"
)

// errUpstreamBudget is returned instead of calling Open-Meteo once this
// minute's budget is spent; the caches then serve what they have.
var errUpstreamBudget = errors.New("open-meteo: this server's request budget is spent")

// UpstreamService budgets the calls this server makes to Open-Meteo (forecast,
// air quality, geocoding). Its free tier allows 600 calls a minute and counts a
// forecast with many variables as several, so one budget for every client keeps
// any client, however many addresses it uses, from tripping that limit and
// leaving every uncached location failing. APP_OPENMETEO_PER_MINUTE sets it
// (default 120). The daily quota (10,000 calls on the free tier) is a deployment
// decision this budget does not enforce.
type UpstreamService struct {
	raptor.Service

	limiter *rate.Limiter
}

func (s *UpstreamService) Setup() error {
	perMinute, err := s.Config.AppInt("openmeteo_per_minute", 120)
	if err != nil {
		return err
	}
	if perMinute < 1 {
		return fmt.Errorf("app config openmeteo_per_minute: %d, want at least 1", perMinute)
	}
	s.limiter = rate.NewLimiter(rate.Limit(float64(perMinute)/60), max(1, perMinute/3))
	return nil
}

// Take spends one call, or returns errUpstreamBudget when none is left.
func (s *UpstreamService) Take() error {
	if !s.limiter.Allow() {
		return errUpstreamBudget
	}
	return nil
}
