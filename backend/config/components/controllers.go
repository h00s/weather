package components

import (
	"github.com/go-raptor/controllers/spa"
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/controllers"
)

func Controllers() raptor.Controllers {
	return raptor.Controllers{
		&controllers.ForecastController{},
		spa.NewSPAController("public", "index.html"),
	}
}
