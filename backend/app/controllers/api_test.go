package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-raptor/raptor/v4"
)

// Probes must not depend on upstreams: an Open-Meteo outage must never restart
// or unroute the app.
func TestProbesAnswer200(t *testing.T) {
	app, u := newApp(t)
	u.fail(http.StatusServiceUnavailable)

	for _, path := range []string{"/healthz", "/readyz"} {
		body := raptor.DecodeJSON[map[string]string](t, app.TestGet(path), http.StatusOK)
		if body["status"] != "ok" {
			t.Errorf("GET %s body = %v", path, body)
		}
	}
}

func TestReadinessFailsOnceShutdownBegins(t *testing.T) {
	app, _ := newApp(t)
	app.Core.BeginShutdown()

	raptor.DecodeJSON[errorBody](t, app.TestGet("/readyz"), http.StatusServiceUnavailable)
	if rec := app.TestGet("/healthz"); rec.Code != http.StatusOK {
		t.Errorf("GET /healthz during shutdown = %d, want 200", rec.Code)
	}
}

// The SPA's catch-all would answer unknown API paths with an empty 404.
func TestUnknownAPIPathIsAJSON404(t *testing.T) {
	app, _ := newApp(t)

	body := raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/nope"), http.StatusNotFound)

	if body.Code != http.StatusNotFound {
		t.Errorf("body = %+v", body)
	}
}

func TestResponsesCarrySecurityHeadersAndARequestID(t *testing.T) {
	app, _ := newApp(t)

	want := map[string]string{
		"X-Content-Type-Options":     "nosniff",
		"Referrer-Policy":            "strict-origin-when-cross-origin",
		"X-Frame-Options":            "DENY",
		"Content-Security-Policy":    "frame-ancestors 'none'",
		"Cross-Origin-Opener-Policy": "same-origin",
		"Permissions-Policy":         "geolocation=(self), camera=(), microphone=()",
	}
	cases := []struct {
		name string
		rec  *httptest.ResponseRecorder
	}{
		{"probe", app.TestGet("/healthz")},
		{"JSON 404", app.TestGet("/api/v1/nope")},
		{"SPA navigation", app.TestGet("/", raptor.WithHeader("Accept", "text/html"))},
	}
	for _, c := range cases {
		for name, value := range want {
			if got := c.rec.Header().Get(name); got != value {
				t.Errorf("%s (%d): %s = %q, want %q", c.name, c.rec.Code, name, got, value)
			}
		}
		if got := c.rec.Header().Get("Strict-Transport-Security"); got != "" {
			t.Errorf("%s: HSTS = %q without secure_hsts_max_age", c.name, got)
		}
		if c.rec.Header().Get("X-Request-Id") == "" {
			t.Errorf("%s: no X-Request-Id", c.name)
		}
	}
}

// The API has no writes today; the csrf middleware still sits in the chain for
// the day it does.
func TestCrossSiteWriteIsRejected(t *testing.T) {
	app, _ := newApp(t)

	rec := app.TestPost("/api/v1/forecast", strings.NewReader(`{}`),
		raptor.WithHeader("Sec-Fetch-Site", "cross-site"),
		raptor.WithHeader("Origin", "https://evil.example"))

	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403", rec.Code)
	}
}

// One client gets a burst of 30 API requests, then 2 a second, so no script can
// churn through cells; another client is unaffected, and probes are never limited.
func TestAPIIsRateLimitedPerClient(t *testing.T) {
	app, _ := newApp(t)
	client := newClient()

	limited := false
	for range 40 {
		if app.TestGet("/api/v1/forecast", client).Code == http.StatusTooManyRequests {
			limited = true
			break
		}
	}
	if !limited {
		t.Fatal("40 requests in a burst were never limited")
	}
	if rec := app.TestGet("/api/v1/forecast", newClient()); rec.Code == http.StatusTooManyRequests {
		t.Error("another client was limited too")
	}
	for range 30 {
		if rec := app.TestGet("/healthz", client); rec.Code != http.StatusOK {
			t.Fatalf("GET /healthz = %d: probes must not be rate limited", rec.Code)
		}
	}
}

// All clients share one Open-Meteo budget (openmeteo_per_minute), so no client,
// however many addresses it uses, can trip the upstream's own limit and leave
// every uncached location failing. Cached cells keep answering.
func TestOpenMeteoCallsShareOneBudget(t *testing.T) {
	app, u := newAppWithConfig(t, map[string]string{"openmeteo_per_minute": "1"})

	raptor.DecodeJSON[map[string]any](t, app.TestGet("/api/v1/forecast?lat=45.59&lon=17.22", newClient()), http.StatusOK)
	raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/forecast?lat=43.51&lon=16.44", newClient()), http.StatusBadGateway)
	raptor.DecodeJSON[map[string]any](t, app.TestGet("/api/v1/forecast?lat=45.59&lon=17.22", newClient()), http.StatusOK)

	if n := u.hits("/forecast"); n != 1 {
		t.Errorf("Open-Meteo was asked %d times, want 1 (the budget)", n)
	}
}
