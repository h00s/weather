package components

import (
	"github.com/go-raptor/controllers/spa/v2"
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/controllers"
)

func Controllers() raptor.Controllers {
	return raptor.Controllers{
		&controllers.ForecastController{},
		spa.NewSPAController(spa.SPAConfig{
			Directory: "public",
		}),
	}
}
