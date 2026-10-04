package components

import (
	"github.com/go-raptor/middlewares/csrf"
	"github.com/go-raptor/middlewares/limiter"
	"github.com/go-raptor/middlewares/logger"
	"github.com/go-raptor/middlewares/requestid"
	"github.com/go-raptor/middlewares/secure"
	"github.com/go-raptor/raptor/v4"
)

// Middlewares run in order, the first outermost: the request id exists before
// anything logs, and secure's headers reach csrf's and the limiter's errors too.
// There is no auth: the API is public and read-only.
func Middlewares() raptor.Middlewares {
	return raptor.Middlewares{
		raptor.Use(&requestid.RequestIDMiddleware{}),
		raptor.Use(&logger.LoggerMiddleware{}),
		// The defaults (no framing, COOP same-origin, nosniff, a referrer policy; HSTS
		// once app.secure_hsts_max_age is set), plus: only this origin may ask for
		// the reader's location.
		raptor.Use(secure.NewSecureMiddleware(secure.SecureConfig{Headers: map[string]string{
			"Permissions-Policy": "geolocation=(self), camera=(), microphone=()",
		}})),
		raptor.Use(&csrf.CSRFMiddleware{}),
		// A cache miss calls Open-Meteo, Meteoalarm or GeoNames: 20 requests a second
		// per client IP keeps one client from draining the upstream quotas.
		raptor.UseOnly(limiter.NewRateLimiterMiddleware(limiter.RateLimiterConfig{}), "Forecast"),
	}
}
