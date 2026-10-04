package controllers

import (
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/weather/app/models"
	"github.com/h00s/weather/app/services"
)

type WarningsController struct {
	raptor.Controller

	Warnings *services.WarningsService
	Places   *services.PlacesService
}

// Index lists DHMZ's warnings for the county at ?lat&lon, in force now or
// starting within 48 hours, most severe first. Outside Croatia it is empty.
func (c *WarningsController) Index(ctx *raptor.Context) error {
	at, err := coordinates(ctx)
	if err != nil {
		return err
	}
	place, ok := c.Places.Nearest(at)
	county, known := models.Counties[place.County]
	if !ok || !known {
		return ctx.Data([]models.WarningResponse{})
	}
	alerts, err := c.Warnings.Alerts(ctx.Request().Context())
	if err != nil {
		return errs.NewErrorBadGateway("Meteoalarm is unavailable")
	}
	return ctx.Data(models.ActiveWarnings(alerts, county.Areas, time.Now()))
}
