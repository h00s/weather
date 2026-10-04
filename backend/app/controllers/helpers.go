package controllers

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/weather/app/models"
)

// coordinates reads the lat and lon query parameters, rounded to the ~1 km cell
// that every upstream request and cache key uses. A missing or invalid one is a 400.
func coordinates(ctx *raptor.Context) (models.Coordinates, error) {
	at, err := models.ParseCoordinates(ctx.QueryParam("lat"), ctx.QueryParam("lon"))
	if err != nil {
		return models.Coordinates{}, errs.NewErrorBadRequest(err.Error())
	}
	return at.Rounded(), nil
}
