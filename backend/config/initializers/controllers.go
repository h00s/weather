package initializers

import (
	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/weather/app/controllers"
)

func Controllers() raptor.Controllers {
	return raptor.Controllers{
		&controllers.ForecastController{},
	}
}
