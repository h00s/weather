package controllers

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/weather/app/models"
	"github.com/h00s/weather/app/services"
)

type PlacesController struct {
	raptor.Controller

	Places *services.PlacesService
}

// Nearest names the Croatian settlement at ?lat&lon; abroad or at sea it is a 404.
func (c *PlacesController) Nearest(ctx *raptor.Context) error {
	at, err := coordinates(ctx)
	if err != nil {
		return err
	}
	place, ok := c.Places.Nearest(at)
	if !ok {
		return errs.NewErrorNotFound("No Croatian place near these coordinates")
	}
	return ctx.Data(models.NewPlaceResponse(place))
}
