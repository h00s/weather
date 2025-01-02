package config

import "github.com/go-raptor/raptor/v3"

func Routes() raptor.Routes {
	return raptor.CollectRoutes(
		raptor.Scope("/api/v1",
			raptor.Scope("/forecast",
				// /api/v1/forecast?lat=1.23&lon=4.56
				raptor.Route("GET", "", "ForecastController", "GetForecast"),
			),
		),
	)
}
