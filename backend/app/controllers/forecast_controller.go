package controllers

import (
	"strconv"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
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
		return errs.NewErrorBadRequest("Lat and Lon query parameters are required")
	}

	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return errs.NewErrorBadRequest("Invalid lat parameter")
	}

	longitude, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return errs.NewErrorBadRequest("Invalid lon parameter")
	}

	forecast, err := fc.Forecast.GetForecast(latitude, longitude)
	if err != nil {
		return errs.NewErrorBadRequest(err.Error())
	}

	return c.Data(forecast)
}
