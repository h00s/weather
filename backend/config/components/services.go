package components

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/services"
)

func Services() raptor.Services {
	return raptor.Services{
		&services.ForecastService{},
	}
}
