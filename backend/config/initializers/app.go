package initializers

import (
	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/weather/config"
)

func App(c *raptor.Config) *raptor.AppInitializer {
	return &raptor.AppInitializer{
		Routes:      config.Routes(),
		Services:    Services(c),
		Middlewares: Middlewares(),
		Controllers: Controllers(),
	}
}
