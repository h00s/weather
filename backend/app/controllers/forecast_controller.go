package controllers

import (
	"strconv"

	"github.com/go-raptor/errs"
	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/weather/app/services"
)

type ForecastController struct {
	raptor.Controller

	Forecast *services.ForecastService
}

func (fc *ForecastController) Get(c *raptor.Context) error {
	lat := c.QueryParam("lat")
	lon := c.QueryParam("lon")

	if lat == "" || lon == "" {
		return c.JSONError(errs.NewErrorBadRequest("lat and lon query parameters are required"))
	}

	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return c.JSONError(errs.NewErrorBadRequest("invalid lat parameter"))
	}

	longitude, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return c.JSONError(errs.NewErrorBadRequest("invalid lon parameter"))
	}

	forecast, err := fc.Forecast.GetForecast(latitude, longitude)
	if err != nil {
		return c.JSONError(errs.NewErrorBadRequest(err.Error()))
	}

	return c.JSONResponse(forecast)
}
