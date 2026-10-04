package controllers

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/weather/app/models"
	"github.com/h00s/weather/app/services"
)

type AirQualityController struct {
	raptor.Controller

	AirQuality *services.AirQualityService
}

// Show is today's European AQI and pollen at ?lat&lon.
func (c *AirQualityController) Show(ctx *raptor.Context) error {
	at, err := coordinates(ctx)
	if err != nil {
		return err
	}
	aq, fetchedAt, err := c.AirQuality.AirQuality(ctx.Request().Context(), at)
	if err != nil {
		return errs.NewErrorBadGateway("Open-Meteo air quality is unavailable")
	}
	return ctx.Data(models.NewAirQualityResponse(aq, fetchedAt))
}
