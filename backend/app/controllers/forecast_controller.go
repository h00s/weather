package controllers

import (
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/weather/app/models"
	"github.com/h00s/weather/app/services"
)

type ForecastController struct {
	raptor.Controller

	Forecast *services.ForecastService
}

// Show is the forecast at ?lat&lon: now, every hour to the end of the week, and
// seven days from today.
func (c *ForecastController) Show(ctx *raptor.Context) error {
	at, err := coordinates(ctx)
	if err != nil {
		return err
	}
	forecast, fetchedAt, err := c.Forecast.Forecast(ctx.Request().Context(), at)
	if err != nil {
		return errs.NewErrorBadGateway("Open-Meteo is unavailable")
	}
	return ctx.Data(models.NewForecastResponse(forecast, time.Now(), fetchedAt))
}
