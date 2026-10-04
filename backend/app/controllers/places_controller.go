package controllers

import (
	"strings"
	"unicode/utf8"

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

// Index finds settlements by name: ?q= of 2 to 60 characters, best match first.
func (c *PlacesController) Index(ctx *raptor.Context) error {
	q := strings.TrimSpace(ctx.QueryParam("q"))
	if n := utf8.RuneCountInString(q); n < 2 || n > 60 {
		return errs.NewErrorBadRequest("q: want 2 to 60 characters")
	}
	places, err := c.Places.Search(ctx.Request().Context(), q)
	if err != nil {
		return errs.NewErrorBadGateway("Open-Meteo geocoding is unavailable")
	}
	return ctx.Data(places)
}
