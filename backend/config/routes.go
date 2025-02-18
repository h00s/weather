package config

import "github.com/go-raptor/raptor/v3/router"

func Routes() router.Routes {
	return router.CollectRoutes(
		router.Scope("/api/v1",
			router.Scope("/forecast",
				// /api/v1/forecast?lat=1.23&lon=4.56
				router.Get("", "Forecast.GetForecast"),
			),
		),
	)
}
