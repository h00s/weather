# vrijeme.app Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Bring vrijeme.app onto the current go-raptor ecosystem and rebuild it as a Croatian weather app for any location (GPS, search or one of 20 cities, remembered on the device). It gains DHMZ warnings, air quality and pollen, a "living blue sky" design and installation as a PWA.

**Architecture:** The Go/Raptor backend serves a read-only JSON API under `/api/v1`: forecast, air quality, warnings, place search and nearest place. Each upstream sits behind a stale-on-failure cache keyed by a ~1 km cell, and the same server serves the SvelteKit SPA from `public/`. The SvelteKit 2 + Svelte 5 SPA (adapter-static) keeps the chosen location and the last forecast in localStorage and polls while visible. It renders over a sky gradient computed from the location's sun times and the current weather.

**Tech Stack:**
- Backend: Go 1.27, raptor/v4 v4.6.1, spa/v2 v2.1.0, and the go-raptor middlewares requestid 1.0, logger 1.4, secure 1.0, csrf 1.1 and limiter 1.1. Upstream data comes through goopenmeteo v1.1.0.
- Frontend: Bun 1.4, SvelteKit 2.70, Svelte 5.57, Vite 8, TypeScript 6, Tailwind 4 and Vitest 5. UI pieces: @bybas/weather-icons (Meteocons), @lucide/svelte and Inter Variable.

**Spec:** `docs/superpowers/specs/2026-10-04-vrijeme-redesign-design.md`

## Global Constraints

**Repository**
- Work on branch `redesign`, and commit after each task. Each commit message ends with `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`. Never push, tag or merge to `production` without asking.
- Reference code lives in `~/dev/go/husak-dashboard`, written **dashboard:** below. Port exactly what each task says, never wholesale.

**Code conventions**
- The backend is module `github.com/h00s/weather` in `backend/`. Follow the raptor-api-conventions skill:
  - every DTO field gets an explicit camelCase json tag;
  - errors are `errs.*` values with English messages meant for developers;
  - app config is read with the typed getters;
  - request work gets the request's context.
- The frontend follows the sveltekit-conventions skill:
  - every network call goes component → service → client;
  - import through `$lib` only, with no custom aliases;
  - no `$env`, and the API base is the constant `"/api/v1"`;
  - props are destructured;
  - `bun run check` ends with 0 errors and 0 warnings.

**UI copy and formats**
- UI copy is Croatian only and hardcoded, and the page is `<html lang="hr">`.
- Numbers and dates go through `hr-HR` Intl. Times are shown in the forecast's IANA `timezone`, never the device's.

**Privacy, origin and caching**
- Coordinates are rounded to 2 decimals (~1 km) before they're stored on the device, sent to the API or used as a cache key.
- The SPA and the API share one origin: no CORS anywhere. Vite proxies `/api` with the object form `{ target: "http://localhost:3000" }`.
- Tests point each upstream at a stub through app config: `openmeteo_url`, `airquality_url`, `geocoding_url` and `meteoalarm_url`. In the environment these are `APP_OPENMETEO_URL`, `APP_AIRQUALITY_URL`, `APP_GEOCODING_URL` and `APP_METEOALARM_URL`.
- Backend cache TTLs: forecast 10 min, air quality 30 min, warnings 15 min, place search 24 h. A per-key cache holds at most 5000 keys (2000 for searches).
- Tracked YAML holds no secrets, and `max_body_bytes` is never 0.
- The footer credits Open-Meteo (CC BY 4.0), DHMZ via Meteoalarm, and GeoNames (CC BY 4.0).

## Review Focus

1. **A location outside Croatia** (GPS abroad, or a point at sea). The forecast and air quality still load, the name reads "Moja lokacija", and the warnings list is empty rather than an error. Pinned by:
   - Task 4: Sarajevo, the open Adriatic and Berlin get no place;
   - Task 6: `TestWarningsIndexIsEmptyOutsideCroatia`;
   - Task 10: `locationLabel` returns "Moja lokacija" for a null name.
2. **A week that spans a DST change.** Hourly times must keep the right offset after the switch; otherwise every hour label after it is off by one. Pinned by Task 3, `TestNewForecastResponseUsesTheZonesOwnOffsetAcrossDST`.
3. **Device storage that's blocked, missing or holds an older or corrupt shape** (Safari private mode, a previous app version). The app starts normally on the welcome screen. Pinned by Task 10's storage tests (throwing `getItem`/`setItem`, missing `localStorage`, invalid JSON, a failed validator) and its `isSavedLocation` tests.
4. **Geolocation denied, unavailable or timing out, or a browser without `navigator.geolocation`.** The reader gets a Croatian message and the picker stays usable. Pinned by Task 10's `positionErrorMessage` tests and the manual checks in Task 12.
5. **An upstream that's down.** A warnings or air-quality outage must not blank the forecast, and a forecast outage with a cached copy shows that copy marked stale. Pinned by:
   - Tasks 3, 6 and 7: each endpoint has its own 502 test;
   - Task 11: the store applies each `Promise.allSettled` result on its own (reviewed in code);
   - Task 14: the offline check.

---

## File Structure

**Backend (`backend/`)**

| File | Responsibility |
|---|---|
| `main.go` | Entry point; also embeds `time/tzdata` |
| `config/routes.yaml` | Every route |
| `config/components/{controllers,services,middlewares}.go` | Registration, with the middleware order and scopes |
| `.raptor.dev.yaml`, `.raptor.test.yaml`, `.gitignore` | Local and test config; `public/` is ignored |
| `cmd/genplaces/main.go` | Builds the places dataset from GeoNames |
| `app/models/coordinates.go` | `Coordinates`, `ParseCoordinates`, `Rounded`, `Key` |
| `app/models/forecast.go` | Forecast DTOs and `NewForecastResponse` |
| `app/models/place.go` | `Place`, `ParsePlaces`, `NearestPlace`, `DistanceKm`, `PlaceResponse` |
| `app/models/county.go` | The 21 counties and their Meteoalarm area names |
| `app/models/geocoding.go` | `GeocodingResponse`, `SearchResults` |
| `app/models/warning.go` (+ `testdata/meteoalarm-croatia.json`) | Ported from dashboard: Meteoalarm alerts → `WarningResponse` |
| `app/models/airquality.go` | Ported from dashboard: EAQI and pollen |
| `app/services/cache.go` | `staleCache` (ported) and the new `keyedCache` |
| `app/services/{forecast,places,warnings,airquality}_service.go` | Upstream clients behind the caches |
| `app/services/data/hr-places.tsv` | The embedded GeoNames dataset |
| `app/controllers/helpers.go` | `coordinates(ctx)` |
| `app/controllers/{forecast,places,warnings,airquality}_controller.go` | Actions |
| `app/controllers/{setup,fixtures,api,forecast,places,warnings,airquality}_test.go` | Integration tests against httptest upstreams |

Every model and service file has a `_test.go` beside it.

**Frontend (`frontend/`)**

| File | Responsibility |
|---|---|
| `package.json`, `vite.config.ts`, `tsconfig.json` | Dependencies and config; kit config lives in `vite.config.ts` |
| `src/app.html` | Shell: `lang="hr"`, the theme-color meta, icons and the manifest |
| `src/routes/layout.css` | Tailwind and the theme tokens, plus the `panel`, `text-lift` and `scrollbar-thin` utilities |
| `src/routes/{+layout.ts,+layout.svelte,+error.svelte,+page.svelte}` | SPA mode, the deploy reload, the error page and the single page |
| `src/lib/api/client.ts` | Ported from dashboard: the GET-only client |
| `src/lib/types/{forecast,place,location,warning,airquality}.ts` | Mirrors of the backend DTOs, plus `SavedLocation` |
| `src/lib/services/{forecast,places,warnings,airquality}.ts` | One function per endpoint |
| `src/lib/data/cities.ts` | The 20 cities |
| `src/lib/helpers/*.ts` | Pure logic, each with `*.test.ts`: format, weather, sparkline, alerts, preview, sky, storage, geo, location, wind, uv, air, svg, forecast |
| `src/lib/stores/{clock,location,weather}.svelte.ts` | Runes stores |
| `src/lib/components/*.svelte` | UI: the sky, picker, sheet, welcome screen, header, panels and tiles, plus pull-to-refresh |
| `src/service-worker.ts` | Precaches the app shell |
| `static/` | `manifest.json`, `favicon.svg`, `favicon.png`, `apple-touch-icon.png`, `icons/*.png` |
| `icons/` | `maskable.svg` and `render.sh`, which generates the PNGs |

**Root:** `Dockerfile` and `README.md`. `frontend/README.md` is deleted.

---

## Task 1: Backend ecosystem update and test harness

**Files:**
- Modify: `backend/go.mod`, `backend/go.sum`, `backend/config/components/middlewares.go`, `backend/config/components/controllers.go`, `backend/config/routes.yaml`, `backend/.gitignore`, `Dockerfile`
- Create: `backend/.raptor.test.yaml`
- Track (it exists but is git-ignored): `backend/.raptor.dev.yaml`
- Test: `backend/app/controllers/setup_test.go`, `backend/app/controllers/fixtures_test.go`, `backend/app/controllers/api_test.go`

**Interfaces:**
- Produces, in package `controllers_test`:
  - `newApp(t) (*raptor.Raptor, *upstreams)`;
  - `(*upstreams).fail(status int)`, `(*upstreams).hits(path string) int` and `(*upstreams).query(path string) url.Values`;
  - `canned(path string, now time.Time) ([]byte, bool)`, a `switch` that later tasks extend with their upstream's path;
  - `newClient() raptor.TestRequestOption` and `errorBody{Code, Message}`.
- Produces in `config/components/middlewares.go`: the limiter is scoped with `raptor.UseOnly(…, "Forecast")`. Tasks 4, 6 and 7 append their controller names to that list.

- [ ] **Step 1: Update the dependencies**

```bash
cd backend
go get github.com/go-raptor/raptor/v4@v4.6.1 github.com/go-raptor/controllers/spa/v2@v2.1.0 \
  github.com/go-raptor/middlewares/logger@v1.4.0 github.com/go-raptor/middlewares/requestid@v1.0.0 \
  github.com/go-raptor/middlewares/secure@v1.0.0 github.com/go-raptor/middlewares/csrf@v1.1.0 \
  github.com/go-raptor/middlewares/limiter@v1.1.0 github.com/h00s/goopenmeteo@v1.1.0
```

- [ ] **Step 2: Write the test harness**

`backend/app/controllers/setup_test.go`:

```go
package controllers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-raptor/raptor/v4"
	rconfig "github.com/go-raptor/raptor/v4/config"
	"github.com/h00s/weather/config"
	"github.com/h00s/weather/config/components"
)

// upstreams stands in for Open-Meteo (forecast, air quality, geocoding) and
// Meteoalarm on one httptest server. Each test builds its own app against its
// own upstreams, so cached answers never leak between tests.
type upstreams struct {
	srv *httptest.Server

	mu      sync.Mutex
	status  int                   // when set, every path answers it; 0 serves canned bodies
	counts  map[string]int        // requests per path
	queries map[string]url.Values // the last query per path
}

func newApp(t *testing.T) (*raptor.Raptor, *upstreams) {
	t.Helper()
	u := &upstreams{counts: map[string]int{}, queries: map[string]url.Values{}}
	u.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u.mu.Lock()
		u.counts[r.URL.Path]++
		u.queries[r.URL.Path] = r.URL.Query()
		status := u.status
		u.mu.Unlock()
		if status != 0 {
			w.WriteHeader(status)
			return
		}
		body, ok := canned(r.URL.Path, time.Now())
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}))
	t.Cleanup(u.srv.Close)

	app := raptor.NewTestApp(components.New(), config.Routes(), raptor.WithConfig(&rconfig.Config{
		AppConfig: map[string]string{
			"openmeteo_url":  u.srv.URL,
			"airquality_url": u.srv.URL,
			"geocoding_url":  u.srv.URL,
			"meteoalarm_url": u.srv.URL,
		},
	}))
	return app, u
}

// fail makes every upstream answer status from now on.
func (u *upstreams) fail(status int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.status = status
}

// hits counts the requests that reached path.
func (u *upstreams) hits(path string) int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.counts[path]
}

// query is the last query string sent to path.
func (u *upstreams) query(path string) url.Values {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.queries[path]
}

var clientIPs atomic.Uint32

// newClient gives a request its own client address, and so its own rate-limit
// bucket: every httptest request otherwise comes from 192.0.2.1.
func newClient() raptor.TestRequestOption {
	n := clientIPs.Add(1)
	return raptor.WithRemoteAddr(fmt.Sprintf("10.%d.%d.%d", byte(n>>16), byte(n>>8), byte(n)))
}

type errorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
```

`backend/app/controllers/fixtures_test.go`:

```go
package controllers_test

import "time"

// canned answers each upstream path a service requests. Every task that adds an
// upstream adds its path here; ok is false for a path no service should call.
func canned(path string, now time.Time) (body []byte, ok bool) {
	switch path {
	}
	return nil, false
}
```

- [ ] **Step 3: Write the failing API tests**

`backend/app/controllers/api_test.go`:

```go
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

// One client gets at most 20 API requests a second (burst 20), so nobody can
// drain the upstream quotas; another client is unaffected, and probes are never limited.
func TestAPIIsRateLimitedPerClient(t *testing.T) {
	app, _ := newApp(t)
	client := newClient()

	limited := false
	for range 30 {
		if app.TestGet("/api/v1/forecast", client).Code == http.StatusTooManyRequests {
			limited = true
			break
		}
	}
	if !limited {
		t.Fatal("30 requests in a burst were never limited")
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
```

- [ ] **Step 4: Run the tests and watch them fail**

Run: `cd backend && go test ./app/controllers/`

Expected: FAIL. `/healthz` is a 404 (no route), headers are missing, and nothing answers 429.

- [ ] **Step 5: Wire the middlewares, the SPA and the routes**

`backend/config/components/middlewares.go`:

```go
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
```

`backend/config/components/controllers.go`:

```go
package components

import (
	"github.com/go-raptor/controllers/spa/v2"
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/controllers"
)

func Controllers() raptor.Controllers {
	return raptor.Controllers{
		&controllers.ForecastController{},
		&spa.SPAController{}, // serves public/; spa_optional in .raptor.dev.yaml and .raptor.test.yaml
	}
}
```

`backend/config/routes.yaml`:

```yaml
routes:
  /healthz:
    GET: Health.Live            # built in: the process serves requests
  /readyz:
    GET: Health.Ready           # built in: 503 from the moment shutdown begins
  /api/v1:
    /forecast:
      GET: Forecast.Get
    /{path...}: Errors.NotFound # unknown API paths answer a JSON 404, not the SPA's empty one

  /: SPA.Index                  # catch-all for the frontend
```

`backend/.gitignore` stops ignoring the dev config. It also ignores a local copy of the frontend build, which Task 15 uses:

```
/bin
/public
```

`backend/.raptor.dev.yaml` (now tracked):

```yaml
# Local development only. Production is configured by environment variables (see the Dockerfile).
general:
  log_level: debug

server:
  address: "127.0.0.1"
  port: 3000

app:
  spa_optional: "true" # the API boots before the frontend is built; Vite serves it in dev
```

`backend/.raptor.test.yaml`:

```yaml
# Loaded by raptor.NewTestApp together with .raptor.yaml (if any). Tests point
# every upstream at httptest stubs via raptor.WithConfig.
app:
  spa_optional: "true" # tests have no frontend build
```

Then: `cd backend && go mod tidy`. This drops `middlewares/cors`.

- [ ] **Step 6: Run the tests and vet**

Run: `cd backend && go vet ./... && go test ./...`

Expected: PASS.

- [ ] **Step 7: Update the Dockerfile**

Replace `Dockerfile` with:

```dockerfile
FROM golang:alpine AS backend
WORKDIR /app
ENV CGO_ENABLED=0
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend ./
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/weather

FROM oven/bun:latest AS frontend
WORKDIR /app
COPY frontend/package.json frontend/bun.lock ./
RUN bun install --frozen-lockfile
COPY frontend ./
RUN bun run build

FROM gcr.io/distroless/static-debian13:latest
WORKDIR /app
COPY --from=backend /out/weather ./
COPY --from=frontend /app/build ./public
# Raptor binds 127.0.0.1 by default, which is unreachable from outside the container.
ENV SERVER_ADDRESS=0.0.0.0
# Behind the HTTPS reverse proxy: rate-limit and log the client's address, not the proxy's.
ENV SERVER_IP_EXTRACTOR=x-forwarded-for
# HSTS starts small; raise it to 31536000 once HTTPS is settled.
ENV APP_SECURE_HSTS_MAX_AGE=300

EXPOSE 3000

ENTRYPOINT ["./weather"]
```

- [ ] **Step 8: Commit**

```bash
git add backend Dockerfile
git commit -m "Update Raptor to v4.6.1 with requestid, secure, csrf and limiter

Health probes, JSON 404s for unknown API paths, the SPA via spa/v2.1 config,
no CORS (same origin), tracked dev and test configs, and an integration-test
harness with stubbed upstreams.

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 2: Coordinates and per-key caches

**Files:**
- Create: `backend/app/models/coordinates.go`, `backend/app/models/coordinates_test.go`
- Create, ported from dashboard and then extended: `backend/app/services/cache.go`, `backend/app/services/cache_test.go`
- Delete: `backend/app/models/.keep`

**Interfaces:**
- Produces:
  - `models.Coordinates{Latitude, Longitude float64}`;
  - `models.ParseCoordinates(lat, lon string) (models.Coordinates, error)`;
  - `(models.Coordinates).Rounded() models.Coordinates` and `(models.Coordinates).Key() string`, which returns something like `"45.59,17.23"`;
  - `newStaleCache[T](ttl)` and `(*staleCache[T]).Get(ctx, fetch) (T, time.Time, error)`;
  - `newKeyedCache[K, T](ttl, max)` and `(*keyedCache[K, T]).Get(ctx, key, fetch) (T, time.Time, error)`;
  - for tests: `fakeClock` and `counter(values...)`.

- [ ] **Step 1: Write the failing coordinate tests**

`backend/app/models/coordinates_test.go`:

```go
package models

import "testing"

func TestParseCoordinates(t *testing.T) {
	ok := []struct {
		lat, lon string
		want     Coordinates
	}{
		{"45.5936", "17.2251", Coordinates{45.5936, 17.2251}},
		{"-33.9", "151.2", Coordinates{-33.9, 151.2}},
		{"90", "-180", Coordinates{90, -180}},
	}
	for _, c := range ok {
		got, err := ParseCoordinates(c.lat, c.lon)
		if err != nil || got != c.want {
			t.Errorf("ParseCoordinates(%q, %q) = %v, %v; want %v", c.lat, c.lon, got, err, c.want)
		}
	}

	bad := []struct{ lat, lon string }{
		{"", "17"}, {"45", ""}, {"abc", "17"}, {"45", "1e999"}, {"NaN", "17"},
		{"Inf", "17"}, {"91", "17"}, {"45", "180.5"}, {"-90.01", "0"},
	}
	for _, c := range bad {
		if got, err := ParseCoordinates(c.lat, c.lon); err == nil {
			t.Errorf("ParseCoordinates(%q, %q) = %v, want an error", c.lat, c.lon, got)
		}
	}
}

func TestRoundedSnapsToAKilometreCell(t *testing.T) {
	got := Coordinates{45.5936, 17.2251}.Rounded()
	if got != (Coordinates{45.59, 17.23}) {
		t.Errorf("Rounded = %v, want {45.59 17.23}", got)
	}
	if k := got.Key(); k != "45.59,17.23" {
		t.Errorf("Key = %q", k)
	}
	if a, b := (Coordinates{45.5912, 17.2263}).Key(), (Coordinates{45.5936, 17.2251}).Key(); a != b {
		t.Errorf("neighbours in one cell got keys %q and %q", a, b)
	}
	if k := (Coordinates{-0.001, 0.004}).Key(); k != "0.00,0.00" {
		t.Errorf("Key = %q, want no negative zero", k)
	}
}
```

- [ ] **Step 2: Run it and watch it fail**

Run: `cd backend && go test ./app/models/`

Expected: FAIL with `undefined: ParseCoordinates`.

- [ ] **Step 3: Implement the coordinates**

`backend/app/models/coordinates.go`:

```go
package models

import (
	"fmt"
	"math"
	"strconv"
)

// Coordinates is a point on Earth in decimal degrees.
type Coordinates struct {
	Latitude  float64
	Longitude float64
}

// ParseCoordinates reads lat and lon as sent in a query string. Both must be
// finite numbers in range; the error names the one that isn't.
func ParseCoordinates(lat, lon string) (Coordinates, error) {
	latitude, err := parseDegrees(lat, 90)
	if err != nil {
		return Coordinates{}, fmt.Errorf("lat %q: %w", lat, err)
	}
	longitude, err := parseDegrees(lon, 180)
	if err != nil {
		return Coordinates{}, fmt.Errorf("lon %q: %w", lon, err)
	}
	return Coordinates{latitude, longitude}, nil
}

func parseDegrees(s string, limit float64) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(v) || math.Abs(v) > limit {
		return 0, fmt.Errorf("want a number from -%g to %g", limit, limit)
	}
	return v, nil
}

// Rounded snaps c to 2 decimals, a cell of about 1 km. Every upstream request
// and cache key uses it, so neighbours share one cached answer and no precise
// position is passed on.
func (c Coordinates) Rounded() Coordinates {
	return Coordinates{round2(c.Latitude), round2(c.Longitude)}
}

// Key names c's cell, such as "45.59,17.23".
func (c Coordinates) Key() string {
	r := c.Rounded()
	return strconv.FormatFloat(r.Latitude, 'f', 2, 64) + "," + strconv.FormatFloat(r.Longitude, 'f', 2, 64)
}

func round2(v float64) float64 {
	r := math.Round(v*100) / 100
	if r == 0 {
		return 0 // never "-0.00"
	}
	return r
}
```

Then: `git rm -q backend/app/models/.keep`

- [ ] **Step 4: Port the stale cache and write failing tests for the keyed cache**

```bash
cp ~/dev/go/husak-dashboard/backend/app/services/cache.go backend/app/services/cache.go
cp ~/dev/go/husak-dashboard/backend/app/services/cache_test.go backend/app/services/cache_test.go
```

Append to `backend/app/services/cache_test.go`:

```go
func newTestKeyedCache(ttl time.Duration, max int) (*keyedCache[string, int], *fakeClock) {
	clock := &fakeClock{now: time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)}
	c := newKeyedCache[string, int](ttl, max)
	c.now = clock.Now
	return c, clock
}

func TestKeyedCacheFetchesEachKeyOnce(t *testing.T) {
	c, _ := newTestKeyedCache(time.Minute, 10)
	fetchA, callsA := counter(1)
	fetchB, callsB := counter(2)

	for range 3 {
		if v, _, err := c.Get(t.Context(), "a", fetchA); v != 1 || err != nil {
			t.Fatalf("Get(a) = %d, %v", v, err)
		}
		if v, _, err := c.Get(t.Context(), "b", fetchB); v != 2 || err != nil {
			t.Fatalf("Get(b) = %d, %v", v, err)
		}
	}
	if *callsA != 1 || *callsB != 1 {
		t.Errorf("fetches: a %d, b %d; want one each", *callsA, *callsB)
	}
}

func TestKeyedCacheRefetchesAKeyAfterTTL(t *testing.T) {
	c, clock := newTestKeyedCache(time.Minute, 10)
	fetch, calls := counter(1, 2)

	c.Get(t.Context(), "a", fetch)
	clock.Advance(2 * time.Minute)
	if v, _, _ := c.Get(t.Context(), "a", fetch); v != 2 || *calls != 2 {
		t.Errorf("after TTL: value %d after %d fetches, want 2 after 2", v, *calls)
	}
}

func TestKeyedCacheDropsTheLeastRecentlyUsedKeyWhenFull(t *testing.T) {
	c, clock := newTestKeyedCache(time.Hour, 2)
	fetch, calls := counter(7)
	get := func(key string) {
		clock.Advance(time.Second)
		c.Get(t.Context(), key, fetch)
	}

	get("a")
	get("b")
	get("a") // b is now the least recently used
	get("c") // full: drops b
	before := *calls
	get("a")
	if *calls != before {
		t.Error("a was dropped, want b dropped")
	}
	get("b")
	if *calls != before+1 {
		t.Error("b was kept, want it dropped")
	}
}
```

- [ ] **Step 5: Run them and watch them fail**

Run: `cd backend && go test ./app/services/`

Expected: FAIL with `undefined: newKeyedCache`.

- [ ] **Step 6: Implement the keyed cache**

Append to `backend/app/services/cache.go`:

```go
// keyedCache keeps a staleCache per key, for upstreams asked per location or
// per query. It holds at most max keys: adding one more drops the least
// recently used. Callers with different keys never wait on each other.
type keyedCache[K comparable, T any] struct {
	ttl time.Duration
	max int
	now func() time.Time

	mu      sync.Mutex
	entries map[K]*keyedEntry[T]
}

type keyedEntry[T any] struct {
	cache    *staleCache[T]
	lastUsed time.Time
}

func newKeyedCache[K comparable, T any](ttl time.Duration, max int) *keyedCache[K, T] {
	return &keyedCache[K, T]{ttl: ttl, max: max, now: time.Now, entries: map[K]*keyedEntry[T]{}}
}

// Get returns key's value as staleCache.Get does.
func (c *keyedCache[K, T]) Get(ctx context.Context, key K, fetch func(context.Context) (T, error)) (T, time.Time, error) {
	return c.entry(key).Get(ctx, fetch)
}

func (c *keyedCache[K, T]) entry(key K) *staleCache[T] {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		if len(c.entries) >= c.max {
			c.evictLeastRecentlyUsed()
		}
		e = &keyedEntry[T]{cache: &staleCache[T]{ttl: c.ttl, now: c.now}}
		c.entries[key] = e
	}
	e.lastUsed = c.now()
	return e.cache
}

func (c *keyedCache[K, T]) evictLeastRecentlyUsed() {
	var oldest K
	var oldestAt time.Time
	first := true
	for k, e := range c.entries {
		if first || e.lastUsed.Before(oldestAt) {
			oldest, oldestAt, first = k, e.lastUsed, false
		}
	}
	delete(c.entries, oldest)
}
```

- [ ] **Step 7: Run all backend tests**

Run: `cd backend && go vet ./... && go test ./...`

Expected: PASS.

- [ ] **Step 8: Commit**

```bash
git add -A backend/app/models backend/app/services
git commit -m "Add coordinates parsing and a bounded per-key stale cache

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 3: Forecast endpoint

**Files:**
- Create: `backend/app/models/forecast.go`, `backend/app/models/forecast_test.go`, `backend/app/controllers/helpers.go`, `backend/app/controllers/forecast_test.go`
- Rewrite: `backend/app/services/forecast_service.go`, `backend/app/controllers/forecast_controller.go`
- Modify: `backend/config/routes.yaml`, `backend/main.go`, `backend/app/controllers/fixtures_test.go`

**Interfaces:**
- Consumes: `models.Coordinates`, `newKeyedCache` (Task 2), `newApp` and `canned` (Task 1).
- Produces:
  - `models.ForecastResponse` and its parts, `models.NewForecastResponse(f, now, fetchedAt)` and `models.ForecastDays = 7`;
  - `services.ForecastService.Forecast(ctx, at models.Coordinates) (*goopenmeteo.Forecast, time.Time, error)`;
  - the constant `maxCachedLocations = 5000` in package services;
  - `coordinates(ctx *raptor.Context) (models.Coordinates, error)` in package controllers, which returns rounded coordinates or a 400;
  - for tests in package models, `var zagreb = time.FixedZone("GMT+2", 7200)` (Task 6 reuses it).
- The frontend's `types/forecast.ts` (Task 8) mirrors these json tags exactly.

- [ ] **Step 1: Write the failing model tests**

`backend/app/models/forecast_test.go`:

```go
package models

import (
	"testing"
	"time"

	"github.com/h00s/goopenmeteo"
)

// zagreb is how goopenmeteo places a CEST response: a fixed zone at today's offset.
var zagreb = time.FixedZone("GMT+2", 7200)

func hours(start time.Time, n int) []time.Time {
	times := make([]time.Time, n)
	for i := range times {
		times[i] = start.Add(time.Duration(i) * time.Hour)
	}
	return times
}

func seq(n int, from float64) []float64 {
	values := make([]float64, n)
	for i := range values {
		values[i] = from + float64(i)
	}
	return values
}

func testForecast() *goopenmeteo.Forecast {
	midnight := time.Date(2026, 9, 29, 0, 0, 0, 0, zagreb)
	days := []time.Time{midnight, midnight.AddDate(0, 0, 1), midnight.AddDate(0, 0, 2)}
	return &goopenmeteo.Forecast{
		Timezone:  "Europe/Zagreb",
		Elevation: 161,
		Current: goopenmeteo.Current{
			Time: time.Date(2026, 9, 29, 14, 15, 0, 0, zagreb),
			Data: map[string]float64{
				goopenmeteo.Temperature2M:       17.4,
				goopenmeteo.ApparentTemperature: 15.9,
				goopenmeteo.RelativeHumidity2M:  71,
				goopenmeteo.DewPoint2M:          12.1,
				goopenmeteo.Precipitation:       0.2,
				goopenmeteo.WeatherCode:         61,
				goopenmeteo.CloudCover:          90,
				goopenmeteo.PressureMsl:         1013.2,
				goopenmeteo.WindSpeed10M:        12.5,
				goopenmeteo.WindDirection10M:    225,
				goopenmeteo.WindGusts10M:        31,
				goopenmeteo.Visibility:          24000,
				goopenmeteo.UVIndex:             2.6,
				goopenmeteo.IsDay:               1,
			},
		},
		Hourly: goopenmeteo.TimeseriesData{
			Time: hours(midnight, 72),
			Data: map[string][]float64{
				goopenmeteo.Temperature2M:            seq(72, 0),
				goopenmeteo.ApparentTemperature:      seq(72, -2),
				goopenmeteo.PrecipitationProbability: seq(72, 100),
				goopenmeteo.Precipitation:            seq(72, 200),
				goopenmeteo.WeatherCode:              seq(72, 0),
				goopenmeteo.WindSpeed10M:             seq(72, 300),
				goopenmeteo.WindDirection10M:         seq(72, 400),
				goopenmeteo.IsDay:                    seq(72, 0), // 0 at midnight, then 1 and up: "day"
			},
		},
		Daily: goopenmeteo.TimeseriesData{
			Time: days,
			Data: map[string][]float64{
				goopenmeteo.Temperature2MMin:            {9, 8, 7},
				goopenmeteo.Temperature2MMax:            {19, 18, 17},
				goopenmeteo.WeatherCode:                 {61, 3, 0},
				goopenmeteo.PrecipitationSum:            {4.2, 0, 0},
				goopenmeteo.PrecipitationProbabilityMax: {80, 10, 0},
				goopenmeteo.WindSpeed10MMax:             {25, 14, 9},
				goopenmeteo.WindGusts10MMax:             {48, 30, 20},
				goopenmeteo.WindDirection10MDominant:    {200, 45, 90},
				goopenmeteo.UVIndexMax:                  {3.1, 4, 5},
				goopenmeteo.DaylightDuration:            {42600, 42400, 42200},
			},
			Times: map[string][]time.Time{
				goopenmeteo.Sunrise: {days[0].Add(6*time.Hour + 46*time.Minute), days[1].Add(6*time.Hour + 47*time.Minute), days[2].Add(6*time.Hour + 49*time.Minute)},
				goopenmeteo.Sunset:  {days[0].Add(18*time.Hour + 35*time.Minute), days[1].Add(18*time.Hour + 33*time.Minute), days[2].Add(18*time.Hour + 31*time.Minute)},
			},
		},
	}
}

func TestNewForecastResponseMapsCurrentConditions(t *testing.T) {
	now := time.Date(2026, 9, 29, 14, 20, 0, 0, zagreb)
	fetchedAt := now.Add(-time.Minute)

	res := NewForecastResponse(testForecast(), now, fetchedAt)

	want := CurrentForecastResponse{
		Time:                time.Date(2026, 9, 29, 14, 15, 0, 0, zagreb),
		Temperature:         17.4,
		ApparentTemperature: 15.9,
		Humidity:            71,
		DewPoint:            12.1,
		Precipitation:       0.2,
		WeatherCode:         61,
		CloudCover:          90,
		Pressure:            1013.2,
		WindSpeed:           12.5,
		WindDirection:       225,
		WindGusts:           31,
		Visibility:          24000,
		UVIndex:             2.6,
		IsDay:               true,
	}
	got := res.Current
	if !got.Time.Equal(want.Time) {
		t.Errorf("current time = %v, want %v", got.Time, want.Time)
	}
	got.Time = want.Time
	if got != want {
		t.Errorf("current = %+v\nwant      %+v", got, want)
	}
	if res.Timezone != "Europe/Zagreb" || res.Elevation != 161 || !res.FetchedAt.Equal(fetchedAt) {
		t.Errorf("timezone %q, elevation %v, fetchedAt %v", res.Timezone, res.Elevation, res.FetchedAt)
	}
}

func TestNewForecastResponseHourlyStartsAtTheCurrentHour(t *testing.T) {
	now := time.Date(2026, 9, 29, 14, 20, 0, 0, zagreb)

	res := NewForecastResponse(testForecast(), now, now)

	if len(res.Hourly) != 72-14 {
		t.Fatalf("got %d hours, want %d (from 14:00 to the end)", len(res.Hourly), 72-14)
	}
	first := res.Hourly[0]
	if !first.Time.Equal(time.Date(2026, 9, 29, 14, 0, 0, 0, zagreb)) {
		t.Errorf("first hour = %v, want 14:00", first.Time)
	}
	if first.Temperature != 14 || first.ApparentTemperature != 12 || first.PrecipitationProbability != 114 ||
		first.Precipitation != 214 || first.WeatherCode != 14 || first.WindSpeed != 314 || first.WindDirection != 414 {
		t.Errorf("first hour values = %+v, want index 14 of each series", first)
	}
	if !res.Hourly[1].IsDay {
		t.Error("IsDay: a series value of 1 must read as day")
	}
}

func TestNewForecastResponseDailyStartsToday(t *testing.T) {
	now := time.Date(2026, 9, 30, 0, 30, 0, 0, zagreb)

	res := NewForecastResponse(testForecast(), now, now)

	if len(res.Daily) != 2 || res.Daily[0].Date != "2026-09-30" {
		t.Fatalf("daily = %+v, want 2 days from 2026-09-30", res.Daily)
	}
	d := res.Daily[0]
	if d.TemperatureMin != 8 || d.TemperatureMax != 18 || d.WeatherCode != 3 || d.PrecipitationProbabilityMax != 10 ||
		d.WindSpeedMax != 14 || d.WindGustsMax != 30 || d.WindDirectionDominant != 45 || d.UVIndexMax != 4 || d.DaylightSeconds != 42400 {
		t.Errorf("day = %+v", d)
	}
	if !d.Sunrise.Equal(time.Date(2026, 9, 30, 6, 47, 0, 0, zagreb)) || !d.Sunset.Equal(time.Date(2026, 9, 30, 18, 33, 0, 0, zagreb)) {
		t.Errorf("sun = %v – %v", d.Sunrise, d.Sunset)
	}
}

func TestNewForecastResponseToleratesShortSeries(t *testing.T) {
	f := testForecast()
	f.Hourly.Data[goopenmeteo.Temperature2M] = []float64{1, 2}
	delete(f.Daily.Times, goopenmeteo.Sunset)
	now := time.Date(2026, 9, 29, 0, 10, 0, 0, zagreb)

	res := NewForecastResponse(f, now, now)

	if res.Hourly[1].Temperature != 2 || res.Hourly[5].Temperature != 0 {
		t.Errorf("a short series should read as zeros past its end: %v, %v", res.Hourly[1].Temperature, res.Hourly[5].Temperature)
	}
	if !res.Daily[0].Sunset.IsZero() {
		t.Errorf("missing sunset = %v, want the zero time", res.Daily[0].Sunset)
	}
}

// Open-Meteo sends wall-clock times; goopenmeteo places them at the offset the
// response had, so after the switch to CET (2026-10-25 03:00) every hour would
// be an hour off unless it is placed in Europe/Zagreb itself.
func TestNewForecastResponseUsesTheZonesOwnOffsetAcrossDST(t *testing.T) {
	start := time.Date(2026, 10, 24, 0, 0, 0, 0, zagreb)
	f := &goopenmeteo.Forecast{
		Timezone: "Europe/Zagreb",
		Hourly: goopenmeteo.TimeseriesData{
			Time: hours(start, 48), // wall clock 24 Oct 00:00 … 25 Oct 23:00, all labelled +02:00
			Data: map[string][]float64{goopenmeteo.Temperature2M: seq(48, 0)},
		},
	}
	now := time.Date(2026, 10, 24, 0, 30, 0, 0, zagreb)

	res := NewForecastResponse(f, now, now)

	for _, h := range res.Hourly {
		wall := h.Time.Format("2006-01-02 15:04")
		if wall == "2026-10-24 12:00" && !h.Time.Equal(time.Date(2026, 10, 24, 10, 0, 0, 0, time.UTC)) {
			t.Errorf("24 Oct 12:00 = %v, want 10:00 UTC (CEST)", h.Time.UTC())
		}
		if wall == "2026-10-25 04:00" && !h.Time.Equal(time.Date(2026, 10, 25, 3, 0, 0, 0, time.UTC)) {
			t.Errorf("25 Oct 04:00 = %v, want 03:00 UTC (CET)", h.Time.UTC())
		}
	}
}
```

- [ ] **Step 2: Run them and watch them fail**

Run: `cd backend && go test ./app/models/`

Expected: FAIL with `undefined: NewForecastResponse`.

- [ ] **Step 3: Implement the forecast model**

`backend/app/models/forecast.go`:

```go
package models

import (
	"time"

	"github.com/h00s/goopenmeteo"
)

// The Open-Meteo variables NewForecastResponse reads; ForecastService requests
// exactly these.
var (
	ForecastCurrentVariables = goopenmeteo.WeatherVariables{
		goopenmeteo.Temperature2M,
		goopenmeteo.ApparentTemperature,
		goopenmeteo.RelativeHumidity2M,
		goopenmeteo.DewPoint2M,
		goopenmeteo.Precipitation,
		goopenmeteo.WeatherCode,
		goopenmeteo.CloudCover,
		goopenmeteo.PressureMsl,
		goopenmeteo.WindSpeed10M,
		goopenmeteo.WindDirection10M,
		goopenmeteo.WindGusts10M,
		goopenmeteo.Visibility,
		goopenmeteo.UVIndex, // declared with the air-quality variables; the forecast API has it too
		goopenmeteo.IsDay,
	}
	ForecastHourlyVariables = goopenmeteo.WeatherVariables{
		goopenmeteo.Temperature2M,
		goopenmeteo.ApparentTemperature,
		goopenmeteo.PrecipitationProbability,
		goopenmeteo.Precipitation,
		goopenmeteo.WeatherCode,
		goopenmeteo.WindSpeed10M,
		goopenmeteo.WindDirection10M,
		goopenmeteo.IsDay,
	}
	ForecastDailyVariables = goopenmeteo.WeatherVariables{
		goopenmeteo.Temperature2MMin,
		goopenmeteo.Temperature2MMax,
		goopenmeteo.WeatherCode,
		goopenmeteo.PrecipitationSum,
		goopenmeteo.PrecipitationProbabilityMax,
		goopenmeteo.WindSpeed10MMax,
		goopenmeteo.WindGusts10MMax,
		goopenmeteo.WindDirection10MDominant,
		goopenmeteo.UVIndexMax,
		goopenmeteo.Sunrise,
		goopenmeteo.Sunset,
		goopenmeteo.DaylightDuration,
	}
)

// ForecastDays is how many days are requested and returned, today first.
const ForecastDays = 7

type ForecastResponse struct {
	Current   CurrentForecastResponse  `json:"current"`
	Hourly    []HourlyForecastResponse `json:"hourly"`    // from the current hour to the end of the last day
	Daily     []DailyForecastResponse  `json:"daily"`     // ForecastDays days from today
	Timezone  string                   `json:"timezone"`  // the location's IANA zone, e.g. "Europe/Zagreb"
	Elevation float64                  `json:"elevation"` // metres
	FetchedAt time.Time                `json:"fetchedAt"`
}

type CurrentForecastResponse struct {
	Time                time.Time `json:"time"`
	Temperature         float64   `json:"temperature"`         // °C
	ApparentTemperature float64   `json:"apparentTemperature"` // °C
	Humidity            float64   `json:"humidity"`            // %
	DewPoint            float64   `json:"dewPoint"`            // °C
	Precipitation       float64   `json:"precipitation"`       // mm
	WeatherCode         int       `json:"weatherCode"`         // WMO
	CloudCover          float64   `json:"cloudCover"`          // %
	Pressure            float64   `json:"pressure"`            // hPa at sea level
	WindSpeed           float64   `json:"windSpeed"`           // km/h
	WindDirection       float64   `json:"windDirection"`       // degrees, where it blows from
	WindGusts           float64   `json:"windGusts"`           // km/h
	Visibility          float64   `json:"visibility"`          // m
	UVIndex             float64   `json:"uvIndex"`
	IsDay               bool      `json:"isDay"`
}

type HourlyForecastResponse struct {
	Time                     time.Time `json:"time"`
	Temperature              float64   `json:"temperature"`
	ApparentTemperature      float64   `json:"apparentTemperature"`
	PrecipitationProbability float64   `json:"precipitationProbability"`
	Precipitation            float64   `json:"precipitation"`
	WeatherCode              int       `json:"weatherCode"`
	WindSpeed                float64   `json:"windSpeed"`
	WindDirection            float64   `json:"windDirection"`
	IsDay                    bool      `json:"isDay"`
}

type DailyForecastResponse struct {
	Date                        string    `json:"date"` // YYYY-MM-DD in the location's zone
	TemperatureMin              float64   `json:"temperatureMin"`
	TemperatureMax              float64   `json:"temperatureMax"`
	WeatherCode                 int       `json:"weatherCode"`
	PrecipitationSum            float64   `json:"precipitationSum"`
	PrecipitationProbabilityMax float64   `json:"precipitationProbabilityMax"`
	WindSpeedMax                float64   `json:"windSpeedMax"`
	WindGustsMax                float64   `json:"windGustsMax"`
	WindDirectionDominant       float64   `json:"windDirectionDominant"`
	UVIndexMax                  float64   `json:"uvIndexMax"`
	Sunrise                     time.Time `json:"sunrise"`
	Sunset                      time.Time `json:"sunset"`
	DaylightSeconds             float64   `json:"daylightSeconds"`
}

// NewForecastResponse maps a forecast to the response: hourly from the hour
// containing now, daily from today. Times are placed in the location's own
// zone, so they keep the right offset across a DST change within the week:
// Open-Meteo sends wall-clock times, and goopenmeteo can only give them the
// response's current offset.
func NewForecastResponse(f *goopenmeteo.Forecast, now, fetchedAt time.Time) ForecastResponse {
	zone := zoneOf(f)
	cur := f.Current.Data
	res := ForecastResponse{
		Current: CurrentForecastResponse{
			Time:                zone(f.Current.Time),
			Temperature:         cur[goopenmeteo.Temperature2M],
			ApparentTemperature: cur[goopenmeteo.ApparentTemperature],
			Humidity:            cur[goopenmeteo.RelativeHumidity2M],
			DewPoint:            cur[goopenmeteo.DewPoint2M],
			Precipitation:       cur[goopenmeteo.Precipitation],
			WeatherCode:         int(cur[goopenmeteo.WeatherCode]),
			CloudCover:          cur[goopenmeteo.CloudCover],
			Pressure:            cur[goopenmeteo.PressureMsl],
			WindSpeed:           cur[goopenmeteo.WindSpeed10M],
			WindDirection:       cur[goopenmeteo.WindDirection10M],
			WindGusts:           cur[goopenmeteo.WindGusts10M],
			Visibility:          cur[goopenmeteo.Visibility],
			UVIndex:             cur[goopenmeteo.UVIndex],
			IsDay:               cur[goopenmeteo.IsDay] == 1,
		},
		Hourly:    make([]HourlyForecastResponse, 0, len(f.Hourly.Time)),
		Daily:     make([]DailyForecastResponse, 0, ForecastDays),
		Timezone:  f.Timezone,
		Elevation: f.Elevation,
		FetchedAt: fetchedAt,
	}

	hourly := f.Hourly.Data
	for i, t := range f.Hourly.Time {
		t = zone(t)
		if !t.Add(time.Hour).After(now) { // this hour is already over
			continue
		}
		res.Hourly = append(res.Hourly, HourlyForecastResponse{
			Time:                     t,
			Temperature:              at(hourly[goopenmeteo.Temperature2M], i),
			ApparentTemperature:      at(hourly[goopenmeteo.ApparentTemperature], i),
			PrecipitationProbability: at(hourly[goopenmeteo.PrecipitationProbability], i),
			Precipitation:            at(hourly[goopenmeteo.Precipitation], i),
			WeatherCode:              int(at(hourly[goopenmeteo.WeatherCode], i)),
			WindSpeed:                at(hourly[goopenmeteo.WindSpeed10M], i),
			WindDirection:            at(hourly[goopenmeteo.WindDirection10M], i),
			IsDay:                    at(hourly[goopenmeteo.IsDay], i) >= 1,
		})
	}

	daily, times := f.Daily.Data, f.Daily.Times
	for i, t := range f.Daily.Time {
		if len(res.Daily) == ForecastDays {
			break
		}
		t = zone(t)
		if !t.AddDate(0, 0, 1).After(now) { // this day is already over
			continue
		}
		res.Daily = append(res.Daily, DailyForecastResponse{
			Date:                        t.Format(time.DateOnly),
			TemperatureMin:              at(daily[goopenmeteo.Temperature2MMin], i),
			TemperatureMax:              at(daily[goopenmeteo.Temperature2MMax], i),
			WeatherCode:                 int(at(daily[goopenmeteo.WeatherCode], i)),
			PrecipitationSum:            at(daily[goopenmeteo.PrecipitationSum], i),
			PrecipitationProbabilityMax: at(daily[goopenmeteo.PrecipitationProbabilityMax], i),
			WindSpeedMax:                at(daily[goopenmeteo.WindSpeed10MMax], i),
			WindGustsMax:                at(daily[goopenmeteo.WindGusts10MMax], i),
			WindDirectionDominant:       at(daily[goopenmeteo.WindDirection10MDominant], i),
			UVIndexMax:                  at(daily[goopenmeteo.UVIndexMax], i),
			Sunrise:                     zone(at(times[goopenmeteo.Sunrise], i)),
			Sunset:                      zone(at(times[goopenmeteo.Sunset], i)),
			DaylightSeconds:             at(daily[goopenmeteo.DaylightDuration], i),
		})
	}

	return res
}

// zoneOf re-places wall-clock times in the forecast's IANA zone. With no zone
// (a GMT forecast) or an unknown one, times stay as decoded.
func zoneOf(f *goopenmeteo.Forecast) func(time.Time) time.Time {
	loc, err := time.LoadLocation(f.Timezone)
	if f.Timezone == "" || err != nil {
		return func(t time.Time) time.Time { return t }
	}
	return func(t time.Time) time.Time {
		if t.IsZero() {
			return t
		}
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
	}
}

// at returns series[i], or the zero value when the series is shorter than the
// time axis (a variable the response didn't include).
func at[T any](series []T, i int) T {
	if i < len(series) {
		return series[i]
	}
	var zero T
	return zero
}
```

In `backend/main.go`, add this import beside the others:

```go
	_ "time/tzdata" // forecasts are placed in their IANA zone wherever the binary runs
```

- [ ] **Step 4: Run the model tests**

Run: `cd backend && go test ./app/models/`

Expected: PASS.

- [ ] **Step 5: Write the failing endpoint tests**

Add to `backend/app/controllers/fixtures_test.go`. Use imports `encoding/json`, `fmt` and `time`, add the case to `canned`, and add the function:

```go
	case "/forecast":
		return forecastJSON(now), true
```

```go
// forecastJSON is an Open-Meteo forecast for Europe/Zagreb covering the seven
// days from today, so the hourly window always contains now.
func forecastJSON(now time.Time) []byte {
	zagreb, _ := time.LoadLocation("Europe/Zagreb")
	now = now.In(zagreb)
	_, offset := now.Zone()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, zagreb)

	var hourlyTimes []string
	var hourlyTemps []float64
	for i := range 7 * 24 {
		hourlyTimes = append(hourlyTimes, midnight.Add(time.Duration(i)*time.Hour).Format("2006-01-02T15:04"))
		hourlyTemps = append(hourlyTemps, float64(i%24))
	}
	var days, sunrises, sunsets []string
	var maxes []float64
	for i := range 7 {
		day := midnight.AddDate(0, 0, i)
		days = append(days, day.Format(time.DateOnly))
		sunrises = append(sunrises, day.Add(6*time.Hour+46*time.Minute).Format("2006-01-02T15:04"))
		sunsets = append(sunsets, day.Add(18*time.Hour+35*time.Minute).Format("2006-01-02T15:04"))
		maxes = append(maxes, 20+float64(i))
	}

	body, _ := json.Marshal(map[string]any{
		"latitude":              45.59,
		"longitude":             17.23,
		"elevation":             161,
		"utc_offset_seconds":    offset,
		"timezone":              "Europe/Zagreb",
		"timezone_abbreviation": fmt.Sprintf("GMT+%d", offset/3600),
		"current": map[string]any{
			"time": now.Truncate(15 * time.Minute).Format("2006-01-02T15:04"), "interval": 900,
			"temperature_2m": 17.4, "weather_code": 61, "is_day": 1, "wind_direction_10m": 225,
		},
		"hourly": map[string]any{"time": hourlyTimes, "temperature_2m": hourlyTemps},
		"daily":  map[string]any{"time": days, "temperature_2m_max": maxes, "sunrise": sunrises, "sunset": sunsets},
	})
	return body
}
```

`backend/app/controllers/forecast_test.go`:

```go
package controllers_test

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/models"
)

func TestForecastShowReturnsTheForecastForTheCell(t *testing.T) {
	app, u := newApp(t)

	res := raptor.DecodeJSON[models.ForecastResponse](t, app.TestGet("/api/v1/forecast?lat=45.5936&lon=17.2251"), http.StatusOK)

	c := res.Current
	if c.Temperature != 17.4 || c.WeatherCode != 61 || !c.IsDay || c.WindDirection != 225 {
		t.Errorf("current = %+v", c)
	}
	if len(res.Daily) != models.ForecastDays || len(res.Hourly) == 0 || res.Timezone != "Europe/Zagreb" || res.FetchedAt.IsZero() {
		t.Errorf("daily %d, hourly %d, timezone %q, fetchedAt %v", len(res.Daily), len(res.Hourly), res.Timezone, res.FetchedAt)
	}
	if first, now := res.Hourly[0].Time, time.Now(); first.After(now) || !first.Add(time.Hour).After(now) {
		t.Errorf("hourly starts at %v, want the hour containing %v", first, now)
	}

	q := u.query("/forecast")
	lat, _ := strconv.ParseFloat(q.Get("latitude"), 64)
	lon, _ := strconv.ParseFloat(q.Get("longitude"), 64)
	if lat != 45.59 || lon != 17.23 {
		t.Errorf("Open-Meteo got %v,%v; want the rounded cell 45.59,17.23", lat, lon)
	}
	if q.Get("timezone") != "auto" || q.Get("forecast_days") != "7" || q.Get("daily") == "" || q.Get("current") == "" {
		t.Errorf("Open-Meteo query = %v", q)
	}
}

func TestForecastShowSharesOneUpstreamRequestPerCell(t *testing.T) {
	app, u := newApp(t)

	for _, q := range []string{"lat=45.5936&lon=17.2251", "lat=45.5912&lon=17.2263"} {
		raptor.DecodeJSON[models.ForecastResponse](t, app.TestGet("/api/v1/forecast?"+q), http.StatusOK)
	}

	if n := u.hits("/forecast"); n != 1 {
		t.Errorf("Open-Meteo was asked %d times for one cell, want 1", n)
	}
}

func TestForecastShowRejectsBadCoordinates(t *testing.T) {
	app, u := newApp(t)

	for _, q := range []string{"", "lat=45.59", "lat=abc&lon=17", "lat=91&lon=17", "lat=45&lon=-180.5", "lat=NaN&lon=17"} {
		body := raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/forecast?"+q), http.StatusBadRequest)
		if body.Code != http.StatusBadRequest || body.Message == "" {
			t.Errorf("%q: body = %+v", q, body)
		}
	}
	if n := u.hits("/forecast"); n != 0 {
		t.Errorf("Open-Meteo was asked %d times for invalid coordinates", n)
	}
}

func TestForecastShowAnswers502WhenOpenMeteoFailsAndNothingIsCached(t *testing.T) {
	app, u := newApp(t)
	u.fail(http.StatusServiceUnavailable)

	body := raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/forecast?lat=45.59&lon=17.23"), http.StatusBadGateway)

	if body.Code != http.StatusBadGateway {
		t.Errorf("body = %+v", body)
	}
}
```

- [ ] **Step 6: Run them and watch them fail**

Run: `cd backend && go test ./app/controllers/ -run Forecast`

Expected: FAIL. The old controller answers raw goopenmeteo JSON, and the upstream stub is never called (it still targets the real API).

- [ ] **Step 7: Implement the service, the helper, the controller and the route**

`backend/app/services/forecast_service.go`:

```go
package services

import (
	"context"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/goopenmeteo"
	"github.com/h00s/weather/app/models"
)

const (
	// forecastTTL keeps "now" current without asking Open-Meteo twice in a few
	// minutes for the same square kilometre.
	forecastTTL = 10 * time.Minute
	// maxCachedLocations bounds every per-location cache (cells of ~1 km).
	maxCachedLocations = 5000
)

// ForecastService fetches forecasts from Open-Meteo, cached per ~1 km cell.
// APP_OPENMETEO_URL overrides the API root (tests).
type ForecastService struct {
	raptor.Service

	client *goopenmeteo.OpenMeteo
	cache  *keyedCache[string, *goopenmeteo.Forecast]
}

func (s *ForecastService) Setup() error {
	s.client = goopenmeteo.NewOpenMeteo()
	s.client.BaseURL = s.Config.AppString("openmeteo_url", goopenmeteo.BaseURL)
	s.cache = newKeyedCache[string, *goopenmeteo.Forecast](forecastTTL, maxCachedLocations)
	return nil
}

// Forecast returns the forecast for at's cell and when it was fetched. When
// Open-Meteo fails it returns the last forecast it has for the cell; it errors
// only when it has none.
func (s *ForecastService) Forecast(ctx context.Context, at models.Coordinates) (*goopenmeteo.Forecast, time.Time, error) {
	at = at.Rounded()
	return s.cache.Get(ctx, at.Key(), func(ctx context.Context) (*goopenmeteo.Forecast, error) {
		forecast, err := s.client.ForecastContext(ctx, goopenmeteo.ForecastOptions{
			Latitude:     at.Latitude,
			Longitude:    at.Longitude,
			Timezone:     "auto", // the location's own zone
			Current:      models.ForecastCurrentVariables,
			Hourly:       models.ForecastHourlyVariables,
			Daily:        models.ForecastDailyVariables,
			ForecastDays: models.ForecastDays,
		})
		if err != nil {
			s.Log.Warn("Open-Meteo forecast request failed", "cell", at.Key(), "error", err)
		}
		return forecast, err
	})
}
```

`backend/app/controllers/helpers.go`:

```go
package controllers

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/weather/app/models"
)

// coordinates reads the lat and lon query parameters, rounded to the ~1 km cell
// that every upstream request and cache key uses. A missing or invalid one is a 400.
func coordinates(ctx *raptor.Context) (models.Coordinates, error) {
	at, err := models.ParseCoordinates(ctx.QueryParam("lat"), ctx.QueryParam("lon"))
	if err != nil {
		return models.Coordinates{}, errs.NewErrorBadRequest(err.Error())
	}
	return at.Rounded(), nil
}
```

`backend/app/controllers/forecast_controller.go`:

```go
package controllers

import (
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/weather/app/models"
	"github.com/h00s/weather/app/services"
)

type ForecastController struct {
	raptor.Controller

	Forecast *services.ForecastService
}

// Show is the forecast at ?lat&lon: now, every hour to the end of the week, and
// seven days from today.
func (c *ForecastController) Show(ctx *raptor.Context) error {
	at, err := coordinates(ctx)
	if err != nil {
		return err
	}
	forecast, fetchedAt, err := c.Forecast.Forecast(ctx.Request().Context(), at)
	if err != nil {
		return errs.NewErrorBadGateway("Open-Meteo is unavailable")
	}
	return ctx.Data(models.NewForecastResponse(forecast, time.Now(), fetchedAt))
}
```

In `backend/config/routes.yaml`, replace the forecast entry:

```yaml
    /forecast:
      GET: Forecast.Show        # ?lat&lon: now, hourly to the end of the week, 7 days
```

- [ ] **Step 8: Run all backend tests**

Run: `cd backend && go vet ./... && go test ./...`

Expected: PASS.

- [ ] **Step 9: Commit**

```bash
git add -A backend
git commit -m "Serve a typed forecast for any location, cached per 1 km cell

Hourly to the end of the week and seven days, placed in the location's own
zone so hours keep the right offset across a DST change.

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 4: Places dataset and nearest place

**Files:**
- Create: `backend/cmd/genplaces/main.go`, `backend/app/services/data/hr-places.tsv` (generated), `backend/app/models/place.go`, `backend/app/models/place_test.go`, `backend/app/models/county.go`, `backend/app/models/county_test.go`, `backend/app/services/places_service.go`, `backend/app/services/places_service_test.go`, `backend/app/controllers/places_controller.go`, `backend/app/controllers/places_test.go`
- Modify: `backend/config/components/services.go`, `backend/config/components/controllers.go`, `backend/config/components/middlewares.go`, `backend/config/routes.yaml`

**Interfaces:**
- Consumes: `models.Coordinates` and `coordinates(ctx)`.
- Produces:
  - `models.Place{Name, Coordinates, County, Population}`;
  - `models.ParsePlaces(io.Reader) ([]models.Place, error)`;
  - `models.NearestPlace(places, at) (models.Place, bool)` and `models.DistanceKm(a, b) float64`;
  - `models.PlaceResponse{Name, County, Latitude, Longitude}`, json `name`, `county`, `latitude`, `longitude`, and `models.NewPlaceResponse(p)`;
  - `models.County{Code, Name, Areas}` and `models.Counties map[string]models.County`, keyed by GeoNames admin1 code `"01"`…`"21"`;
  - `services.PlacesService.Nearest(at) (models.Place, bool)`;
  - `GET /api/v1/places/nearest?lat&lon`, which answers 200 `PlaceResponse` or a JSON 404 abroad.

- [ ] **Step 1: Write the generator and produce the dataset**

`backend/cmd/genplaces/main.go`:

```go
// Command genplaces writes app/services/data/hr-places.tsv, the settlements
// PlacesService names locations after, from GeoNames' Croatia dump (CC BY 4.0).
//
//	cd backend && go run ./cmd/genplaces          # downloads HR.zip
//	cd backend && go run ./cmd/genplaces HR.zip   # or reads a local copy
//
// It keeps populated places (feature class P) except sections of a city (PPLX:
// they would name a point in Maksimir "Donji Bukovec" instead of Zagreb) and
// historical, abandoned or destroyed places.
package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"cmp"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

const (
	dumpURL = "https://download.geonames.org/export/dump/HR.zip"
	output  = "app/services/data/hr-places.tsv"
)

var skipped = map[string]bool{"PPLX": true, "PPLH": true, "PPLQ": true, "PPLW": true, "PPLCH": true}

func main() {
	archive, err := readArchive(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	rows, err := places(archive)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	out.WriteString("# Croatian settlements from GeoNames (https://www.geonames.org, CC BY 4.0), written by cmd/genplaces.\n")
	out.WriteString("# name\tlatitude\tlongitude\tcounty (GeoNames admin1)\tpopulation\n")
	for _, r := range rows {
		out.WriteString(strings.Join(r, "\t") + "\n")
	}
	if err := os.WriteFile(output, out.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: %d places\n", output, len(rows))
}

func readArchive(args []string) ([]byte, error) {
	if len(args) > 0 {
		return os.ReadFile(args[0])
	}
	resp, err := http.Get(dumpURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", dumpURL, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// places reads HR.txt, tab-separated with the columns of
// https://download.geonames.org/export/dump/readme.txt.
func places(archive []byte) ([][]string, error) {
	zr, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return nil, err
	}
	f, err := zr.Open("HR.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var rows [][]string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1<<20), 1<<20)
	for scanner.Scan() {
		c := strings.Split(scanner.Text(), "\t")
		if len(c) < 15 {
			continue
		}
		name, lat, lon, class, code, admin1, population := c[1], c[4], c[5], c[6], c[7], c[10], c[14]
		if class != "P" || skipped[code] || admin1 == "" {
			continue
		}
		if population == "" {
			population = "0"
		}
		rows = append(rows, []string{name, lat, lon, admin1, population})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	slices.SortFunc(rows, func(a, b []string) int {
		return cmp.Or(strings.Compare(a[0], b[0]), strings.Compare(a[1], b[1]), strings.Compare(a[2], b[2]))
	})
	return rows, nil
}
```

Run: `cd backend && mkdir -p app/services/data && go run ./cmd/genplaces`

Expected: `app/services/data/hr-places.tsv: 11365 places`. The count may drift slightly as GeoNames updates; anything above 11,000 is right. Check the file with `head -4 app/services/data/hr-places.tsv` and `du -h app/services/data/hr-places.tsv`, which should be about 400 KB.

- [ ] **Step 2: Write the failing model tests**

`backend/app/models/place_test.go`:

```go
package models

import (
	"strings"
	"testing"
)

func TestParsePlaces(t *testing.T) {
	places, err := ParsePlaces(strings.NewReader("# comment\nDaruvar\t45.59056\t17.225\t01\t7440\n\nZagreb\t45.81444\t15.97798\t21\t663592\n"))
	if err != nil {
		t.Fatal(err)
	}
	want := Place{Name: "Daruvar", Coordinates: Coordinates{45.59056, 17.225}, County: "01", Population: 7440}
	if len(places) != 2 || places[0] != want || places[1].Name != "Zagreb" {
		t.Errorf("places = %+v", places)
	}

	for _, bad := range []string{"Daruvar\t45.59\n", "Daruvar\tx\t17.2\t01\t1\n", "Daruvar\t45.5\t17.2\t01\tmany\n"} {
		if _, err := ParsePlaces(strings.NewReader(bad)); err == nil {
			t.Errorf("ParsePlaces(%q) accepted a bad line", bad)
		}
	}
}

func TestDistanceKm(t *testing.T) {
	zagreb, split := Coordinates{45.81444, 15.97798}, Coordinates{43.50891, 16.43915}
	if d := DistanceKm(zagreb, split); d < 255 || d > 262 {
		t.Errorf("Zagreb–Split = %.1f km, want about 259", d)
	}
	if d := DistanceKm(zagreb, zagreb); d != 0 {
		t.Errorf("distance to itself = %v", d)
	}
}

func TestNearestPlaceLetsLargerPlacesReachFurther(t *testing.T) {
	// One degree of longitude at 45°N is 78.7 km: the village sits 3.9 km east of the city.
	places := []Place{
		{Name: "City", Coordinates: Coordinates{45.0, 16.0}, County: "21", Population: 600_000}, // radius ≈ 14.7 km
		{Name: "Village", Coordinates: Coordinates{45.0, 16.05}, County: "20"},                  // radius 1 km
	}
	cases := []struct {
		at   Coordinates
		want string // "" for no place
	}{
		{Coordinates{45.0, 16.04}, "City"},     // 0.8 km from the village, but well inside the city's reach
		{Coordinates{45.0, 16.049}, "Village"}, // right at the village
		{Coordinates{45.0, 16.25}, "Village"},  // the city is beyond 15 km
		{Coordinates{45.0, 16.3}, ""},          // nothing within 15 km
		{Coordinates{46.5, 16.0}, ""},
	}
	for _, c := range cases {
		got, ok := NearestPlace(places, c.at)
		if c.want == "" && ok {
			t.Errorf("NearestPlace(%v) = %s, want none", c.at, got.Name)
		}
		if c.want != "" && (!ok || got.Name != c.want) {
			t.Errorf("NearestPlace(%v) = %s (%v), want %s", c.at, got.Name, ok, c.want)
		}
	}
}

func TestNewPlaceResponseNamesTheCounty(t *testing.T) {
	got := NewPlaceResponse(Place{Name: "Daruvar", Coordinates: Coordinates{45.59056, 17.225}, County: "01"})
	want := PlaceResponse{Name: "Daruvar", County: "Bjelovarsko-bilogorska županija", Latitude: 45.59056, Longitude: 17.225}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
```

`backend/app/models/county_test.go`:

```go
package models

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestCountiesCoverEveryGeoNamesCode(t *testing.T) {
	if len(Counties) != 21 {
		t.Fatalf("%d counties, want 21", len(Counties))
	}
	for i := 1; i <= 21; i++ {
		code := fmt.Sprintf("%02d", i)
		c, ok := Counties[code]
		if !ok || c.Code != code || c.Name == "" || len(c.Areas) < 3 {
			t.Errorf("county %s = %+v", code, c)
		}
	}
}

// The coastal names and NUTS codes are pinned by the recorded Meteoalarm feed
// (testdata/meteoalarm-croatia.json, Task 6).
func TestCountiesUseMeteoalarmsAreaNames(t *testing.T) {
	for code, area := range map[string]string{"03": "Dubrovačko-neretvanska", "08": "Ličko-senjska", "12": "Primorsko-goranska", "13": "Šibensko-kninska", "15": "Splitsko-dalmatinska", "19": "Zadarska"} {
		if !slices.Contains(Counties[code].Areas, area) {
			t.Errorf("county %s areas %v lack %q", code, Counties[code].Areas, area)
		}
	}
	for _, c := range Counties {
		if !slices.ContainsFunc(c.Areas, func(a string) bool { return strings.HasPrefix(a, "HR0") }) {
			t.Errorf("%s has no NUTS3 code", c.Name)
		}
	}
}
```

- [ ] **Step 3: Run them and watch them fail**

Run: `cd backend && go test ./app/models/`

Expected: FAIL with `undefined: ParsePlaces` and `undefined: Counties`.

- [ ] **Step 4: Implement places and counties**

`backend/app/models/place.go`:

```go
package models

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// Place is a Croatian settlement from the GeoNames dataset.
type Place struct {
	Name string
	Coordinates
	County     string // GeoNames admin1 code, a key of Counties
	Population int
}

// PlaceResponse is a named place to show weather for.
type PlaceResponse struct {
	Name      string  `json:"name"`
	County    string  `json:"county"` // e.g. "Bjelovarsko-bilogorska županija"; "" when unknown
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewPlaceResponse(p Place) PlaceResponse {
	return PlaceResponse{Name: p.Name, County: Counties[p.County].Name, Latitude: p.Latitude, Longitude: p.Longitude}
}

// ParsePlaces reads the tab-separated dataset cmd/genplaces writes: name,
// latitude, longitude, county code and population per line; # starts a comment.
func ParsePlaces(r io.Reader) ([]Place, error) {
	var places []Place
	scanner := bufio.NewScanner(r)
	for n := 1; scanner.Scan(); n++ {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) != 5 {
			return nil, fmt.Errorf("line %d: %d fields, want 5", n, len(f))
		}
		lat, errLat := strconv.ParseFloat(f[1], 64)
		lon, errLon := strconv.ParseFloat(f[2], 64)
		population, errPop := strconv.Atoi(f[4])
		if err := cmp.Or(errLat, errLon, errPop); err != nil {
			return nil, fmt.Errorf("line %d: %w", n, err)
		}
		places = append(places, Place{Name: f[0], Coordinates: Coordinates{lat, lon}, County: f[3], Population: population})
	}
	return places, scanner.Err()
}

// placeReachKm is how far a point may be from every settlement and still be
// named: farther, it is abroad or out at sea.
const placeReachKm = 15

// NearestPlace names the settlement a point is in. Larger places reach further:
// distance counts in units of a place's radius, max(1 km, 0.6 km·√(population/1000)),
// so a point in Maksimir is Zagreb (3 km from a city of 660,000) and not the
// village 1.5 km away, while a field near Daruvar is the hamlet beside it.
// ok is false when no settlement lies within placeReachKm.
func NearestPlace(places []Place, at Coordinates) (place Place, ok bool) {
	best := math.Inf(1)
	for _, p := range places {
		// A cheap box first: 0.2° of latitude and 0.3° of longitude both exceed
		// 15 km across Croatia, so this only skips places out of reach.
		if math.Abs(p.Latitude-at.Latitude) > 0.2 || math.Abs(p.Longitude-at.Longitude) > 0.3 {
			continue
		}
		d := DistanceKm(at, p.Coordinates)
		if d > placeReachKm {
			continue
		}
		if score := d / p.radiusKm(); score < best {
			place, best, ok = p, score, true
		}
	}
	return place, ok
}

func (p Place) radiusKm() float64 {
	return max(1, 0.6*math.Sqrt(float64(p.Population)/1000))
}

// DistanceKm is the great-circle distance between a and b.
func DistanceKm(a, b Coordinates) float64 {
	const earthRadiusKm = 6371
	rad := math.Pi / 180
	dLat := (b.Latitude - a.Latitude) * rad
	dLon := (b.Longitude - a.Longitude) * rad
	h := math.Pow(math.Sin(dLat/2), 2) + math.Cos(a.Latitude*rad)*math.Cos(b.Latitude*rad)*math.Pow(math.Sin(dLon/2), 2)
	return 2 * earthRadiusKm * math.Asin(math.Sqrt(h))
}
```

`backend/app/models/county.go`:

```go
package models

// County is one of Croatia's 21 counties (županije). Areas name it as
// Meteoalarm's warning areas may: the areaDesc DHMZ uses ("Zadarska"), the full
// name, and its NUTS 2021 code. ActiveWarnings matches any of them.
type County struct {
	Code  string // GeoNames admin1 code
	Name  string // e.g. "Bjelovarsko-bilogorska županija"
	Areas []string
}

func county(code, short, nuts string) County {
	name := short + " županija"
	return County{Code: code, Name: name, Areas: []string{short, name, nuts}}
}

// Counties by GeoNames admin1 code. Coastal areaDesc and NUTS codes are checked
// against a recorded feed; inland ones follow the same pattern. If a live
// warning for an inland county ever fails to show, compare its areaDesc here.
var Counties = map[string]County{
	"01": county("01", "Bjelovarsko-bilogorska", "HR021"),
	"02": county("02", "Brodsko-posavska", "HR024"),
	"03": county("03", "Dubrovačko-neretvanska", "HR037"),
	"04": county("04", "Istarska", "HR036"),
	"05": county("05", "Karlovačka", "HR027"),
	"06": county("06", "Koprivničko-križevačka", "HR063"),
	"07": county("07", "Krapinsko-zagorska", "HR064"),
	"08": county("08", "Ličko-senjska", "HR032"),
	"09": county("09", "Međimurska", "HR061"),
	"10": county("10", "Osječko-baranjska", "HR025"),
	"11": county("11", "Požeško-slavonska", "HR023"),
	"12": county("12", "Primorsko-goranska", "HR031"),
	"13": county("13", "Šibensko-kninska", "HR034"),
	"14": county("14", "Sisačko-moslavačka", "HR028"),
	"15": county("15", "Splitsko-dalmatinska", "HR035"),
	"16": county("16", "Varaždinska", "HR062"),
	"17": county("17", "Virovitičko-podravska", "HR022"),
	"18": county("18", "Vukovarsko-srijemska", "HR026"),
	"19": county("19", "Zadarska", "HR033"),
	"20": county("20", "Zagrebačka", "HR065"),
	"21": {Code: "21", Name: "Grad Zagreb", Areas: []string{"Grad Zagreb", "Zagreb", "HR050"}},
}
```

- [ ] **Step 5: Run the model tests**

Run: `cd backend && go test ./app/models/`

Expected: PASS.

- [ ] **Step 6: Write the failing dataset test, then the service**

`backend/app/services/places_service_test.go`:

```go
package services

import (
	"strings"
	"testing"

	"github.com/h00s/weather/app/models"
)

// The embedded dataset names real points the way a reader would.
func TestEmbeddedPlacesNameRealPoints(t *testing.T) {
	places, err := models.ParsePlaces(strings.NewReader(placesTSV))
	if err != nil {
		t.Fatal(err)
	}
	if len(places) < 10_000 {
		t.Fatalf("%d places, want the whole country", len(places))
	}

	cases := []struct {
		name         string
		at           models.Coordinates
		want, county string
	}{
		{"Zagreb, main square", models.Coordinates{45.813, 15.977}, "Zagreb", "21"},
		{"Zagreb, Maksimir", models.Coordinates{45.82, 16.02}, "Zagreb", "21"},
		{"Sesvete", models.Coordinates{45.83, 16.11}, "Sesvete", "21"},
		{"Velika Gorica", models.Coordinates{45.71, 16.07}, "Velika Gorica", "20"},
		{"Daruvar", models.Coordinates{45.59, 17.22}, "Daruvar", "01"},
		{"a field near Daruvar", models.Coordinates{45.62, 17.30}, "Gornja Vrijeska", "01"},
		{"Split, Riva", models.Coordinates{43.508, 16.44}, "Split", "15"},
		{"Hvar", models.Coordinates{43.172, 16.442}, "Hvar", "15"},
	}
	for _, c := range cases {
		got, ok := models.NearestPlace(places, c.at)
		if !ok || got.Name != c.want || got.County != c.county {
			t.Errorf("%s: got %q in %s (%v), want %q in %s", c.name, got.Name, got.County, ok, c.want, c.county)
		}
	}

	for name, at := range map[string]models.Coordinates{
		"Sarajevo":        {43.86, 18.41},
		"the open Adriatic": {43.3, 15.5},
		"Berlin":          {52.52, 13.4},
	} {
		if got, ok := models.NearestPlace(places, at); ok {
			t.Errorf("%s was named %q, want no Croatian place", name, got.Name)
		}
	}
}
```

Run: `cd backend && go test ./app/services/`

Expected: FAIL with `undefined: placesTSV`.

`backend/app/services/places_service.go`:

```go
package services

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/models"
)

// placesTSV is written by cmd/genplaces from GeoNames (CC BY 4.0).
//
//go:embed data/hr-places.tsv
var placesTSV string

// PlacesService names locations after Croatian settlements.
type PlacesService struct {
	raptor.Service

	places []models.Place
}

func (s *PlacesService) Setup() error {
	places, err := models.ParsePlaces(strings.NewReader(placesTSV))
	if err != nil {
		return fmt.Errorf("places dataset: %w", err)
	}
	s.places = places
	return nil
}

// Nearest is the settlement at at; ok is false abroad or out at sea.
func (s *PlacesService) Nearest(at models.Coordinates) (models.Place, bool) {
	return models.NearestPlace(s.places, at)
}
```

Run: `cd backend && go test ./app/services/`

Expected: PASS.

- [ ] **Step 7: Write the failing endpoint tests**

`backend/app/controllers/places_test.go`:

```go
package controllers_test

import (
	"net/http"
	"testing"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/models"
)

func TestPlacesNearestNamesTheSettlement(t *testing.T) {
	app, _ := newApp(t)

	got := raptor.DecodeJSON[models.PlaceResponse](t, app.TestGet("/api/v1/places/nearest?lat=45.59&lon=17.22"), http.StatusOK)

	if got.Name != "Daruvar" || got.County != "Bjelovarsko-bilogorska županija" {
		t.Errorf("got %+v", got)
	}
}

func TestPlacesNearestIs404AbroadAnd400ForBadCoordinates(t *testing.T) {
	app, _ := newApp(t)

	raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/places/nearest?lat=43.86&lon=18.41"), http.StatusNotFound)
	raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/places/nearest?lat=x&lon=18"), http.StatusBadRequest)
}
```

Run: `cd backend && go test ./app/controllers/ -run Places`

Expected: FAIL, because the route is unknown (JSON 404 on both).

- [ ] **Step 8: Add the controller, the registration and the route**

`backend/app/controllers/places_controller.go`:

```go
package controllers

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/weather/app/models"
	"github.com/h00s/weather/app/services"
)

type PlacesController struct {
	raptor.Controller

	Places *services.PlacesService
}

// Nearest names the Croatian settlement at ?lat&lon; abroad or at sea it is a 404.
func (c *PlacesController) Nearest(ctx *raptor.Context) error {
	at, err := coordinates(ctx)
	if err != nil {
		return err
	}
	place, ok := c.Places.Nearest(at)
	if !ok {
		return errs.NewErrorNotFound("No Croatian place near these coordinates")
	}
	return ctx.Data(models.NewPlaceResponse(place))
}
```

- In `backend/config/components/services.go`, register `&services.PlacesService{}` after `&services.ForecastService{}`.
- In `controllers.go`, register `&controllers.PlacesController{}` after the forecast controller.
- In `middlewares.go`, change the limiter scope to `"Forecast", "Places"`.
- In `routes.yaml`, under `/api/v1`, add:

```yaml
    /places:
      /nearest:
        GET: Places.Nearest     # ?lat&lon: the Croatian settlement at a point (404 abroad)
```

- [ ] **Step 9: Run all backend tests**

Run: `cd backend && go vet ./... && go test ./...`

Expected: PASS.

- [ ] **Step 10: Commit**

```bash
git add -A backend
git commit -m "Name locations after the nearest Croatian settlement

An embedded GeoNames dataset (cmd/genplaces) where larger places reach
further, the 21 counties with their Meteoalarm area names, and
GET /api/v1/places/nearest.

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 5: Place search

**Files:**
- Create: `backend/app/models/geocoding.go`, `backend/app/models/geocoding_test.go`
- Modify: `backend/app/services/places_service.go`, `backend/app/controllers/places_controller.go`, `backend/app/controllers/fixtures_test.go`, `backend/app/controllers/places_test.go`, `backend/config/routes.yaml`

**Interfaces:**
- Consumes: `models.PlaceResponse` and `newKeyedCache`.
- Produces:
  - `models.GeocodingResponse`, `models.GeocodingResult` and `models.SearchResults(g, limit) []models.PlaceResponse`;
  - `services.PlacesService.Search(ctx, query) ([]models.PlaceResponse, error)`;
  - `GET /api/v1/places?q=`, which answers 200 with `[]PlaceResponse`, a 400 for a bad `q`, or a 502.

- [ ] **Step 1: Write the failing model test**

`backend/app/models/geocoding_test.go`:

```go
package models

import (
	"encoding/json"
	"slices"
	"testing"
)

// Recorded from geocoding-api.open-meteo.com (?name=Daruvar&language=hr&countryCode=HR),
// plus a Bosnian and an abandoned place the filter must drop.
const daruvarSearch = `{"results":[
 {"id":3202184,"name":"Daruvar","latitude":45.59056,"longitude":17.225,"feature_code":"PPLA2","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija","population":7440},
 {"id":11500092,"name":"Daruvar","latitude":45.58507,"longitude":17.2114,"feature_code":"AIRF","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija"},
 {"id":12509853,"name":"Daruvarski Vinogradi","latitude":45.60251,"longitude":17.25084,"feature_code":"PPL","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija","population":182},
 {"id":1,"name":"Daruvarac","latitude":44.5,"longitude":18.1,"feature_code":"PPL","country_code":"BA","admin1":"Federacija Bosne i Hercegovine"},
 {"id":2,"name":"Daruvarska Stara","latitude":45.6,"longitude":17.2,"feature_code":"PPLQ","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija"}
],"generationtime_ms":0.31}`

func TestSearchResultsKeepsCroatianSettlementsInOrder(t *testing.T) {
	var g GeocodingResponse
	if err := json.Unmarshal([]byte(daruvarSearch), &g); err != nil {
		t.Fatal(err)
	}

	want := []PlaceResponse{
		{Name: "Daruvar", County: "Bjelovarsko-bilogorska županija", Latitude: 45.59056, Longitude: 17.225},
		{Name: "Daruvarski Vinogradi", County: "Bjelovarsko-bilogorska županija", Latitude: 45.60251, Longitude: 17.25084},
	}
	if got := SearchResults(g, 10); !slices.Equal(got, want) {
		t.Errorf("got %+v\nwant %+v", got, want)
	}
	if got := SearchResults(g, 1); len(got) != 1 || got[0].Name != "Daruvar" {
		t.Errorf("limit 1: %+v", got)
	}
	if got := SearchResults(GeocodingResponse{}, 10); got == nil || len(got) != 0 {
		t.Errorf("no results = %#v, want an empty, non-nil slice", got)
	}
}
```

Run: `cd backend && go test ./app/models/ -run SearchResults`

Expected: FAIL with `undefined: GeocodingResponse`.

- [ ] **Step 2: Implement the geocoding model**

`backend/app/models/geocoding.go`:

```go
package models

import "strings"

// GeocodingResponse is Open-Meteo's geocoding answer
// (https://open-meteo.com/en/docs/geocoding-api), reduced to what search uses.
// results is absent when nothing matches.
type GeocodingResponse struct {
	Results []GeocodingResult `json:"results"`
}

type GeocodingResult struct {
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	FeatureCode string  `json:"feature_code"` // GeoNames feature code: PPL… for settlements
	CountryCode string  `json:"country_code"`
	Admin1      string  `json:"admin1"` // the county, in Croatian with language=hr
}

// gone are settlements that no longer exist: historical, abandoned, destroyed.
var gone = map[string]bool{"PPLH": true, "PPLQ": true, "PPLW": true, "PPLCH": true}

// SearchResults keeps the Croatian settlements, at most limit, in Open-Meteo's
// order: the best match first, larger places before smaller ones of one name.
func SearchResults(g GeocodingResponse, limit int) []PlaceResponse {
	places := make([]PlaceResponse, 0, min(limit, len(g.Results)))
	for _, r := range g.Results {
		if len(places) == limit {
			break
		}
		if r.CountryCode != "HR" || !strings.HasPrefix(r.FeatureCode, "PPL") || gone[r.FeatureCode] {
			continue
		}
		places = append(places, PlaceResponse{Name: r.Name, County: r.Admin1, Latitude: r.Latitude, Longitude: r.Longitude})
	}
	return places
}
```

Run: `cd backend && go test ./app/models/`

Expected: PASS.

- [ ] **Step 3: Write the failing endpoint tests**

In `backend/app/controllers/fixtures_test.go`, add the case and the constant:

```go
	case "/search":
		return []byte(searchJSON), true
```

```go
// searchJSON is Open-Meteo's geocoding answer for "Daru": two settlements and an airfield.
const searchJSON = `{"results":[
 {"id":3202184,"name":"Daruvar","latitude":45.59056,"longitude":17.225,"feature_code":"PPLA2","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija"},
 {"id":11500092,"name":"Daruvar","latitude":45.58507,"longitude":17.2114,"feature_code":"AIRF","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija"},
 {"id":12509853,"name":"Daruvarski Vinogradi","latitude":45.60251,"longitude":17.25084,"feature_code":"PPL","country_code":"HR","admin1":"Bjelovarsko-bilogorska županija"}
],"generationtime_ms":0.3}`
```

Append to `backend/app/controllers/places_test.go` (add `strings` to its imports):

```go
func TestPlacesIndexSearchesCroatianSettlements(t *testing.T) {
	app, u := newApp(t)

	got := raptor.DecodeJSON[[]models.PlaceResponse](t, app.TestGet("/api/v1/places?q=%20Daru%20"), http.StatusOK)

	if len(got) != 2 || got[0].Name != "Daruvar" || got[1].Name != "Daruvarski Vinogradi" {
		t.Errorf("got %+v", got)
	}
	q := u.query("/search")
	if q.Get("name") != "Daru" || q.Get("language") != "hr" || q.Get("countryCode") != "HR" {
		t.Errorf("geocoding query = %v", q)
	}
}

func TestPlacesIndexCachesAQueryWhateverItsCase(t *testing.T) {
	app, u := newApp(t)

	for _, q := range []string{"Daru", "daru", "DARU"} {
		raptor.DecodeJSON[[]models.PlaceResponse](t, app.TestGet("/api/v1/places?q="+q), http.StatusOK)
	}

	if n := u.hits("/search"); n != 1 {
		t.Errorf("geocoding asked %d times, want 1", n)
	}
}

func TestPlacesIndexRejectsTooShortOrTooLongQueries(t *testing.T) {
	app, u := newApp(t)

	for _, q := range []string{"", "q=", "q=a", "q=%20%20a%20", "q=" + strings.Repeat("a", 61)} {
		raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/places?"+q), http.StatusBadRequest)
	}
	if n := u.hits("/search"); n != 0 {
		t.Errorf("geocoding asked %d times for invalid queries", n)
	}
}

func TestPlacesIndexAnswers502WhenGeocodingFails(t *testing.T) {
	app, u := newApp(t)
	u.fail(http.StatusInternalServerError)

	raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/places?q=Daru"), http.StatusBadGateway)
}
```

Run: `cd backend && go test ./app/controllers/ -run PlacesIndex`

Expected: FAIL, because `/api/v1/places` has no GET route.

- [ ] **Step 4: Implement search in the service, the controller and the route**

In `backend/app/services/places_service.go`:
- the imports become `context`, `embed` (blank), `encoding/json`, `fmt`, `net/http`, `net/url`, `strings`, `time`, raptor and models;
- the struct and `Setup` become:

```go
const (
	geocodingURL      = "https://geocoding-api.open-meteo.com/v1"
	searchTTL         = 24 * time.Hour // place names don't move
	maxCachedSearches = 2000
	searchLimit       = 10
)

// PlacesService names locations after Croatian settlements and finds them by
// name through Open-Meteo's geocoding API. APP_GEOCODING_URL overrides its root
// (tests).
type PlacesService struct {
	raptor.Service

	places   []models.Place
	client   *http.Client
	baseURL  string
	searches *keyedCache[string, []models.PlaceResponse]
}

func (s *PlacesService) Setup() error {
	places, err := models.ParsePlaces(strings.NewReader(placesTSV))
	if err != nil {
		return fmt.Errorf("places dataset: %w", err)
	}
	s.places = places
	s.client = &http.Client{Timeout: 10 * time.Second}
	s.baseURL = strings.TrimSuffix(s.Config.AppString("geocoding_url", geocodingURL), "/")
	s.searches = newKeyedCache[string, []models.PlaceResponse](searchTTL, maxCachedSearches)
	return nil
}
```

Then append:

```go
// Search finds Croatian settlements by name, the best match first (Open-Meteo
// matches two letters exactly and three or more fuzzily). When the geocoding
// API fails it returns the last answer to the same query, or the error when
// there is none.
func (s *PlacesService) Search(ctx context.Context, query string) ([]models.PlaceResponse, error) {
	places, _, err := s.searches.Get(ctx, strings.ToLower(query), func(ctx context.Context) ([]models.PlaceResponse, error) {
		places, err := s.geocode(ctx, query)
		if err != nil {
			s.Log.Warn("Open-Meteo geocoding request failed", "error", err)
		}
		return places, err
	})
	return places, err
}

func (s *PlacesService) geocode(ctx context.Context, query string) ([]models.PlaceResponse, error) {
	params := url.Values{
		"name":        {query},
		"count":       {"20"}, // the filter drops airfields, stations and other features
		"language":    {"hr"},
		"countryCode": {"HR"},
		"format":      {"json"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/search?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding: unexpected status %d", resp.StatusCode)
	}
	var g models.GeocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, fmt.Errorf("geocoding: %w", err)
	}
	return models.SearchResults(g, searchLimit), nil
}
```

Add to `backend/app/controllers/places_controller.go`, importing `strings` and `unicode/utf8`:

```go
// Index finds settlements by name: ?q= of 2 to 60 characters, best match first.
func (c *PlacesController) Index(ctx *raptor.Context) error {
	q := strings.TrimSpace(ctx.QueryParam("q"))
	if n := utf8.RuneCountInString(q); n < 2 || n > 60 {
		return errs.NewErrorBadRequest("q: want 2 to 60 characters")
	}
	places, err := c.Places.Search(ctx.Request().Context(), q)
	if err != nil {
		return errs.NewErrorBadGateway("Open-Meteo geocoding is unavailable")
	}
	return ctx.Data(places)
}
```

`backend/config/routes.yaml`. The `/places` block becomes:

```yaml
    /places:
      GET: Places.Index         # ?q=: Croatian settlements by name, best match first
      /nearest:
        GET: Places.Nearest     # ?lat&lon: the Croatian settlement at a point (404 abroad)
```

- [ ] **Step 5: Run all backend tests**

Run: `cd backend && go vet ./... && go test ./...`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add -A backend
git commit -m "Search Croatian settlements by name through Open-Meteo geocoding

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 6: DHMZ warnings

**Files:**
- Create, ported from dashboard: `backend/app/models/warning.go`, `backend/app/models/warning_test.go`, `backend/app/models/testdata/meteoalarm-croatia.json`, `backend/app/services/warnings_service.go`
- Create: `backend/app/controllers/warnings_controller.go`, `backend/app/controllers/warnings_test.go`
- Modify: `backend/app/controllers/fixtures_test.go`, `backend/config/components/{services,controllers,middlewares}.go`, `backend/config/routes.yaml`

**Interfaces:**
- Consumes: `models.Counties[code].Areas`, `services.PlacesService.Nearest` and the test zone `zagreb` (Task 3).
- Produces:
  - `models.ActiveWarnings(alerts, areas []string, now) []models.WarningResponse` and `models.WarningLookahead = 48h`;
  - `services.WarningsService.Alerts(ctx)`;
  - `GET /api/v1/warnings?lat&lon`, which answers 200 with `[]WarningResponse` (json `level`, `type`, `event`, `description`, `onset`, `expires`) or a 502.

- [ ] **Step 1: Port the model and its tests, then change them to the county's list of areas**

```bash
cp ~/dev/go/husak-dashboard/backend/app/models/warning.go ~/dev/go/husak-dashboard/backend/app/models/warning_test.go backend/app/models/
mkdir -p backend/app/models/testdata
cp ~/dev/go/husak-dashboard/backend/app/models/testdata/meteoalarm-croatia.json backend/app/models/testdata/
sed -i -E 's/ActiveWarnings\(alerts, ([^,]+), /ActiveWarnings(alerts, []string{\1}, /' backend/app/models/warning_test.go
sed -i 's/warnNow.Add(30\*time.Hour), warnNow.Add(40\*time.Hour)/warnNow.Add(50*time.Hour), warnNow.Add(60*time.Hour)/' backend/app/models/warning_test.go
```

The first `sed` wraps each area argument in a slice. The second moves the "far" alert beyond the new 48-hour window.

Append to `backend/app/models/warning_test.go`:

```go
// The app shows today and tomorrow: a warning starting 30 hours ahead is listed.
func TestActiveWarningsKeepsTomorrowsWarnings(t *testing.T) {
	area := "Bjelovarsko-bilogorska"
	alerts := []MeteoalarmAlert{alert("tomorrow", "Alert", area, "2; yellow; Moderate", "1; Wind", warnNow.Add(30*time.Hour), warnNow.Add(36*time.Hour), "")}

	if got := ActiveWarnings(alerts, []string{area}, warnNow); len(got) != 1 {
		t.Errorf("got %d warnings, want tomorrow's", len(got))
	}
}

func TestActiveWarningsMatchesAnyOfACountysAreas(t *testing.T) {
	alerts := []MeteoalarmAlert{alert("a", "Alert", "Bjelovarsko-bilogorska", "2; yellow; Moderate", "1; Wind", warnNow, warnNow.Add(time.Hour), "")}

	if got := ActiveWarnings(alerts, Counties["01"].Areas, warnNow); len(got) != 1 {
		t.Errorf("Bjelovarsko-bilogorska matched %d warnings, want 1", len(got))
	}
	if got := ActiveWarnings(alerts, Counties["19"].Areas, warnNow); len(got) != 0 {
		t.Errorf("Zadarska matched %d warnings, want 0", len(got))
	}
}
```

Run: `cd backend && go test ./app/models/`

Expected: FAIL. `ActiveWarnings` still takes a string (the build fails).

- [ ] **Step 2: Change the model**

In `backend/app/models/warning.go`:
- replace the lookahead constant:

```go
// WarningLookahead is how far ahead an upcoming warning is already shown: today and tomorrow.
const WarningLookahead = 48 * time.Hour
```

- change the signature and its doc comment:

```go
// ActiveWarnings returns the yellow, orange and red warnings for any of areas
// (a county's areaDesc, full name or NUTS3 code; see County.Areas) that are in
// force now or start within WarningLookahead, most severe first. Cancelled
// alerts and alerts a later update replaced are dropped.
func ActiveWarnings(alerts []MeteoalarmAlert, areas []string, now time.Time) []WarningResponse {
```

- inside it, change `!info.covers(area)` to `!info.covers(areas)`;
- replace `covers` with:

```go
func (info *MeteoalarmInfo) covers(areas []string) bool {
	for _, a := range info.Area {
		for _, area := range areas {
			if strings.EqualFold(a.AreaDesc, area) {
				return true
			}
			for _, g := range a.Geocode {
				if strings.EqualFold(g.Value, area) {
					return true
				}
			}
		}
	}
	return false
}
```

Run: `cd backend && go test ./app/models/`

Expected: PASS. The recorded-feed test still finds exactly one Zadarska warning with the 48-hour window; this was checked against the fixture.

- [ ] **Step 3: Write the failing endpoint tests**

In `backend/app/controllers/fixtures_test.go`, add the case and the function:

```go
	case "/api/v1/warnings/feeds-croatia":
		return meteoalarmJSON(now), true
```

```go
// meteoalarmJSON is DHMZ's Meteoalarm feed with a yellow wind warning in force
// for Bjelovarsko-bilogorska and an orange one for Zadarska.
func meteoalarmJSON(now time.Time) []byte {
	warning := func(id, area, level string) map[string]any {
		return map[string]any{"alert": map[string]any{
			"identifier": id, "msgType": "Alert",
			"info": []map[string]any{{
				"language": "hr-HR", "event": "Upozorenje", "description": "Opis",
				"onset": now.Add(-time.Hour).Format(time.RFC3339), "expires": now.Add(5 * time.Hour).Format(time.RFC3339),
				"responseType": []string{"Monitor"},
				"parameter": []map[string]string{
					{"valueName": "awareness_level", "value": level},
					{"valueName": "awareness_type", "value": "1; Wind"},
				},
				"area": []map[string]any{{"areaDesc": area, "geocode": []map[string]string{{"valueName": "EMMA_ID", "value": "X"}}}},
			}},
		}}
	}
	body, _ := json.Marshal(map[string]any{"warnings": []any{
		warning("here", "Bjelovarsko-bilogorska", "2; yellow; Moderate"),
		warning("coast", "Zadarska", "3; orange; Severe"),
	}})
	return body
}
```

`backend/app/controllers/warnings_test.go`:

```go
package controllers_test

import (
	"net/http"
	"testing"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/models"
)

func TestWarningsIndexListsTheCountysWarnings(t *testing.T) {
	app, _ := newApp(t)

	daruvar := raptor.DecodeJSON[[]models.WarningResponse](t, app.TestGet("/api/v1/warnings?lat=45.59&lon=17.22"), http.StatusOK)
	zadar := raptor.DecodeJSON[[]models.WarningResponse](t, app.TestGet("/api/v1/warnings?lat=44.12&lon=15.23"), http.StatusOK)

	if len(daruvar) != 1 || daruvar[0].Level != "yellow" || daruvar[0].Type != "wind" {
		t.Errorf("Daruvar: %+v", daruvar)
	}
	if len(zadar) != 1 || zadar[0].Level != "orange" {
		t.Errorf("Zadar: %+v", zadar)
	}
}

func TestWarningsIndexIsEmptyOutsideCroatia(t *testing.T) {
	app, u := newApp(t)

	got := raptor.DecodeJSON[[]models.WarningResponse](t, app.TestGet("/api/v1/warnings?lat=43.86&lon=18.41"), http.StatusOK)

	if got == nil || len(got) != 0 {
		t.Errorf("Sarajevo: %#v, want []", got)
	}
	if n := u.hits("/api/v1/warnings/feeds-croatia"); n != 0 {
		t.Errorf("Meteoalarm asked %d times for a point outside Croatia", n)
	}
}

func TestWarningsIndexAnswers502WhenMeteoalarmFailsAndNothingIsCached(t *testing.T) {
	app, u := newApp(t)
	u.fail(http.StatusServiceUnavailable)

	raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/warnings?lat=45.59&lon=17.22"), http.StatusBadGateway)
}
```

Run: `cd backend && go test ./app/controllers/ -run Warnings`

Expected: FAIL, because the route is unknown.

- [ ] **Step 4: Implement the service, the controller, the registration and the route**

`backend/app/services/warnings_service.go`:

```go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/models"
)

const (
	meteoalarmURL  = "https://feeds.meteoalarm.org"
	meteoalarmFeed = "croatia" // DHMZ's warnings
	warningsTTL    = 15 * time.Minute
)

// WarningsService reads DHMZ's weather warnings from Meteoalarm's Croatian
// feed: one feed for the whole country, cached once. APP_METEOALARM_URL
// overrides the feed host (tests).
type WarningsService struct {
	raptor.Service

	client  *http.Client
	baseURL string
	cache   *staleCache[[]models.MeteoalarmAlert]
}

func (s *WarningsService) Setup() error {
	s.client = &http.Client{Timeout: 10 * time.Second}
	s.baseURL = strings.TrimSuffix(s.Config.AppString("meteoalarm_url", meteoalarmURL), "/")
	s.cache = newStaleCache[[]models.MeteoalarmAlert](warningsTTL)
	return nil
}

// Alerts returns every alert in the feed. When Meteoalarm fails it returns the
// last alerts it has; it errors only when it has none.
func (s *WarningsService) Alerts(ctx context.Context) ([]models.MeteoalarmAlert, error) {
	alerts, _, err := s.cache.Get(ctx, s.fetch)
	return alerts, err
}

func (s *WarningsService) fetch(ctx context.Context) ([]models.MeteoalarmAlert, error) {
	alerts, err := s.request(ctx)
	if err != nil {
		s.Log.Warn("Meteoalarm request failed", "error", err)
	}
	return alerts, err
}

func (s *WarningsService) request(ctx context.Context) ([]models.MeteoalarmAlert, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/api/v1/warnings/feeds-"+meteoalarmFeed, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meteoalarm: unexpected status %d", resp.StatusCode)
	}

	var feed models.MeteoalarmFeed
	if err := json.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("meteoalarm: %w", err)
	}
	alerts := make([]models.MeteoalarmAlert, len(feed.Warnings))
	for i, w := range feed.Warnings {
		alerts[i] = w.Alert
	}
	return alerts, nil
}
```

`backend/app/controllers/warnings_controller.go`:

```go
package controllers

import (
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/weather/app/models"
	"github.com/h00s/weather/app/services"
)

type WarningsController struct {
	raptor.Controller

	Warnings *services.WarningsService
	Places   *services.PlacesService
}

// Index lists DHMZ's warnings for the county at ?lat&lon, in force now or
// starting within 48 hours, most severe first. Outside Croatia it is empty.
func (c *WarningsController) Index(ctx *raptor.Context) error {
	at, err := coordinates(ctx)
	if err != nil {
		return err
	}
	place, ok := c.Places.Nearest(at)
	county, known := models.Counties[place.County]
	if !ok || !known {
		return ctx.Data([]models.WarningResponse{})
	}
	alerts, err := c.Warnings.Alerts(ctx.Request().Context())
	if err != nil {
		return errs.NewErrorBadGateway("Meteoalarm is unavailable")
	}
	return ctx.Data(models.ActiveWarnings(alerts, county.Areas, time.Now()))
}
```

Registration:
- `services.go`: add `&services.WarningsService{}` after `PlacesService`;
- `controllers.go`: add `&controllers.WarningsController{}`;
- `middlewares.go`: the limiter scope becomes `"Forecast", "Places", "Warnings"`;
- `routes.yaml`, under `/api/v1`:

```yaml
    /warnings:
      GET: Warnings.Index       # ?lat&lon: DHMZ warnings for that county, now or within 48 h
```

- [ ] **Step 5: Run all backend tests**

Run: `cd backend && go vet ./... && go test ./...`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add -A backend
git commit -m "Serve DHMZ warnings for the county at a location

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 7: Air quality and pollen

**Files:**
- Create, ported verbatim from dashboard: `backend/app/models/airquality.go`, `backend/app/models/airquality_test.go`
- Create: `backend/app/services/airquality_service.go`, `backend/app/controllers/airquality_controller.go`, `backend/app/controllers/airquality_test.go`
- Modify: `backend/app/controllers/fixtures_test.go`, `backend/config/components/{services,controllers,middlewares}.go`, `backend/config/routes.yaml`

**Interfaces:**
- Consumes: `newKeyedCache`, `maxCachedLocations` and `coordinates(ctx)`.
- Produces `GET /api/v1/air-quality?lat&lon`, which answers 200 `models.AirQualityResponse` (json `aqi{value,level}`, `pollen[]{type,current,todayMax,level}`, `fetchedAt`) or a 502.

- [ ] **Step 1: Port the model**

```bash
cp ~/dev/go/husak-dashboard/backend/app/models/airquality.go ~/dev/go/husak-dashboard/backend/app/models/airquality_test.go backend/app/models/
```

Run: `cd backend && go test ./app/models/`

Expected: PASS. The model is self-contained.

- [ ] **Step 2: Write the failing endpoint tests**

In `backend/app/controllers/fixtures_test.go`, add the case and the function:

```go
	case "/air-quality":
		return airQualityJSON(now), true
```

```go
// airQualityJSON is an Open-Meteo air quality answer for today: moderate AQI and
// a ragweed peak of 23.5 grains/m³ in the afternoon.
func airQualityJSON(now time.Time) []byte {
	zagreb, _ := time.LoadLocation("Europe/Zagreb")
	now = now.In(zagreb)
	_, offset := now.Zone()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, zagreb)
	var times []string
	ragweed := make([]float64, 24)
	for i := range 24 {
		times = append(times, midnight.Add(time.Duration(i)*time.Hour).Format("2006-01-02T15:04"))
	}
	ragweed[15] = 23.5

	body, _ := json.Marshal(map[string]any{
		"utc_offset_seconds":    offset,
		"timezone":              "Europe/Zagreb",
		"timezone_abbreviation": fmt.Sprintf("GMT+%d", offset/3600),
		"current": map[string]any{
			"time": now.Truncate(time.Hour).Format("2006-01-02T15:04"), "interval": 3600,
			"european_aqi": 43, "ragweed_pollen": 3.2, "birch_pollen": nil,
		},
		"hourly": map[string]any{"time": times, "ragweed_pollen": ragweed},
	})
	return body
}
```

`backend/app/controllers/airquality_test.go`:

```go
package controllers_test

import (
	"net/http"
	"testing"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/app/models"
)

func TestAirQualityShowReturnsAQIAndPollen(t *testing.T) {
	app, u := newApp(t)

	res := raptor.DecodeJSON[models.AirQualityResponse](t, app.TestGet("/api/v1/air-quality?lat=45.5936&lon=17.2251"), http.StatusOK)

	if res.AQI.Value != 43 || res.AQI.Level != "moderate" {
		t.Errorf("aqi = %+v", res.AQI)
	}
	var ragweed models.PollenResponse
	for _, p := range res.Pollen {
		if p.Type == "ragweed" {
			ragweed = p
		}
	}
	if ragweed.TodayMax != 23.5 || ragweed.Level != "high" || len(res.Pollen) != len(models.PollenTypes) {
		t.Errorf("pollen = %+v", res.Pollen)
	}
	q := u.query("/air-quality")
	if q.Get("timezone") != "auto" || q.Get("forecast_days") != "1" || q.Get("latitude") == "" {
		t.Errorf("air quality query = %v", q)
	}
}

func TestAirQualityShowAnswers502WhenOpenMeteoFailsAndNothingIsCached(t *testing.T) {
	app, u := newApp(t)
	u.fail(http.StatusServiceUnavailable)

	raptor.DecodeJSON[errorBody](t, app.TestGet("/api/v1/air-quality?lat=45.59&lon=17.22"), http.StatusBadGateway)
}
```

Run: `cd backend && go test ./app/controllers/ -run AirQuality`

Expected: FAIL, because the route is unknown.

- [ ] **Step 3: Implement the service, the controller, the registration and the route**

`backend/app/services/airquality_service.go`:

```go
package services

import (
	"context"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/goopenmeteo"
	"github.com/h00s/weather/app/models"
)

// airQualityTTL: the CAMS models behind it update a few times a day.
const airQualityTTL = 30 * time.Minute

// AirQualityService fetches today's air quality and pollen from Open-Meteo,
// cached per ~1 km cell. APP_AIRQUALITY_URL overrides the API root (tests).
type AirQualityService struct {
	raptor.Service

	client *goopenmeteo.OpenMeteo
	cache  *keyedCache[string, *goopenmeteo.AirQuality]
}

func (s *AirQualityService) Setup() error {
	s.client = goopenmeteo.NewOpenMeteo()
	s.client.AirQualityURL = s.Config.AppString("airquality_url", goopenmeteo.AirQualityBaseURL)
	s.cache = newKeyedCache[string, *goopenmeteo.AirQuality](airQualityTTL, maxCachedLocations)
	return nil
}

// AirQuality returns today's air quality at at's cell and when it was fetched.
// When Open-Meteo fails it returns the last one it has; it errors only when it
// has none.
func (s *AirQualityService) AirQuality(ctx context.Context, at models.Coordinates) (*goopenmeteo.AirQuality, time.Time, error) {
	at = at.Rounded()
	return s.cache.Get(ctx, at.Key(), func(ctx context.Context) (*goopenmeteo.AirQuality, error) {
		vars := models.AirQualityVariables()
		aq, err := s.client.AirQualityContext(ctx, goopenmeteo.AirQualityOptions{
			Latitude:     at.Latitude,
			Longitude:    at.Longitude,
			Timezone:     "auto",
			Current:      vars,
			Hourly:       vars,
			ForecastDays: 1,
		})
		if err != nil {
			s.Log.Warn("Open-Meteo air quality request failed", "cell", at.Key(), "error", err)
		}
		return aq, err
	})
}
```

`backend/app/controllers/airquality_controller.go`:

```go
package controllers

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/weather/app/models"
	"github.com/h00s/weather/app/services"
)

type AirQualityController struct {
	raptor.Controller

	AirQuality *services.AirQualityService
}

// Show is today's European AQI and pollen at ?lat&lon.
func (c *AirQualityController) Show(ctx *raptor.Context) error {
	at, err := coordinates(ctx)
	if err != nil {
		return err
	}
	aq, fetchedAt, err := c.AirQuality.AirQuality(ctx.Request().Context(), at)
	if err != nil {
		return errs.NewErrorBadGateway("Open-Meteo air quality is unavailable")
	}
	return ctx.Data(models.NewAirQualityResponse(aq, fetchedAt))
}
```

Registration:
- `services.go`: add `&services.AirQualityService{}`;
- `controllers.go`: add `&controllers.AirQualityController{}`;
- `middlewares.go`: the limiter scope becomes `"Forecast", "Places", "Warnings", "AirQuality"`;
- `routes.yaml`, under `/api/v1`:

```yaml
    /air-quality:
      GET: AirQuality.Show      # ?lat&lon: today's European AQI and pollen
```

- [ ] **Step 4: Run all backend tests**

Run: `cd backend && go vet ./... && go test ./...`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add -A backend
git commit -m "Serve air quality and pollen for any location

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 8: Frontend foundation: config, client, types and services

**Files:**
- Delete: `frontend/svelte.config.js`, `frontend/.env.production`, `frontend/src/app.css`, `frontend/src/theme.css`, `frontend/src/lib/index.ts`, `frontend/src/routes/+page.ts`, `frontend/src/lib/components/app/*`, `frontend/src/lib/components/layouts/*`, `frontend/src/lib/services/forecast.ts`, `frontend/src/lib/types/Forecast.ts`, `frontend/src/lib/utils/*`
- Modify: `frontend/package.json`, `frontend/bun.lock`, `frontend/vite.config.ts`, `frontend/tsconfig.json`, `frontend/src/app.html`, `frontend/src/routes/+layout.svelte`, `frontend/src/routes/+page.svelte`
- Create: `frontend/src/routes/layout.css`, `frontend/src/routes/+layout.ts`, `frontend/src/routes/+error.svelte`, `frontend/src/lib/api/client.ts`, `frontend/src/lib/api/client.test.ts`, `frontend/src/lib/types/{forecast,place,location,warning,airquality}.ts`, `frontend/src/lib/services/{forecast,places,warnings,airquality}.ts`

**Interfaces:**
- Produces from the client: `api.get<T>(endpoint, { fetch?, query?, signal? })`, `ApiError(message, status, body)` and `isAbortError(e)`.
- Produces the types:
  - `Forecast`, `CurrentForecast`, `HourlyForecast` and `DailyForecast` in `$lib/types/forecast`;
  - `Place` in `$lib/types/place`;
  - `Coordinates` and `SavedLocation { kind: "gps" | "place"; name: string | null; county: string | null; latitude; longitude }` in `$lib/types/location`;
  - `Warning` in `$lib/types/warning`;
  - `AirQuality`, `AQILevel`, `PollenLevel` and `PollenType` in `$lib/types/airquality`.
- Produces the services:
  - `fetchForecast(at, signal?, fetch?)`, `fetchAirQuality(at, signal?, fetch?)` and `fetchWarnings(at, signal?, fetch?)`;
  - `searchPlaces(q, signal?, fetch?)` and `nearestPlace(at, signal?, fetch?)`.
- Produces in CSS:
  - the tokens `--color-navy`, `--color-foreground`, `--color-muted-foreground`, `--color-panel`, `--color-panel-edge`, `--color-sheet`, `--color-rain` and `--color-warn-{yellow,orange,red}`;
  - the utilities `panel`, `text-lift` and `scrollbar-thin`;
  - the registered properties `--sky-top` and `--sky-bottom`.

- [ ] **Step 1: Remove the old app**

```bash
cd frontend
git rm -q svelte.config.js .env.production src/app.css src/theme.css src/lib/index.ts src/routes/+page.ts \
  src/lib/components/app/CurrentForecast.svelte src/lib/components/app/DailyForecast.svelte \
  src/lib/components/app/HourlyForecast.svelte src/lib/components/layouts/Container.svelte \
  src/lib/services/forecast.ts src/lib/types/Forecast.ts src/lib/utils/datetime.ts src/lib/utils/weather.ts
```

- [ ] **Step 2: Update the dependencies**

```bash
cd frontend
bun remove @tailwindcss/vite tailwindcss
bun add -d @sveltejs/adapter-static@^3.0.10 @sveltejs/kit@^2.70.3 @sveltejs/vite-plugin-svelte@^7.3.1 \
  svelte@^5.57.1 svelte-check@^4.7.6 typescript@~6 vite@^8.3.1 vitest@^5.0.2 \
  @tailwindcss/vite@^4.3.3 tailwindcss@^4.3.3 @fontsource-variable/inter@^5.3.0 \
  @bybas/weather-icons@^2.0.0 @lucide/svelte@latest
```

In `frontend/package.json`, add `"test": "vitest run"` to `scripts`. `dependencies` should now be gone; everything is a devDependency, as in dashboard.

- [ ] **Step 3: Configure Vite, SvelteKit and TypeScript**

`frontend/vite.config.ts`:

```ts
import tailwindcss from "@tailwindcss/vite";
import adapter from "@sveltejs/adapter-static";
import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vitest/config";

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit({
      // SPA: every route falls back to the client shell, which the Go server serves.
      adapter: adapter({ fallback: "index.html" }),
      // A phone keeps the app open for days and rarely navigates: poll for a new deploy, and the
      // root layout reloads into it the next time the page is shown.
      version: { pollInterval: 15 * 60_000 },
      compilerOptions: {
        // Runes everywhere except in libraries (the default in Svelte 6).
        runes: ({ filename }) => (filename.split(/[/\\]/).includes("node_modules") ? undefined : true),
      },
    }),
  ],
  test: {
    include: ["src/**/*.test.ts"],
  },
  server: {
    // Same origin in development too: the browser talks only to Vite, which forwards /api to the
    // Go server. The object form keeps Host intact, which the backend's csrf middleware checks.
    proxy: {
      "/api": { target: "http://localhost:3000" },
    },
  },
});
```

```bash
cp ~/dev/go/husak-dashboard/frontend/tsconfig.json frontend/tsconfig.json
```

- [ ] **Step 4: Write the shell, the styles and the layouts**

`frontend/src/app.html`:

```html
<!doctype html>
<html lang="hr">
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover" />
		<meta
			name="description"
			content="Vremenska prognoza za vaše mjesto: sada, po satima i za 7 dana, s upozorenjima DHMZ-a, kvalitetom zraka i peludi."
		/>
		<meta name="theme-color" content="#012a4a" />
		<link rel="icon" href="%sveltekit.assets%/favicon.png" />
		%sveltekit.head%
	</head>
	<body data-sveltekit-preload-data="hover">
		<div style="display: contents">%sveltekit.body%</div>
	</body>
</html>
```

`frontend/src/routes/layout.css`:

```css
@import "tailwindcss";
@import "@fontsource-variable/inter";

@theme {
  --font-sans: "Inter Variable", ui-sans-serif, system-ui, sans-serif;

  /* The brand navy: the night sky, the PWA splash, and what shows before the sky is computed. */
  --color-navy: #012a4a;
  --color-foreground: oklch(0.985 0 0);
  --color-muted-foreground: oklch(0.985 0 0 / 75%);

  /* Frosted panels over the sky */
  --color-panel: oklch(0.22 0.05 248 / 35%);
  --color-panel-edge: oklch(1 0 0 / 10%);
  /* The location sheet sits above everything, so it is nearly opaque. */
  --color-sheet: oklch(0.25 0.06 248 / 97%);

  /* Data marks */
  --color-rain: oklch(0.82 0.1 230);

  /* Warning levels: always shown with an icon and text */
  --color-warn-yellow: oklch(0.88 0.16 95);
  --color-warn-orange: oklch(0.77 0.16 60);
  --color-warn-red: oklch(0.6 0.21 27);
}

/* Registered so the sky's colours can transition; a gradient itself can't. */
@property --sky-top {
  syntax: "<color>";
  inherits: false;
  initial-value: #012a4a;
}
@property --sky-bottom {
  syntax: "<color>";
  inherits: false;
  initial-value: #012a4a;
}

@utility panel {
  border-radius: var(--radius-3xl);
  background-color: var(--color-panel);
  box-shadow: inset 0 0 0 1px var(--color-panel-edge);
  -webkit-backdrop-filter: blur(12px);
  backdrop-filter: blur(12px);
}

/* Keeps white text crisp against the sky: a tight shadow for the edges, a soft one for the glow. */
@utility text-lift {
  text-shadow:
    0 1px 2px oklch(0 0 0 / 40%),
    0 2px 18px oklch(0 0 0 / 30%);
}

/* Horizontal strips: a thin, quiet scrollbar where the platform shows one. */
@utility scrollbar-thin {
  scrollbar-width: thin;
  scrollbar-color: oklch(1 0 0 / 25%) transparent;
}

@layer base {
  :root {
    color-scheme: dark;
  }

  html {
    background-color: var(--color-navy);
    /* The page's own pull-to-refresh replaces the browser's reload-on-pull. */
    overscroll-behavior-y: contain;
  }

  body {
    @apply text-foreground font-sans antialiased;
    min-height: 100dvh;
  }

  :focus-visible {
    outline: 2px solid white;
    outline-offset: 2px;
  }
}
```

`frontend/src/routes/+layout.ts`:

```ts
/** A client-only SPA (adapter-static): the API exists only in the browser. */
export const ssr = false;
```

`frontend/src/routes/+layout.svelte`:

```svelte
<script lang="ts">
  import "./layout.css";
  import { updated } from "$app/state";
  import type { LayoutProps } from "./$types";

  let { children }: LayoutProps = $props();

  // SvelteKit applies a new deploy on the next navigation, and this one-page app rarely navigates.
  // When the version poll (vite.config.ts) sees a new build, reload the next time the page comes
  // back into view: never under the reader's eyes. A full reload is the point here, hence location.
  $effect(() => {
    if (!updated.current) return;
    const reload = () => {
      if (document.visibilityState === "visible") location.reload();
    };
    document.addEventListener("visibilitychange", reload);
    return () => document.removeEventListener("visibilitychange", reload);
  });
</script>

{@render children()}
```

`frontend/src/routes/+error.svelte`:

```svelte
<script lang="ts">
  import { page } from "$app/state";
</script>

<svelte:head><title>Pogreška {page.status} · Vrijeme</title></svelte:head>

<main class="flex min-h-dvh flex-col items-center justify-center gap-4 px-4 text-center">
  <p class="text-muted-foreground text-lg">Pogreška {page.status}</p>
  <h1 class="text-3xl font-semibold">{page.status === 404 ? "Ova stranica ne postoji." : "Nešto je pošlo po zlu."}</h1>
  <a href="/" class="panel px-5 py-2.5 font-medium">Na prognozu</a>
</main>
```

`frontend/src/routes/+page.svelte` is a placeholder until Task 12:

```svelte
<svelte:head><title>Vrijeme</title></svelte:head>

<main class="grid min-h-dvh place-items-center text-3xl font-extralight">Vrijeme</main>
```

- [ ] **Step 5: Write the failing API client test**

`frontend/src/lib/api/client.test.ts`:

```ts
import { describe, expect, it, vi } from "vitest";
import { api, ApiError, isAbortError } from "./client";

const respond = (status: number, body?: unknown) =>
  vi.fn<typeof globalThis.fetch>(
    async () =>
      new Response(body === undefined ? null : JSON.stringify(body), {
        status,
        headers: { "Content-Type": "application/json" },
      }),
  );

describe("api.get", () => {
  it("calls /api/v1 with the query, dropping empty values", async () => {
    const fetch = respond(200, { ok: true });

    const got = await api.get("/forecast", { fetch, query: { lat: 45.59, lon: 17.23, q: null, x: undefined } });

    expect(got).toEqual({ ok: true });
    expect(fetch.mock.calls[0][0]).toBe("/api/v1/forecast?lat=45.59&lon=17.23");
  });

  it("throws an ApiError carrying the status and Raptor's envelope", async () => {
    const fetch = respond(502, { code: 502, message: "Open-Meteo is unavailable" });

    const error = await api.get("/forecast", { fetch }).catch((e: unknown) => e);

    expect(error).toBeInstanceOf(ApiError);
    expect((error as ApiError).status).toBe(502);
    expect((error as ApiError).body).toEqual({ code: 502, message: "Open-Meteo is unavailable" });
  });

  it("recognizes the caller's own abort", () => {
    expect(isAbortError(new DOMException("aborted", "AbortError"))).toBe(true);
    expect(isAbortError(new Error("network"))).toBe(false);
  });
});
```

Run: `cd frontend && bun run test`

Expected: FAIL, because `./client` does not exist.

- [ ] **Step 6: Port the client**

```bash
mkdir -p frontend/src/lib/api
cp ~/dev/go/husak-dashboard/frontend/src/lib/api/client.ts frontend/src/lib/api/client.ts
sed -i 's|// The kiosk only reads, and has no session|// The app only reads, and has no session|' frontend/src/lib/api/client.ts
```

Run: `cd frontend && bun run test`

Expected: PASS (3 tests).

- [ ] **Step 7: Write the types and services**

`frontend/src/lib/types/forecast.ts`:

```ts
/** Mirrors the backend's ForecastResponse (app/models/forecast.go). Times are ISO 8601 with the
 *  location's offset; dates are YYYY-MM-DD on the location's calendar. */
export interface CurrentForecast {
  time: string;
  /** °C */
  temperature: number;
  /** °C */
  apparentTemperature: number;
  /** % */
  humidity: number;
  /** °C */
  dewPoint: number;
  /** mm */
  precipitation: number;
  /** WMO weather interpretation code */
  weatherCode: number;
  /** % */
  cloudCover: number;
  /** hPa at sea level */
  pressure: number;
  /** km/h */
  windSpeed: number;
  /** Degrees, where the wind blows from */
  windDirection: number;
  /** km/h */
  windGusts: number;
  /** m */
  visibility: number;
  uvIndex: number;
  isDay: boolean;
}

export interface HourlyForecast {
  time: string;
  temperature: number;
  apparentTemperature: number;
  /** % */
  precipitationProbability: number;
  /** mm */
  precipitation: number;
  weatherCode: number;
  windSpeed: number;
  windDirection: number;
  isDay: boolean;
}

export interface DailyForecast {
  /** YYYY-MM-DD on the location's calendar */
  date: string;
  temperatureMin: number;
  temperatureMax: number;
  weatherCode: number;
  precipitationSum: number;
  precipitationProbabilityMax: number;
  windSpeedMax: number;
  windGustsMax: number;
  windDirectionDominant: number;
  uvIndexMax: number;
  sunrise: string;
  sunset: string;
  daylightSeconds: number;
}

export interface Forecast {
  current: CurrentForecast;
  /** From the current hour to the end of the last day. */
  hourly: HourlyForecast[];
  /** Seven days from today. */
  daily: DailyForecast[];
  /** The location's IANA zone: every time on the page is shown in it. */
  timezone: string;
  /** Metres */
  elevation: number;
  /** When the backend fetched it; much older than 10 minutes means Open-Meteo is failing. */
  fetchedAt: string;
}
```

`frontend/src/lib/types/place.ts`:

```ts
/** Mirrors the backend's PlaceResponse (app/models/place.go). */
export interface Place {
  name: string;
  /** e.g. "Bjelovarsko-bilogorska županija"; "" when unknown */
  county: string;
  latitude: number;
  longitude: number;
}
```

`frontend/src/lib/types/location.ts`:

```ts
export interface Coordinates {
  latitude: number;
  longitude: number;
}

/** What the app shows weather for, as stored on this device. Coordinates are rounded to 2 decimals
 *  (~1 km) before they are stored or sent. */
export interface SavedLocation extends Coordinates {
  kind: "gps" | "place";
  /** A settlement; null for a position with no Croatian place nearby ("Moja lokacija"). */
  name: string | null;
  county: string | null;
}
```

```bash
cp ~/dev/go/husak-dashboard/frontend/src/lib/types/warning.ts ~/dev/go/husak-dashboard/frontend/src/lib/types/airquality.ts frontend/src/lib/types/
```

`frontend/src/lib/services/forecast.ts`:

```ts
import { api } from "$lib/api/client";
import type { Forecast } from "$lib/types/forecast";
import type { Coordinates } from "$lib/types/location";

export const fetchForecast = (at: Coordinates, signal?: AbortSignal, fetch?: typeof globalThis.fetch) =>
  api.get<Forecast>("/forecast", { query: { lat: at.latitude, lon: at.longitude }, signal, fetch });
```

`frontend/src/lib/services/airquality.ts`:

```ts
import { api } from "$lib/api/client";
import type { AirQuality } from "$lib/types/airquality";
import type { Coordinates } from "$lib/types/location";

export const fetchAirQuality = (at: Coordinates, signal?: AbortSignal, fetch?: typeof globalThis.fetch) =>
  api.get<AirQuality>("/air-quality", { query: { lat: at.latitude, lon: at.longitude }, signal, fetch });
```

`frontend/src/lib/services/warnings.ts`:

```ts
import { api } from "$lib/api/client";
import type { Coordinates } from "$lib/types/location";
import type { Warning } from "$lib/types/warning";

export const fetchWarnings = (at: Coordinates, signal?: AbortSignal, fetch?: typeof globalThis.fetch) =>
  api.get<Warning[]>("/warnings", { query: { lat: at.latitude, lon: at.longitude }, signal, fetch });
```

`frontend/src/lib/services/places.ts`:

```ts
import { api } from "$lib/api/client";
import type { Coordinates } from "$lib/types/location";
import type { Place } from "$lib/types/place";

/** Croatian settlements by name (2 to 60 characters), best match first. */
export const searchPlaces = (q: string, signal?: AbortSignal, fetch?: typeof globalThis.fetch) =>
  api.get<Place[]>("/places", { query: { q }, signal, fetch });

/** The settlement at a point; rejects with an ApiError 404 abroad or at sea. */
export const nearestPlace = (at: Coordinates, signal?: AbortSignal, fetch?: typeof globalThis.fetch) =>
  api.get<Place>("/places/nearest", { query: { lat: at.latitude, lon: at.longitude }, signal, fetch });
```

- [ ] **Step 8: Check, test and build**

Run: `cd frontend && bun run check && bun run test && bun run build`

Expected: 0 errors and 0 warnings, 3 tests passing, and `build/index.html` written.

- [ ] **Step 9: Commit**

```bash
git add -A frontend
git commit -m "Modernize the frontend: kit config in Vite, same-origin API client, typed services

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 9: Formatting, weather and sky helpers

**Files:**
- Create: `frontend/src/lib/helpers/{format,weather,alerts,preview,sky}.ts`, each with its `*.test.ts`
- Create, ported verbatim from dashboard: `frontend/src/lib/helpers/sparkline.ts`, `frontend/src/lib/helpers/sparkline.test.ts`

**Interfaces:**
- Consumes: the types from Task 8.
- Produces. Every time-related function takes the location's `timeZone: string`.
  - From `format`: `formatClock(date, tz)`, `formatHour(date, tz)`, `isoDateOf(date, tz)`, `minutesOfDay(date, tz)`, `dayLabel(date: string, now, tz)`, `formatDegrees(v)`, `formatPercent(v)`, `formatMillimetres(v)`, `formatSpeed(kmh)`, `formatPressure(hpa)`, `formatDistance(m)`, `formatDuration(s)` and `formatAgo(then, now)`.
  - From `weather`: `weatherInfo(code, isDay?) → { label, icon, precipitation }`, `type Precipitation`, `type WeatherInfo` and `rainHint(hourly, now, tz)`.
  - From `alerts`: `warningWhen(w, now, tz)`.
  - From `preview`: `parsePreview(params) → { at?, code?, cloud? }` and `previewOffset(at, now, tz)`.
  - From `sky`: `type Oklch`, `type Sky`, `skyAt(now, sunrise, sunset, cloudCover, precipitation)`, `oklchCss(c)`, `oklchToHex(c)`, `contrastWithWhite(c)` and `MIN_CONTRAST`.
  - From `sparkline`: `extent`, `scale`, `linePath`, `areaPath` and `type Series`.

- [ ] **Step 1: Port the sparkline helpers**

```bash
mkdir -p frontend/src/lib/helpers
cp ~/dev/go/husak-dashboard/frontend/src/lib/helpers/sparkline.ts ~/dev/go/husak-dashboard/frontend/src/lib/helpers/sparkline.test.ts frontend/src/lib/helpers/
```

- [ ] **Step 2: Write the failing format tests**

`frontend/src/lib/helpers/format.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import {
  dayLabel,
  formatAgo,
  formatClock,
  formatDegrees,
  formatDistance,
  formatDuration,
  formatHour,
  formatMillimetres,
  formatPercent,
  formatPressure,
  formatSpeed,
  isoDateOf,
  minutesOfDay,
} from "./format";

const tz = "Europe/Zagreb";
const d = new Date("2026-09-29T07:05:00+02:00");

describe("times on the location's clock", () => {
  it("formats clock times and hours in the given zone", () => {
    expect(formatClock(d, tz)).toBe("07:05");
    expect(formatHour(d, tz)).toBe("7 h");
    expect(formatClock(d, "UTC")).toBe("05:05");
    expect(minutesOfDay(d, tz)).toBe(425);
  });

  it("dates by the location's calendar, not the device's", () => {
    expect(isoDateOf(new Date("2026-09-29T23:30:00Z"), tz)).toBe("2026-09-30");
  });

  it("names days relative to today, then by weekday and date", () => {
    expect(dayLabel("2026-09-29", d, tz)).toBe("Danas");
    expect(dayLabel("2026-09-30", d, tz)).toBe("Sutra");
    expect(dayLabel("2026-10-01", d, tz)).toBe("Čet 1. 10.");
    expect(dayLabel("2026-10-04", d, tz)).toBe("Ned 4. 10.");
  });
});

describe("values", () => {
  it("rounds degrees, never showing -0", () => {
    expect(formatDegrees(17.46)).toBe("17°");
    expect(formatDegrees(-0.4)).toBe("0°");
    expect(formatDegrees(null)).toBe("–");
  });

  it("formats measurements the Croatian way", () => {
    expect(formatPercent(47.6)).toBe("48 %");
    expect(formatMillimetres(4.25)).toBe("4,3 mm");
    expect(formatMillimetres(0)).toBe("0 mm");
    expect(formatSpeed(12.4)).toBe("12 km/h");
    expect(formatPressure(1013.2)).toBe("1.013 hPa");
    expect(formatDistance(41_300)).toBe("41 km");
    expect(formatDistance(2_500)).toBe("2,5 km");
    expect(formatDistance(820)).toBe("800 m");
    expect(formatDuration(41_590)).toBe("11 h 33 min");
  });

  it("says how long ago", () => {
    const now = new Date("2026-09-29T12:00:00Z");
    const ago = (ms: number) => formatAgo(new Date(now.getTime() - ms), now);
    expect(ago(20_000)).toBe("upravo");
    expect(ago(5 * 60_000)).toBe("prije 5 min");
    expect(ago(125 * 60_000)).toBe("prije 2 h");
    expect(ago(24 * 3_600_000)).toBe("prije 1 dan");
    expect(ago(3 * 24 * 3_600_000)).toBe("prije 3 dana");
  });
});
```

Run: `cd frontend && bun run test`

Expected: FAIL, because `./format` does not exist.

- [ ] **Step 3: Implement format**

`frontend/src/lib/helpers/format.ts`:

```ts
// Times and dates are shown on the forecast location's clock and calendar: pass its IANA zone
// (Forecast.timezone), never the device's.

const dateFormats = new Map<string, Intl.DateTimeFormat>();

/** An Intl.DateTimeFormat, built once per locale and options: constructing one is slow. */
function dateFormat(locale: string, options: Intl.DateTimeFormatOptions): Intl.DateTimeFormat {
  const key = `${locale}|${JSON.stringify(options)}`;
  let format = dateFormats.get(key);
  if (!format) dateFormats.set(key, (format = new Intl.DateTimeFormat(locale, options)));
  return format;
}

const decimal = new Intl.NumberFormat("hr-HR", { maximumFractionDigits: 1 });
const integer = new Intl.NumberFormat("hr-HR", { maximumFractionDigits: 0 });

/** 14:27 */
export const formatClock = (date: Date, timeZone: string) =>
  dateFormat("hr-HR", { hour: "2-digit", minute: "2-digit", hourCycle: "h23", timeZone }).format(date);

/** 7 h */
export const formatHour = (date: Date, timeZone: string) =>
  `${Number.parseInt(dateFormat("hr-HR", { hour: "numeric", hourCycle: "h23", timeZone }).format(date), 10)} h`;

/** YYYY-MM-DD on the location's calendar. */
export const isoDateOf = (date: Date, timeZone: string) => dateFormat("en-CA", { timeZone }).format(date);

/** Minutes since midnight on the location's clock. */
export function minutesOfDay(date: Date, timeZone: string): number {
  const [h, m] = formatClock(date, timeZone).split(":").map(Number);
  return h * 60 + m;
}

/** "Danas", "Sutra", then the short weekday and date ("Čet 1. 10.") for a YYYY-MM-DD date. */
export function dayLabel(date: string, now: Date, timeZone: string): string {
  const days = Math.round((Date.parse(date) - Date.parse(isoDateOf(now, timeZone))) / 86_400_000);
  if (days === 0) return "Danas";
  if (days === 1) return "Sutra";
  const weekday = dateFormat("hr-HR", { weekday: "short", timeZone: "UTC" })
    .format(new Date(`${date}T12:00:00Z`))
    .replace(".", "");
  const [, month, day] = date.split("-").map(Number); // Intl pads them: "01. 10."
  return `${weekday.charAt(0).toUpperCase()}${weekday.slice(1)} ${day}. ${month}.`;
}

/** 17°, rounded; never "-0°"; a dash when missing. */
export const formatDegrees = (value: number | null | undefined) =>
  value == null ? "–" : `${Math.round(value) || 0}°`;

/** 48 % */
export const formatPercent = (value: number) => `${Math.round(value)} %`;

/** 4,3 mm */
export const formatMillimetres = (value: number) => `${decimal.format(value)} mm`;

/** 12 km/h */
export const formatSpeed = (kmh: number) => `${Math.round(kmh)} km/h`;

/** 1.013 hPa */
export const formatPressure = (hpa: number) => `${integer.format(hpa)} hPa`;

/** 41 km, 2,5 km, 800 m */
export function formatDistance(metres: number): string {
  if (metres < 1000) return `${Math.round(metres / 100) * 100} m`;
  return `${(metres < 10_000 ? decimal : integer).format(metres / 1000)} km`;
}

/** 11 h 33 min */
export function formatDuration(seconds: number): string {
  const minutes = Math.round(seconds / 60);
  return `${Math.floor(minutes / 60)} h ${minutes % 60} min`;
}

/** "upravo", "prije 5 min", "prije 2 h", "prije 3 dana" */
export function formatAgo(then: Date, now: Date): string {
  const minutes = Math.floor((now.getTime() - then.getTime()) / 60_000);
  if (minutes < 1) return "upravo";
  if (minutes < 60) return `prije ${minutes} min`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `prije ${hours} h`;
  const days = Math.floor(hours / 24);
  return `prije ${days} ${days % 10 === 1 && days % 100 !== 11 ? "dan" : "dana"}`;
}
```

Run: `cd frontend && bun run test`

Expected: the format and sparkline tests PASS.

- [ ] **Step 4: Write the failing tests for weather, alerts and preview**

`frontend/src/lib/helpers/weather.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import type { HourlyForecast } from "$lib/types/forecast";
import { rainHint, weatherInfo } from "./weather";

const tz = "Europe/Zagreb";

describe("weatherInfo", () => {
  it("labels WMO codes in Croatian with a day or night icon", () => {
    expect(weatherInfo(0, true)).toEqual({ label: "Vedro", icon: "clear-day", precipitation: "none" });
    expect(weatherInfo(0, false).icon).toBe("clear-night");
    expect(weatherInfo(3, true)).toMatchObject({ label: "Oblačno", icon: "overcast-day" });
    expect(weatherInfo(63, true)).toMatchObject({ label: "Kiša", icon: "rain", precipitation: "rain" });
    expect(weatherInfo(73, false)).toMatchObject({ label: "Snijeg", icon: "snow", precipitation: "snow" });
    expect(weatherInfo(95, true)).toMatchObject({ label: "Grmljavina", precipitation: "storm" });
    expect(weatherInfo(51, true).precipitation).toBe("drizzle");
  });

  it("falls back for an unknown code", () => {
    expect(weatherInfo(42, true)).toEqual({ label: "", icon: "not-available", precipitation: "none" });
  });
});

describe("rainHint", () => {
  const start = new Date("2026-09-29T12:00:00+02:00");
  const hours = (wet: number[], code = 61): HourlyForecast[] =>
    Array.from({ length: 24 }, (_, i) => ({
      time: new Date(start.getTime() + i * 3_600_000).toISOString(),
      temperature: 15,
      apparentTemperature: 14,
      precipitationProbability: wet.includes(i) ? 80 : 5,
      precipitation: wet.includes(i) ? 1.2 : 0,
      weatherCode: wet.includes(i) ? code : 2,
      windSpeed: 10,
      windDirection: 180,
      isDay: true,
    }));
  const now = new Date("2026-09-29T12:20:00+02:00");

  it("says when rain starts", () => {
    expect(rainHint(hours([4, 5]), now, tz)).toBe("Kiša od 16 h");
  });

  it("says when ongoing rain stops", () => {
    expect(rainHint(hours([0, 1, 2]), now, tz)).toBe("Kiša do 15 h");
  });

  it("names snow and storms", () => {
    expect(rainHint(hours([3], 73), now, tz)).toBe("Snijeg od 15 h");
    expect(rainHint(hours([3], 95), now, tz)).toBe("Grmljavina od 15 h");
  });

  it("says when there is nothing for 12 hours", () => {
    expect(rainHint(hours([13]), now, tz)).toBe("Bez oborina idućih 12 h");
  });

  it("says rain continues when it doesn't stop within 12 hours", () => {
    expect(rainHint(hours([...Array(24).keys()]), now, tz)).toBe("Kiša i dalje");
  });

  it("reads hours on the location's clock", () => {
    expect(rainHint(hours([4, 5]), now, "UTC")).toBe("Kiša od 14 h");
  });
});
```

`frontend/src/lib/helpers/alerts.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import type { Warning } from "$lib/types/warning";
import { warningWhen } from "./alerts";

describe("warningWhen", () => {
  const tz = "Europe/Zagreb";
  const w = (onset: string, expires: string) => ({ onset, expires }) as Warning;
  const now = new Date("2026-09-29T12:00:00+02:00");

  it("says until when a warning in force lasts", () => {
    expect(warningWhen(w("2026-09-29T10:00:00+02:00", "2026-09-29T17:59:59+02:00"), now, tz)).toBe("do 18 h");
  });

  it("says from when an upcoming warning starts, naming tomorrow", () => {
    expect(warningWhen(w("2026-09-29T16:00:00+02:00", "2026-09-29T22:00:00+02:00"), now, tz)).toBe("od 16 h");
    expect(warningWhen(w("2026-09-30T06:00:00+02:00", "2026-09-30T12:00:00+02:00"), now, tz)).toBe("sutra od 6 h");
  });

  it("names the day when a warning in force ends tomorrow", () => {
    expect(warningWhen(w("2026-09-29T10:00:00+02:00", "2026-09-30T05:59:59+02:00"), now, tz)).toBe("do sutra 6 h");
  });

  it("names a later day by weekday and date", () => {
    expect(warningWhen(w("2026-10-01T06:00:00+02:00", "2026-10-01T12:00:00+02:00"), now, tz)).toBe("čet 1. 10. od 6 h");
  });
});
```

`frontend/src/lib/helpers/preview.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import { parsePreview, previewOffset } from "./preview";

describe("parsePreview", () => {
  it("reads a time, a WMO code and a cloud cover", () => {
    expect(parsePreview(new URLSearchParams("at=19:40&code=61&cloud=90"))).toEqual({ at: 1180, code: 61, cloud: 90 });
  });

  it("ignores what doesn't parse", () => {
    expect(parsePreview(new URLSearchParams("at=25:00&code=abc&cloud=150"))).toEqual({});
    expect(parsePreview(new URLSearchParams(""))).toEqual({});
  });
});

describe("previewOffset", () => {
  it("moves now to the asked time on the location's clock", () => {
    expect(previewOffset(1180, new Date("2026-09-29T19:00:00+02:00"), "Europe/Zagreb")).toBe(40 * 60_000);
  });
});
```

Run: `cd frontend && bun run test`

Expected: FAIL, because `./weather`, `./alerts` and `./preview` don't exist.

- [ ] **Step 5: Implement weather, alerts and preview**

```bash
cp ~/dev/go/husak-dashboard/frontend/src/lib/helpers/weather.ts frontend/src/lib/helpers/weather.ts
```

Then edit `frontend/src/lib/helpers/weather.ts`:
- replace the import `import type { HourlyWeather } from "$lib/types/weather";` with `import type { HourlyForecast } from "$lib/types/forecast";`, and rename every `HourlyWeather` to `HourlyForecast`;
- give `rainHint` a `timeZone` parameter and pass it to both `formatHour` calls, so the function reads:

```ts
export function rainHint(hourly: HourlyForecast[], now: Date, timeZone: string): string {
  const upcoming = hourly.filter((h) => Date.parse(h.time) + 3_600_000 > now.getTime()).slice(0, HORIZON);
  if (upcoming.length === 0) return "";

  const kind = (h: HourlyForecast) => noun[weatherInfo(h.weatherCode).precipitation];

  if (isWet(upcoming[0])) {
    const dry = upcoming.find((h) => !isWet(h));
    return dry ? `${kind(upcoming[0])} do ${formatHour(new Date(dry.time), timeZone)}` : `${kind(upcoming[0])} i dalje`;
  }
  const wet = upcoming.find(isWet);
  return wet ? `${kind(wet)} od ${formatHour(new Date(wet.time), timeZone)}` : `Bez oborina idućih ${HORIZON} h`;
}
```

`frontend/src/lib/helpers/alerts.ts`:

```ts
import type { Warning } from "$lib/types/warning";
import { dayLabel, formatHour, isoDateOf } from "./format";

/** Rounded to the nearest hour, with the day when it isn't today: ["sutra", "6 h"]. */
function when(iso: string, now: Date, timeZone: string): [day: string, hour: string] {
  const t = new Date(Math.round(Date.parse(iso) / 3_600_000) * 3_600_000);
  const day = dayLabel(isoDateOf(t, timeZone), now, timeZone);
  return [day === "Danas" ? "" : day.toLowerCase(), formatHour(t, timeZone)];
}

/** When a warning applies: until when if it is in force, from when if it is still to come. */
export function warningWhen(w: Warning, now: Date, timeZone: string): string {
  if (Date.parse(w.onset) <= now.getTime()) {
    const [day, hour] = when(w.expires, now, timeZone);
    return day ? `do ${day} ${hour}` : `do ${hour}`;
  }
  const [day, hour] = when(w.onset, now, timeZone);
  return day ? `${day} od ${hour}` : `od ${hour}`;
}
```

`frontend/src/lib/helpers/preview.ts`:

```ts
import { minutesOfDay } from "./format";

/** Design preview: `?at=19:40&code=61&cloud=90` shows the page at another time of day and in other
 *  weather, without waiting for it. */
export interface Preview {
  /** Minutes since midnight to show. */
  at?: number;
  /** A WMO weather code to show. */
  code?: number;
  /** Cloud cover in percent to show. */
  cloud?: number;
}

export function parsePreview(params: URLSearchParams): Preview {
  const preview: Preview = {};
  const at = /^(\d{1,2}):(\d{2})$/.exec(params.get("at") ?? "");
  if (at && Number(at[1]) < 24 && Number(at[2]) < 60) preview.at = Number(at[1]) * 60 + Number(at[2]);

  const code = Number.parseInt(params.get("code") ?? "", 10);
  if (Number.isInteger(code) && code >= 0 && code < 100) preview.code = code;

  const cloud = Number(params.get("cloud") ?? Number.NaN);
  if (params.has("cloud") && cloud >= 0 && cloud <= 100) preview.cloud = cloud;
  return preview;
}

/** Milliseconds to add to now so the location's clock reads `at` (minutes since midnight) today. */
export function previewOffset(at: number, now: Date, timeZone: string): number {
  return (at - minutesOfDay(now, timeZone)) * 60_000;
}
```

Run: `cd frontend && bun run test`

Expected: PASS.

- [ ] **Step 6: Write the failing sky tests**

`frontend/src/lib/helpers/sky.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import { contrastWithWhite, oklchCss, oklchToHex, skyAt } from "./sky";

const at = (hhmm: string) => new Date(`2026-09-29T${hhmm}:00+02:00`);
const sunrise = at("06:46");
const sunset = at("18:35");

describe("skyAt", () => {
  const clear = (hhmm: string) => skyAt(at(hhmm), sunrise, sunset, 0, "none");

  it("is darker at night than by day", () => {
    expect(clear("02:00").top.l).toBeLessThan(clear("12:00").top.l - 0.2);
  });

  it("stays blue by day and warms the horizon only around sunrise and sunset", () => {
    const hue = (h: number) => (h > 180 ? h - 360 : h); // warm hues sit near 0–90
    expect(hue(clear("06:30").bottom.h)).toBeLessThan(90);
    expect(hue(clear("18:40").bottom.h)).toBeLessThan(90);
    expect(clear("12:00").top.h).toBeGreaterThan(200);
    expect(clear("12:00").bottom.h).toBeGreaterThan(200);
    expect(clear("02:00").bottom.h).toBeGreaterThan(200);
  });

  it("changes gradually, not in steps", () => {
    const a = clear("06:40").bottom;
    const b = clear("06:41").bottom;
    expect(Math.abs(a.l - b.l)).toBeLessThan(0.02);
  });

  it("greys out under cloud and darkens in rain", () => {
    const sunny = skyAt(at("12:00"), sunrise, sunset, 0, "none");
    const overcast = skyAt(at("12:00"), sunrise, sunset, 100, "none");
    const rainy = skyAt(at("12:00"), sunrise, sunset, 100, "rain");
    expect(overcast.top.c).toBeLessThan(sunny.top.c / 2);
    expect(rainy.top.l).toBeLessThan(overcast.top.l);
  });

  it("formats as CSS", () => {
    expect(oklchCss({ l: 0.6, c: 0.13, h: 240 })).toBe("oklch(0.600 0.130 240.0)");
  });
});

describe("oklchToHex", () => {
  it("converts back to the brand navy", () => {
    const hex = oklchToHex({ l: 0.278, c: 0.073, h: 247.7 }); // #012a4a
    const [r, g, b] = [1, 3, 5].map((i) => Number.parseInt(hex.slice(i, i + 2), 16));
    expect(hex).toMatch(/^#[0-9a-f]{6}$/);
    expect(Math.abs(r - 0x01)).toBeLessThanOrEqual(2);
    expect(Math.abs(g - 0x2a)).toBeLessThanOrEqual(2);
    expect(Math.abs(b - 0x4a)).toBeLessThanOrEqual(2);
  });
});

describe("contrast with white text", () => {
  it("measures WCAG contrast", () => {
    expect(contrastWithWhite({ l: 0, c: 0, h: 0 })).toBeCloseTo(21, 0);
    expect(contrastWithWhite({ l: 1, c: 0, h: 0 })).toBeCloseTo(1, 1);
  });

  it("stays at 4.5:1 or more against every sky of the day, in any weather", () => {
    const worst = { contrast: Infinity, at: "" };
    for (let m = 0; m < 24 * 60; m += 5) {
      const now = new Date(at("00:00").getTime() + m * 60_000);
      for (const cloud of [0, 50, 100]) {
        for (const p of ["none", "snow"] as const) {
          const sky = skyAt(now, sunrise, sunset, cloud, p);
          for (const c of [sky.top, sky.bottom]) {
            const contrast = contrastWithWhite(c);
            if (contrast < worst.contrast) Object.assign(worst, { contrast, at: `${m} min, cloud ${cloud}, ${p}` });
          }
        }
      }
    }
    expect(worst.contrast, worst.at).toBeGreaterThanOrEqual(4.5);
  });

  it("keeps the daytime sky colourful rather than muddy", () => {
    expect(skyAt(at("12:00"), sunrise, sunset, 0, "none").top.c).toBeGreaterThan(0.08);
  });
});
```

Run: `cd frontend && bun run test`

Expected: FAIL, because `./sky` does not exist.

- [ ] **Step 7: Implement the sky, retuned to the brand blues**

`frontend/src/lib/helpers/sky.ts`:

```ts
import type { Precipitation } from "./weather";

export interface Oklch {
  l: number;
  c: number;
  h: number;
}

export interface Sky {
  top: Oklch;
  bottom: Oklch;
}

const MIN = 60_000;

const oklch = (l: number, c: number, h: number): Oklch => ({ l, c, h });

// The sky over a day in the brand's navy-to-teal family: [anchor, minutes from it, zenith, horizon].
// Night settles on the brand navy (#012a4a), the day sits between #01497c and #2c7da0, and only the
// hour around sunrise and sunset warms the horizon. skyAt darkens anything that would leave white
// text below MIN_CONTRAST.
const NIGHT: Sky = { top: oklch(0.18, 0.05, 252), bottom: oklch(0.278, 0.073, 248) };
const DAY: Sky = { top: oklch(0.4, 0.106, 248), bottom: oklch(0.557, 0.094, 231) };
const keyframes: ["sunrise" | "sunset", number, Sky][] = [
  ["sunrise", -70, NIGHT],
  ["sunrise", -25, { top: oklch(0.26, 0.07, 270), bottom: oklch(0.46, 0.1, 15) }], // indigo over a muted rose
  ["sunrise", 15, { top: oklch(0.35, 0.09, 255), bottom: oklch(0.54, 0.1, 55) }], // navy over apricot
  ["sunrise", 75, DAY],
  ["sunset", -100, DAY],
  ["sunset", -35, { top: oklch(0.38, 0.1, 252), bottom: oklch(0.55, 0.1, 65) }], // golden hour
  ["sunset", 5, { top: oklch(0.3, 0.085, 265), bottom: oklch(0.47, 0.12, 35) }], // afterglow
  ["sunset", 35, { top: oklch(0.21, 0.06, 258), bottom: oklch(0.31, 0.07, 270) }], // dusk
  ["sunset", 75, NIGHT],
];

/** The contrast white text keeps against any sky. */
export const MIN_CONTRAST = 4.5;

/** The sky gradient for now: interpolated between the day's keyframes, then greyed by cloud cover
 *  (percent) and darkened by precipitation. */
export function skyAt(now: Date, sunrise: Date, sunset: Date, cloudCover: number, precipitation: Precipitation): Sky {
  const t = now.getTime();
  const frames = keyframes.map(([anchor, minutes, sky]) => ({
    t: (anchor === "sunrise" ? sunrise : sunset).getTime() + minutes * MIN,
    sky,
  }));

  let sky = NIGHT;
  for (let i = 0; i < frames.length - 1; i++) {
    const a = frames[i];
    const b = frames[i + 1];
    if (t >= a.t && t < b.t) {
      const f = (t - a.t) / (b.t - a.t);
      sky = { top: mix(a.sky.top, b.sky.top, f), bottom: mix(a.sky.bottom, b.sky.bottom, f) };
      break;
    }
  }

  const cloud = Math.min(Math.max(cloudCover, 0), 100) / 100;
  const dim = { none: 1, snow: 0.95, drizzle: 0.9, rain: 0.85, storm: 0.75 }[precipitation];
  const weather = (c: Oklch): Oklch => ({ l: c.l * (1 - 0.1 * cloud) * dim, c: c.c * (1 - 0.75 * cloud), h: c.h });
  return { top: legible(weather(sky.top)), bottom: legible(weather(sky.bottom)) };
}

/** Darkens c just enough that white text on it keeps MIN_CONTRAST. */
function legible(c: Oklch): Oklch {
  let l = c.l;
  while (l > 0 && contrastWithWhite({ ...c, l }) < MIN_CONTRAST) l -= 0.005;
  return { ...c, l };
}

/** Interpolates in OKLab, so blue to orange passes through a pale tone rather than green. */
function mix(a: Oklch, b: Oklch, f: number): Oklch {
  const rad = Math.PI / 180;
  const [aa, ab] = [a.c * Math.cos(a.h * rad), a.c * Math.sin(a.h * rad)];
  const [ba, bb] = [b.c * Math.cos(b.h * rad), b.c * Math.sin(b.h * rad)];
  const x = aa + (ba - aa) * f;
  const y = ab + (bb - ab) * f;
  const h = (Math.atan2(y, x) / rad + 360) % 360;
  return { l: a.l + (b.l - a.l) * f, c: Math.hypot(x, y), h };
}

export const oklchCss = ({ l, c, h }: Oklch) => `oklch(${l.toFixed(3)} ${c.toFixed(3)} ${h.toFixed(1)})`;

/** Linear-light sRGB of c, each channel clamped to 0–1 (OKLab → sRGB, Björn Ottosson's matrices). */
function linearSrgb(c: Oklch): [number, number, number] {
  const rad = Math.PI / 180;
  const a = c.c * Math.cos(c.h * rad);
  const b = c.c * Math.sin(c.h * rad);
  const l = (c.l + 0.3963377774 * a + 0.2158037573 * b) ** 3;
  const m = (c.l - 0.1055613458 * a - 0.0638541728 * b) ** 3;
  const s = (c.l - 0.0894841775 * a - 1.291485548 * b) ** 3;
  const clamp = (v: number) => Math.min(1, Math.max(0, v));
  return [
    clamp(4.0767416621 * l - 3.3077115913 * m + 0.2309699292 * s),
    clamp(-1.2684380046 * l + 2.6097574011 * m - 0.3413193965 * s),
    clamp(-0.0041960863 * l - 0.7034186147 * m + 1.707614701 * s),
  ];
}

/** WCAG contrast ratio of white text on c (1–21). */
export function contrastWithWhite(c: Oklch): number {
  const [r, g, b] = linearSrgb(c);
  return 1.05 / (0.2126 * r + 0.7152 * g + 0.0722 * b + 0.05);
}

/** c as #rrggbb, for what takes no oklch(): the theme-color meta. */
export function oklchToHex(c: Oklch): string {
  return `#${linearSrgb(c)
    .map((v) => (v <= 0.0031308 ? 12.92 * v : 1.055 * v ** (1 / 2.4) - 0.055))
    .map((v) => Math.round(Math.min(1, Math.max(0, v)) * 255).toString(16).padStart(2, "0"))
    .join("")}`;
}
```

- [ ] **Step 8: Run the tests and the type check**

Run: `cd frontend && bun run test && bun run check`

Expected: PASS, with 0 errors and 0 warnings.

- [ ] **Step 9: Commit**

```bash
git add -A frontend/src/lib/helpers
git commit -m "Add Croatian formatting, weather labels and a navy-to-teal sky

Times follow the forecast location's zone; the sky keeps white text at 4.5:1.

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 10: Device and location helpers, and the city list

**Files:**
- Create: `frontend/src/lib/helpers/{storage,geo,location,wind,uv,air,svg}.ts`, each with its `*.test.ts`
- Create: `frontend/src/lib/data/cities.ts`

**Interfaces:**
- Consumes: the types from Task 8.
- Produces:
  - from `storage`: `load<T>(key, valid)` and `save(key, value)`;
  - from `geo`: `roundCoordinates(c)` and `cellKey(c)`;
  - from `location`: `MAX_RECENT`, `isSavedLocation(v)`, `isRecentList(v)`, `addRecent(recent, place)`, `locationLabel(l)` and `positionErrorMessage(code)`;
  - from `wind`: `compassPoint(deg)`, `compassName(deg)`, `windArrowDegrees(deg)` and `windStrength(kmh)`;
  - from `uv`: `uvBand(index) → { label, level }`;
  - from `air`: `aqiLabels`, `pollenNames`, `pollenLevelLabels` and `activePollen(aq)`;
  - from `svg`: `stripAnimation(svg)` and `stillSvgUrl(url)`;
  - from `data/cities`: `cities: SavedLocation[]`.

- [ ] **Step 1: Write the failing tests**

`frontend/src/lib/helpers/storage.test.ts`:

```ts
import { afterEach, describe, expect, it, vi } from "vitest";
import { load, save } from "./storage";

const isNumber = (v: unknown): v is number => typeof v === "number";

function fakeStorage(init: Record<string, string> = {}) {
  const data = new Map(Object.entries(init));
  return {
    getItem: (key: string) => data.get(key) ?? null,
    setItem: (key: string, value: string) => void data.set(key, value),
  };
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe("storage", () => {
  it("round-trips a valid value", () => {
    vi.stubGlobal("localStorage", fakeStorage());
    save("n", 42);
    expect(load("n", isNumber)).toBe(42);
  });

  it("treats a missing, corrupt or invalid value as nothing saved", () => {
    vi.stubGlobal("localStorage", fakeStorage({ corrupt: "{", wrong: '"text"' }));
    expect(load("missing", isNumber)).toBeUndefined();
    expect(load("corrupt", isNumber)).toBeUndefined();
    expect(load("wrong", isNumber)).toBeUndefined();
  });

  it("survives storage that throws or doesn't exist", () => {
    vi.stubGlobal("localStorage", {
      getItem: () => {
        throw new DOMException("blocked", "SecurityError");
      },
      setItem: () => {
        throw new DOMException("full", "QuotaExceededError");
      },
    });
    expect(load("n", isNumber)).toBeUndefined();
    expect(() => save("n", 1)).not.toThrow();

    vi.stubGlobal("localStorage", undefined);
    expect(load("n", isNumber)).toBeUndefined();
    expect(() => save("n", 1)).not.toThrow();
  });
});
```

`frontend/src/lib/helpers/geo.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import { cellKey, roundCoordinates } from "./geo";

describe("roundCoordinates", () => {
  it("rounds to the ~1 km cell, never to -0", () => {
    expect(roundCoordinates({ latitude: 45.5936, longitude: 17.2251 })).toEqual({ latitude: 45.59, longitude: 17.23 });
    const zero = roundCoordinates({ latitude: -0.001, longitude: 0.004 });
    expect(Object.is(zero.latitude, 0)).toBe(true);
  });
});

describe("cellKey", () => {
  it("is shared by neighbours in one cell", () => {
    expect(cellKey({ latitude: 45.5936, longitude: 17.2251 })).toBe("45.59,17.23");
    expect(cellKey({ latitude: 45.5912, longitude: 17.2263 })).toBe("45.59,17.23");
    expect(cellKey({ latitude: 45.6, longitude: 17 })).toBe("45.60,17.00");
  });
});
```

`frontend/src/lib/helpers/location.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import { cities } from "$lib/data/cities";
import type { SavedLocation } from "$lib/types/location";
import { addRecent, isSavedLocation, locationLabel, MAX_RECENT, positionErrorMessage } from "./location";

const place = (name: string, latitude: number, longitude: number): SavedLocation => ({
  kind: "place",
  name,
  county: null,
  latitude,
  longitude,
});

describe("isSavedLocation", () => {
  it("accepts what the app stores", () => {
    expect(isSavedLocation(place("Daruvar", 45.59, 17.22))).toBe(true);
    expect(isSavedLocation({ kind: "gps", name: null, county: null, latitude: 43.86, longitude: 18.41 })).toBe(true);
  });

  it("rejects anything else, so a corrupt or older value means nothing saved", () => {
    const bad = [
      null,
      "Daruvar",
      42,
      {},
      { ...place("x", 1, 2), kind: "city" },
      { ...place("x", 1, 2), latitude: "45" },
      place("x", 91, 0),
      place("x", 0, 181),
      { kind: "place", name: "x", latitude: 1, longitude: 2 },
    ];
    for (const v of bad) expect(isSavedLocation(v), JSON.stringify(v)).toBe(false);
  });
});

describe("addRecent", () => {
  it("puts the newest first, once per ~1 km cell, at most MAX_RECENT", () => {
    let recent: SavedLocation[] = [];
    for (let i = 0; i < 7; i++) recent = addRecent(recent, place(`P${i}`, 45 + i / 10, 16));
    expect(recent.map((r) => r.name)).toEqual(["P6", "P5", "P4", "P3", "P2"]);
    expect(recent).toHaveLength(MAX_RECENT);

    recent = addRecent(recent, place("P4 again", 45.4012, 16.0004));
    expect(recent.map((r) => r.name)).toEqual(["P4 again", "P6", "P5", "P3", "P2"]);
  });

  it("leaves GPS fixes out", () => {
    const recent = [place("Daruvar", 45.59, 17.22)];
    expect(addRecent(recent, { kind: "gps", name: "Split", county: null, latitude: 43.51, longitude: 16.44 })).toBe(recent);
  });
});

describe("locationLabel", () => {
  it("calls a position with no Croatian place nearby 'Moja lokacija'", () => {
    expect(locationLabel({ kind: "gps", name: null, county: null, latitude: 43.86, longitude: 18.41 })).toBe("Moja lokacija");
    expect(locationLabel(place("Daruvar", 45.59, 17.22))).toBe("Daruvar");
  });
});

describe("positionErrorMessage", () => {
  it("explains each failure in Croatian", () => {
    const messages = [0, 1, 2, 3].map(positionErrorMessage);
    expect(new Set(messages).size).toBe(4);
    expect(messages[1]).toContain("nije dopušten");
    expect(messages[3]).toContain("predugo");
  });
});

describe("cities", () => {
  it("offers 20 distinct, valid, rounded places", () => {
    expect(cities).toHaveLength(20);
    expect(new Set(cities.map((c) => c.name)).size).toBe(20);
    for (const c of cities) {
      expect(isSavedLocation(c), c.name ?? "").toBe(true);
      expect(Math.round(c.latitude * 100) / 100).toBe(c.latitude);
      expect(Math.round(c.longitude * 100) / 100).toBe(c.longitude);
    }
  });
});
```

`frontend/src/lib/helpers/wind.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import { compassName, compassPoint, windArrowDegrees, windStrength } from "./wind";

describe("compass", () => {
  it("names where the wind blows from, on 8 Croatian points", () => {
    expect([0, 22, 23, 90, 225, 359, -45, 405].map(compassPoint)).toEqual(["S", "S", "SI", "I", "JZ", "S", "SZ", "SI"]);
    expect(compassName(135)).toBe("jugoistok");
  });

  it("points the arrow where the wind blows to", () => {
    expect(windArrowDegrees(225)).toBe(45);
    expect(windArrowDegrees(0)).toBe(180);
  });
});

describe("windStrength", () => {
  it("groups Beaufort forces in Croatian", () => {
    expect([1, 12, 30, 50, 70, 120].map(windStrength)).toEqual([
      "Tiho",
      "Slab vjetar",
      "Umjeren vjetar",
      "Jak vjetar",
      "Olujni vjetar",
      "Orkanski vjetar",
    ]);
  });
});
```

`frontend/src/lib/helpers/uv.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import { uvBand } from "./uv";

describe("uvBand", () => {
  it("follows the WHO bands on the rounded index", () => {
    expect([0, 2.4, 2.55, 5, 6, 8, 11].map((i) => uvBand(i).label)).toEqual([
      "Nizak",
      "Nizak",
      "Umjeren",
      "Umjeren",
      "Visok",
      "Vrlo visok",
      "Ekstreman",
    ]);
    expect(uvBand(9).level).toBe(3);
  });
});
```

`frontend/src/lib/helpers/air.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import type { AirQuality, PollenLevel, PollenType } from "$lib/types/airquality";
import { activePollen, aqiLabels, pollenNames } from "./air";

const aq = (pollen: [PollenType, PollenLevel][]): AirQuality => ({
  aqi: { value: 30, level: "fair" },
  pollen: pollen.map(([type, level]) => ({ type, level, current: 0, todayMax: 0 })),
  fetchedAt: "",
});

describe("activePollen", () => {
  it("keeps the pollen in the air, the highest first", () => {
    const got = activePollen(aq([["grass", "low"], ["birch", "none"], ["ragweed", "very_high"], ["mugwort", "moderate"]]));
    expect(got.map((p) => p.type)).toEqual(["ragweed", "mugwort", "grass"]);
  });

  it("is empty outside the season", () => {
    expect(activePollen(aq([["grass", "none"], ["ragweed", "none"]]))).toEqual([]);
  });
});

describe("labels", () => {
  it("speak Croatian", () => {
    expect(aqiLabels.fair).toBe("Prihvatljiva");
    expect(pollenNames.ragweed).toBe("Ambrozija");
  });
});
```

`frontend/src/lib/helpers/svg.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import { stripAnimation } from "./svg";

describe("stripAnimation", () => {
  it("removes SMIL elements, self-closing or not, and keeps the drawing", () => {
    const svg =
      '<svg><g><path d="M0 0"/><animateTransform attributeName="transform" type="rotate" values="0;45" dur="6s" repeatCount="indefinite"/></g>' +
      '<circle r="2"><animate attributeName="r" values="2;3" dur="1s"></animate></circle><set attributeName="opacity" to="1"/></svg>';
    expect(stripAnimation(svg)).toBe('<svg><g><path d="M0 0"/></g><circle r="2"></circle></svg>');
  });
});
```

Run: `cd frontend && bun run test`

Expected: FAIL, because the modules don't exist.

- [ ] **Step 2: Implement the helpers**

`frontend/src/lib/helpers/storage.ts`:

```ts
// Device storage is a convenience. Private browsing, blocked site data or a full quota make it throw
// or vanish, and an older app version may have left another shape behind, so every read is
// validated and every failure simply means "nothing saved".

/** The value saved under key; undefined when there is none, it doesn't parse, or it fails valid. */
export function load<T>(key: string, valid: (value: unknown) => value is T): T | undefined {
  try {
    const raw = globalThis.localStorage?.getItem(key);
    if (raw == null) return undefined;
    const value: unknown = JSON.parse(raw);
    return valid(value) ? value : undefined;
  } catch {
    return undefined;
  }
}

/** Saves value under key; does nothing when storage is unavailable or full. */
export function save(key: string, value: unknown): void {
  try {
    globalThis.localStorage?.setItem(key, JSON.stringify(value));
  } catch {
    // Best effort: the app works without it, it just won't remember.
  }
}
```

`frontend/src/lib/helpers/geo.ts`:

```ts
import type { Coordinates } from "$lib/types/location";

const round2 = (v: number) => Math.round(v * 100) / 100 || 0; // || 0: never -0

/** Rounded to 2 decimals (~1 km), the backend's cache cell: the only precision stored on the
 *  device or sent to the API. */
export const roundCoordinates = ({ latitude, longitude }: Coordinates): Coordinates => ({
  latitude: round2(latitude),
  longitude: round2(longitude),
});

/** The ~1 km cell, e.g. "45.59,17.23": two locations in one cell get the same forecast. */
export function cellKey(c: Coordinates): string {
  const r = roundCoordinates(c);
  return `${r.latitude.toFixed(2)},${r.longitude.toFixed(2)}`;
}
```

`frontend/src/lib/helpers/location.ts`:

```ts
import type { SavedLocation } from "$lib/types/location";
import { cellKey } from "./geo";

export const MAX_RECENT = 5;

/** A SavedLocation as stored by this version of the app. */
export function isSavedLocation(v: unknown): v is SavedLocation {
  if (typeof v !== "object" || v === null) return false;
  const l = v as Record<string, unknown>;
  return (
    (l.kind === "gps" || l.kind === "place") &&
    (l.name === null || typeof l.name === "string") &&
    (l.county === null || typeof l.county === "string") &&
    typeof l.latitude === "number" &&
    Math.abs(l.latitude) <= 90 &&
    typeof l.longitude === "number" &&
    Math.abs(l.longitude) <= 180
  );
}

export const isRecentList = (v: unknown): v is SavedLocation[] => Array.isArray(v) && v.every(isSavedLocation);

/** The recents with place in front: newest first, one entry per ~1 km cell, at most MAX_RECENT.
 *  GPS fixes aren't recents, since "Koristi moju lokaciju" is always offered on its own. */
export function addRecent(recent: SavedLocation[], place: SavedLocation): SavedLocation[] {
  if (place.kind !== "place") return recent;
  const key = cellKey(place);
  return [place, ...recent.filter((r) => cellKey(r) !== key)].slice(0, MAX_RECENT);
}

/** What the header calls a location. */
export const locationLabel = (l: SavedLocation) => l.name ?? "Moja lokacija";

/** Croatian text for a failed position request: GeolocationPositionError's code, or 0 when the
 *  browser can't locate at all. */
export function positionErrorMessage(code: number): string {
  switch (code) {
    case 1:
      return "Pristup lokaciji nije dopušten. Dopustite ga u postavkama preglednika ili odaberite mjesto.";
    case 2:
      return "Lokacija trenutno nije dostupna. Pokušajte ponovno ili odaberite mjesto.";
    case 3:
      return "Određivanje lokacije traje predugo. Pokušajte ponovno ili odaberite mjesto.";
    default:
      return "Ovaj preglednik ne može odrediti lokaciju. Odaberite mjesto.";
  }
}
```

`frontend/src/lib/helpers/wind.ts`:

```ts
const points = ["S", "SI", "I", "JI", "J", "JZ", "Z", "SZ"];
const names = ["sjever", "sjeveroistok", "istok", "jugoistok", "jug", "jugozapad", "zapad", "sjeverozapad"];
const index = (degrees: number) => Math.round((((degrees % 360) + 360) % 360) / 45) % 8;

/** The 8-point compass point the wind blows from, in Croatian: 0° "S", 45° "SI", 225° "JZ". */
export const compassPoint = (degrees: number) => points[index(degrees)];

/** The same, spelled out: "jugozapad". */
export const compassName = (degrees: number) => names[index(degrees)];

/** Rotation for an arrow drawn pointing up, so it points where the wind blows to. */
export const windArrowDegrees = (degrees: number) => (degrees + 180) % 360;

/** The Beaufort groups in Croatian, by km/h. */
export function windStrength(kmh: number): string {
  if (kmh < 2) return "Tiho";
  if (kmh < 20) return "Slab vjetar";
  if (kmh < 39) return "Umjeren vjetar";
  if (kmh < 62) return "Jak vjetar";
  if (kmh < 89) return "Olujni vjetar";
  return "Orkanski vjetar";
}
```

`frontend/src/lib/helpers/uv.ts`:

```ts
/** The WHO band of a UV index, in Croatian, with a level from 0 (low) to 4 (extreme). */
export function uvBand(index: number): { label: string; level: 0 | 1 | 2 | 3 | 4 } {
  const uv = Math.round(index);
  if (uv < 3) return { label: "Nizak", level: 0 };
  if (uv < 6) return { label: "Umjeren", level: 1 };
  if (uv < 8) return { label: "Visok", level: 2 };
  if (uv < 11) return { label: "Vrlo visok", level: 3 };
  return { label: "Ekstreman", level: 4 };
}
```

`frontend/src/lib/helpers/air.ts`:

```ts
import type { AirQuality, AQILevel, PollenLevel, PollenType } from "$lib/types/airquality";

/** The European Air Quality Index bands, as the EEA names them in Croatian. */
export const aqiLabels: Record<AQILevel, string> = {
  good: "Dobra",
  fair: "Prihvatljiva",
  moderate: "Umjerena",
  poor: "Loša",
  very_poor: "Vrlo loša",
  extremely_poor: "Izrazito loša",
};

export const pollenNames: Record<PollenType, string> = {
  alder: "Joha",
  birch: "Breza",
  olive: "Maslina",
  grass: "Trave",
  mugwort: "Pelin",
  ragweed: "Ambrozija",
};

export const pollenLevelLabels: Record<PollenLevel, string> = {
  none: "nema",
  low: "nisko",
  moderate: "umjereno",
  high: "visoko",
  very_high: "vrlo visoko",
};

const severity: Record<PollenLevel, number> = { none: 0, low: 1, moderate: 2, high: 3, very_high: 4 };

/** Today's pollen in the air, the highest first; empty outside the season. */
export const activePollen = (aq: AirQuality) =>
  aq.pollen.filter((p) => p.level !== "none").sort((a, b) => severity[b.level] - severity[a.level]);
```

`frontend/src/lib/helpers/svg.ts`:

```ts
/** The SVG without its SMIL animation (<animate>, <animateTransform>, <animateMotion>, <set>). An
 *  <img> plays SMIL whatever CSS says, so readers who asked for reduced motion get this still. */
export function stripAnimation(svg: string): string {
  return svg
    .replace(/<(animate|animateTransform|animateMotion|set)\b[^>]*\/>/g, "")
    .replace(/<(animate|animateTransform|animateMotion|set)\b[^>]*>[\s\S]*?<\/\1>/g, "");
}

const stills = new Map<string, Promise<string>>();

/** A data: URL of the still version of the SVG at url, fetched once per url. */
export function stillSvgUrl(url: string): Promise<string> {
  let still = stills.get(url);
  if (!still) {
    still = fetch(url)
      .then((r) => r.text())
      .then((svg) => `data:image/svg+xml,${encodeURIComponent(stripAnimation(svg))}`);
    stills.set(url, still);
    still.catch(() => stills.delete(url)); // retry next time
  }
  return still;
}
```

`frontend/src/lib/data/cities.ts`:

```ts
import type { SavedLocation } from "$lib/types/location";

const city = (name: string, county: string, latitude: number, longitude: number): SavedLocation => ({
  kind: "place",
  name,
  county,
  latitude,
  longitude,
});

/** One city per county (Zagreb stands for both Grad Zagreb and Zagrebačka), largest first: the
 *  choices offered before any search. Coordinates are GeoNames', rounded like every location. */
export const cities: SavedLocation[] = [
  city("Zagreb", "Grad Zagreb", 45.81, 15.98),
  city("Split", "Splitsko-dalmatinska županija", 43.51, 16.44),
  city("Rijeka", "Primorsko-goranska županija", 45.33, 14.44),
  city("Osijek", "Osječko-baranjska županija", 45.55, 18.69),
  city("Zadar", "Zadarska županija", 44.12, 15.23),
  city("Pula", "Istarska županija", 44.87, 13.85),
  city("Slavonski Brod", "Brodsko-posavska županija", 45.16, 18.02),
  city("Karlovac", "Karlovačka županija", 45.49, 15.55),
  city("Varaždin", "Varaždinska županija", 46.3, 16.34),
  city("Šibenik", "Šibensko-kninska županija", 43.73, 15.89),
  city("Sisak", "Sisačko-moslavačka županija", 45.47, 16.38),
  city("Vinkovci", "Vukovarsko-srijemska županija", 45.29, 18.8),
  city("Dubrovnik", "Dubrovačko-neretvanska županija", 42.64, 18.11),
  city("Bjelovar", "Bjelovarsko-bilogorska županija", 45.9, 16.85),
  city("Koprivnica", "Koprivničko-križevačka županija", 46.16, 16.83),
  city("Požega", "Požeško-slavonska županija", 45.34, 17.69),
  city("Čakovec", "Međimurska županija", 46.38, 16.43),
  city("Virovitica", "Virovitičko-podravska županija", 45.83, 17.39),
  city("Gospić", "Ličko-senjska županija", 44.55, 15.37),
  city("Krapina", "Krapinsko-zagorska županija", 46.16, 15.87),
];
```

- [ ] **Step 3: Run the tests and the type check**

Run: `cd frontend && bun run test && bun run check`

Expected: PASS, with 0 errors and 0 warnings.

- [ ] **Step 4: Commit**

```bash
git add -A frontend/src/lib
git commit -m "Add storage, location, wind, UV, air and icon helpers and the 20 cities

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 11: Stores: clock, location and weather

**Files:**
- Create: `frontend/src/lib/helpers/forecast.ts`, `frontend/src/lib/helpers/forecast.test.ts`, `frontend/src/lib/stores/clock.svelte.ts`, `frontend/src/lib/stores/location.svelte.ts`, `frontend/src/lib/stores/weather.svelte.ts`

**Interfaces:**
- Consumes: Task 8's services, `geo`, `location` and `storage` (Task 10), and `isoDateOf` (Task 9).
- Produces:
  - from `helpers/forecast`: `CACHE_VERSION`, `CACHE_MAX_AGE_MS`, `type CachedForecast`, `isCachedForecast`, `usableCache(cache, cell, now)`, `upcomingHours(hourly, now, count)`, `todayOf(forecast, now)`, `hoursOfDay(forecast, date)` and `weekRange(days)`;
  - `clock.now` and `clock.subscribe()`;
  - `locationStore` with getters `current`, `recent`, `locating` and `error`, and the methods `choose(place)`, `locate(): Promise<boolean>` and `relocateIfGranted()`;
  - `weatherStore` with getters `forecast`, `airQuality`, `warnings`, `updatedAt`, `error`, `loading` and `stale`, and the methods `show(location)`, `refresh(): Promise<void>` and `subscribe()`.

- [ ] **Step 1: Write the failing forecast-helper tests**

`frontend/src/lib/helpers/forecast.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import type { DailyForecast, Forecast, HourlyForecast } from "$lib/types/forecast";
import {
  CACHE_MAX_AGE_MS,
  CACHE_VERSION,
  hoursOfDay,
  isCachedForecast,
  todayOf,
  upcomingHours,
  usableCache,
  weekRange,
  type CachedForecast,
} from "./forecast";

const hour = (time: string): HourlyForecast => ({
  time,
  temperature: 10,
  apparentTemperature: 9,
  precipitationProbability: 0,
  precipitation: 0,
  weatherCode: 0,
  windSpeed: 5,
  windDirection: 90,
  isDay: false,
});

const day = (date: string, temperatureMin: number, temperatureMax: number): DailyForecast => ({
  date,
  temperatureMin,
  temperatureMax,
  weatherCode: 0,
  precipitationSum: 0,
  precipitationProbabilityMax: 0,
  windSpeedMax: 10,
  windGustsMax: 20,
  windDirectionDominant: 90,
  uvIndexMax: 3,
  sunrise: `${date}T06:46:00+02:00`,
  sunset: `${date}T18:35:00+02:00`,
  daylightSeconds: 42_600,
});

const forecast = (): Forecast => ({
  current: {
    time: "2026-09-29T23:15:00+02:00",
    temperature: 12,
    apparentTemperature: 11,
    humidity: 70,
    dewPoint: 7,
    precipitation: 0,
    weatherCode: 0,
    cloudCover: 0,
    pressure: 1013,
    windSpeed: 5,
    windDirection: 90,
    windGusts: 9,
    visibility: 20_000,
    uvIndex: 0,
    isDay: false,
  },
  hourly: ["2026-09-29T22:00:00+02:00", "2026-09-29T23:00:00+02:00", "2026-09-30T00:00:00+02:00", "2026-09-30T01:00:00+02:00"].map(hour),
  daily: [day("2026-09-29", 9, 19), day("2026-09-30", 7, 21)],
  timezone: "Europe/Zagreb",
  elevation: 161,
  fetchedAt: "2026-09-29T23:10:00+02:00",
});

describe("upcomingHours", () => {
  it("starts at the hour containing now", () => {
    const got = upcomingHours(forecast().hourly, new Date("2026-09-29T23:20:00+02:00"), 2);
    expect(got.map((h) => h.time)).toEqual(["2026-09-29T23:00:00+02:00", "2026-09-30T00:00:00+02:00"]);
  });
});

describe("todayOf", () => {
  it("picks today on the location's calendar", () => {
    expect(todayOf(forecast(), new Date("2026-09-29T23:30:00+02:00")).date).toBe("2026-09-29");
    // 22:30 UTC is already the 30th in Zagreb, whatever the device's zone.
    expect(todayOf(forecast(), new Date("2026-09-29T22:30:00Z")).date).toBe("2026-09-30");
  });

  it("falls back to the first day", () => {
    expect(todayOf(forecast(), new Date("2026-10-05T12:00:00+02:00")).date).toBe("2026-09-29");
  });
});

describe("hoursOfDay", () => {
  it("keeps one calendar day's hours", () => {
    expect(hoursOfDay(forecast(), "2026-09-30").map((h) => h.time)).toEqual(["2026-09-30T00:00:00+02:00", "2026-09-30T01:00:00+02:00"]);
  });
});

describe("weekRange", () => {
  it("spans the coldest minimum to the warmest maximum", () => {
    expect(weekRange(forecast().daily)).toEqual([7, 21]);
    expect(weekRange([])).toEqual([0, 1]);
  });
});

describe("the cached forecast", () => {
  const cache = (savedAt: string, cell = "45.59,17.23"): CachedForecast => ({ version: CACHE_VERSION, cell, savedAt, forecast: forecast() });
  const now = new Date("2026-09-29T23:30:00+02:00");

  it("is painted only for the same cell while it is recent", () => {
    const fresh = cache("2026-09-29T22:30:00+02:00");
    expect(usableCache(fresh, "45.59,17.23", now)).toBe(fresh);
    expect(usableCache(fresh, "43.51,16.44", now)).toBeUndefined();
    expect(usableCache(cache(new Date(now.getTime() - CACHE_MAX_AGE_MS - 1).toISOString()), "45.59,17.23", now)).toBeUndefined();
    expect(usableCache(cache("2026-09-30T05:00:00+02:00"), "45.59,17.23", now)).toBeUndefined(); // from the future
    expect(usableCache(undefined, "45.59,17.23", now)).toBeUndefined();
  });

  it("is recognized only in this version's shape", () => {
    expect(isCachedForecast(cache("2026-09-29T22:30:00+02:00"))).toBe(true);
    expect(isCachedForecast({ ...cache("2026-09-29T22:30:00+02:00"), version: 0 })).toBe(false);
    expect(isCachedForecast({ ...cache("not a date") })).toBe(false);
    expect(isCachedForecast({ ...cache("2026-09-29T22:30:00+02:00"), forecast: { current: {} } })).toBe(false);
    expect(isCachedForecast(null)).toBe(false);
  });
});
```

Run: `cd frontend && bun run test`

Expected: FAIL, because `./forecast` does not exist.

- [ ] **Step 2: Implement the forecast helpers**

`frontend/src/lib/helpers/forecast.ts`:

```ts
import type { DailyForecast, Forecast, HourlyForecast } from "$lib/types/forecast";
import { isoDateOf } from "./format";

/** Bumped whenever Forecast changes shape, so an older cached copy is ignored. */
export const CACHE_VERSION = 1;

/** How long a cached forecast is still worth painting while the fresh one loads. */
export const CACHE_MAX_AGE_MS = 12 * 3_600_000;

/** The last forecast shown, kept on the device so the next visit paints at once. */
export interface CachedForecast {
  version: number;
  /** cellKey of the location it is for */
  cell: string;
  savedAt: string;
  forecast: Forecast;
}

export function isCachedForecast(v: unknown): v is CachedForecast {
  if (typeof v !== "object" || v === null) return false;
  const c = v as Record<string, unknown>;
  const f = c.forecast as Record<string, unknown> | null | undefined;
  return (
    c.version === CACHE_VERSION &&
    typeof c.cell === "string" &&
    typeof c.savedAt === "string" &&
    !Number.isNaN(Date.parse(c.savedAt)) &&
    typeof f === "object" &&
    f !== null &&
    typeof f.current === "object" &&
    f.current !== null &&
    Array.isArray(f.hourly) &&
    Array.isArray(f.daily) &&
    typeof f.timezone === "string"
  );
}

/** The cached forecast when it is for cell and recent enough to paint. */
export function usableCache(cache: CachedForecast | undefined, cell: string, now: Date): CachedForecast | undefined {
  if (!cache || cache.cell !== cell) return undefined;
  const age = now.getTime() - Date.parse(cache.savedAt);
  return age >= 0 && age < CACHE_MAX_AGE_MS ? cache : undefined;
}

/** At most count hours, from the one containing now. */
export const upcomingHours = (hourly: HourlyForecast[], now: Date, count: number) =>
  hourly.filter((h) => Date.parse(h.time) + 3_600_000 > now.getTime()).slice(0, count);

/** Today on the location's calendar, or the first day the forecast has. */
export const todayOf = (f: Forecast, now: Date): DailyForecast =>
  f.daily.find((d) => d.date === isoDateOf(now, f.timezone)) ?? f.daily[0];

/** The hours of one day (YYYY-MM-DD on the location's calendar). */
export const hoursOfDay = (f: Forecast, date: string) =>
  f.hourly.filter((h) => isoDateOf(new Date(h.time), f.timezone) === date);

/** The coldest minimum and warmest maximum: the shared scale of the week's range bars. */
export function weekRange(days: DailyForecast[]): [number, number] {
  if (days.length === 0) return [0, 1];
  return [Math.min(...days.map((d) => d.temperatureMin)), Math.max(...days.map((d) => d.temperatureMax))];
}
```

Run: `cd frontend && bun run test`

Expected: PASS.

- [ ] **Step 3: Write the stores**

These touch `document`, `navigator` and `localStorage`. They are checked by `bun run check` here and in the browser in Task 12; their decisions live in the tested helpers.

`frontend/src/lib/stores/clock.svelte.ts`:

```ts
import { untrack } from "svelte";

/** The current time, updated at each whole minute while a component is subscribed: the page shows
 *  minutes ("prije 3 min", the hour strip, the sky), never seconds. */
function createClock() {
  let now = $state(new Date());
  let timer: ReturnType<typeof setTimeout> | null = null;
  let consumers = 0;

  function tick() {
    if (timer !== null) clearTimeout(timer);
    now = new Date();
    timer = setTimeout(tick, 60_000 - (now.getTime() % 60_000));
  }

  // A phone suspends timers in the background: catch up as soon as the page is shown again.
  function onVisibility() {
    if (document.visibilityState === "visible") tick();
  }

  return {
    get now() {
      return now;
    },
    /** `$effect(() => clock.subscribe())` ties the timer to the component's lifetime. */
    subscribe(): () => void {
      return untrack(() => {
        if (consumers++ === 0) {
          document.addEventListener("visibilitychange", onVisibility);
          tick();
        }
        return () => {
          if (--consumers > 0) return;
          document.removeEventListener("visibilitychange", onVisibility);
          if (timer !== null) clearTimeout(timer);
          timer = null;
        };
      });
    },
  };
}

export const clock = createClock();
```

`frontend/src/lib/stores/location.svelte.ts`:

```ts
import { untrack } from "svelte";
import { ApiError } from "$lib/api/client";
import { cellKey, roundCoordinates } from "$lib/helpers/geo";
import { addRecent, isRecentList, isSavedLocation, positionErrorMessage } from "$lib/helpers/location";
import { load, save } from "$lib/helpers/storage";
import { nearestPlace } from "$lib/services/places";
import type { SavedLocation } from "$lib/types/location";

const CURRENT_KEY = "vrijeme.location";
const RECENT_KEY = "vrijeme.recent";

/** How old a position the device already has may be before it is asked again. */
const POSITION_MAX_AGE_MS = 15 * 60_000;

function currentPosition(): Promise<GeolocationPosition> {
  return new Promise((resolve, reject) =>
    navigator.geolocation.getCurrentPosition(resolve, reject, { timeout: 10_000, maximumAge: POSITION_MAX_AGE_MS }),
  );
}

function createLocationStore() {
  let current = $state<SavedLocation | null>(load(CURRENT_KEY, isSavedLocation) ?? null);
  let recent = $state<SavedLocation[]>(load(RECENT_KEY, isRecentList) ?? []);
  let locating = $state(false);
  let error = $state<string | null>(null);

  function set(location: SavedLocation) {
    current = location;
    save(CURRENT_KEY, location);
  }

  /** Asks the device where it is, names the spot after the nearest Croatian settlement, and shows
   *  weather there. Resolves false, with error set, when there is no position to be had. */
  async function locate(): Promise<boolean> {
    if (!("geolocation" in navigator)) {
      error = positionErrorMessage(0);
      return false;
    }
    locating = true;
    error = null;
    try {
      const { coords } = await currentPosition();
      const at = roundCoordinates({ latitude: coords.latitude, longitude: coords.longitude });
      const previous = untrack(() => current);
      let name: string | null = null;
      let county: string | null = null;
      try {
        const place = await nearestPlace(at);
        name = place.name;
        county = place.county || null;
      } catch (e) {
        // A 404 means abroad or at sea: "Moja lokacija". Anything else (offline) keeps the name
        // this cell already had.
        const abroad = e instanceof ApiError && e.status === 404;
        if (!abroad && previous && cellKey(previous) === cellKey(at)) {
          name = previous.name;
          county = previous.county;
        }
      }
      set({ kind: "gps", name, county, ...at });
      return true;
    } catch (e) {
      error = positionErrorMessage(e instanceof GeolocationPositionError ? e.code : 0);
      return false;
    } finally {
      locating = false;
    }
  }

  return {
    /** What the app shows weather for; null until the reader picks something (the welcome screen). */
    get current() {
      return current;
    },
    /** Places picked before, newest first. */
    get recent() {
      return recent;
    },
    /** A position request is running. */
    get locating() {
      return locating;
    },
    /** Why the last position request failed, in Croatian; cleared by the next attempt or pick. */
    get error() {
      return error;
    },

    /** Shows weather for a picked place and remembers it among the recents. */
    choose(place: SavedLocation) {
      const location = { ...place, ...roundCoordinates(place) };
      error = null;
      set(location);
      recent = addRecent(untrack(() => recent), location);
      save(RECENT_KEY, recent);
    },

    locate,

    /** In GPS mode, follows the device on open without a prompt, but only when permission is
     *  already granted: a returning reader never meets a cold permission dialog. */
    async relocateIfGranted(): Promise<void> {
      if (untrack(() => current)?.kind !== "gps") return;
      try {
        const status = await navigator.permissions?.query({ name: "geolocation" });
        if (status?.state === "granted") await locate();
      } catch {
        // No Permissions API: stay on the saved position until the reader refreshes.
      }
    },
  };
}

export const locationStore = createLocationStore();
```

`frontend/src/lib/stores/weather.svelte.ts`:

```ts
import { untrack } from "svelte";
import { isAbortError } from "$lib/api/client";
import { CACHE_VERSION, isCachedForecast, usableCache, type CachedForecast } from "$lib/helpers/forecast";
import { cellKey } from "$lib/helpers/geo";
import { load, save } from "$lib/helpers/storage";
import { fetchAirQuality } from "$lib/services/airquality";
import { fetchForecast } from "$lib/services/forecast";
import { fetchWarnings } from "$lib/services/warnings";
import type { AirQuality } from "$lib/types/airquality";
import type { Forecast } from "$lib/types/forecast";
import type { Coordinates } from "$lib/types/location";
import type { Warning } from "$lib/types/warning";

const CACHE_KEY = "vrijeme.forecast";
/** The backend caches the forecast for 10 minutes; asking every 15 keeps an open page current. */
const POLL_MS = 15 * 60_000;
/** Coming back to the page after this long refreshes at once instead of at the next poll. */
const STALE_MS = 10 * 60_000;

function createWeatherStore() {
  let at = $state<Coordinates | null>(null);
  let forecast = $state<Forecast>();
  let airQuality = $state<AirQuality>();
  let warnings = $state<Warning[]>();
  let updatedAt = $state<Date | null>(null);
  let error = $state<unknown>(null);
  let loading = $state(false);

  let controller: AbortController | null = null;
  let timer: ReturnType<typeof setTimeout> | null = null;
  let consumers = 0;

  const age = () => {
    const since = untrack(() => updatedAt);
    return since ? Date.now() - since.getTime() : Infinity;
  };

  /** The next poll, POLL_MS after the last update; none while hidden or unwatched. */
  function schedule() {
    if (timer !== null) clearTimeout(timer);
    timer = null;
    if (consumers === 0 || document.visibilityState === "hidden") return;
    timer = setTimeout(() => void refresh(), Math.max(0, POLL_MS - age()));
  }

  /** Fetches forecast, air quality and warnings for the current cell in parallel. Each lands on its
   *  own, so a Meteoalarm outage never blanks the forecast. A newer refresh, a location change or
   *  the last consumer leaving aborts this one, and its results are dropped. */
  async function refresh(): Promise<void> {
    const target = untrack(() => at);
    if (!target) return;
    controller?.abort();
    const mine = (controller = new AbortController());
    loading = true;

    const [f, a, w] = await Promise.allSettled([
      fetchForecast(target, mine.signal),
      fetchAirQuality(target, mine.signal),
      fetchWarnings(target, mine.signal),
    ]);
    if (controller !== mine) return; // superseded: a newer refresh owns the state now
    controller = null;
    loading = false;

    if (f.status === "fulfilled") {
      forecast = f.value;
      updatedAt = new Date();
      error = null;
      save(CACHE_KEY, {
        version: CACHE_VERSION,
        cell: cellKey(target),
        savedAt: updatedAt.toISOString(),
        forecast: f.value,
      } satisfies CachedForecast);
    } else if (!isAbortError(f.reason)) {
      error = f.reason;
    }
    if (a.status === "fulfilled") airQuality = a.value;
    if (w.status === "fulfilled") warnings = w.value;
    schedule();
  }

  /** Switches to location's ~1 km cell: paints the forecast saved on the device for it, when recent,
   *  then refreshes. The same cell again changes nothing. */
  function show(location: Coordinates) {
    untrack(() => {
      const cell = cellKey(location);
      if (at && cellKey(at) === cell) return;
      at = { latitude: location.latitude, longitude: location.longitude };
      const cached = usableCache(load(CACHE_KEY, isCachedForecast), cell, new Date());
      forecast = cached?.forecast;
      updatedAt = cached ? new Date(cached.savedAt) : null;
      airQuality = undefined;
      warnings = undefined;
      error = null;
      void refresh();
    });
  }

  function onVisibility() {
    if (document.visibilityState === "visible" && age() > STALE_MS) void refresh();
    else schedule();
  }

  function onOnline() {
    void refresh();
  }

  return {
    get forecast() {
      return forecast;
    },
    get airQuality() {
      return airQuality;
    },
    get warnings() {
      return warnings;
    },
    /** When the forecast shown was fetched; a cached one keeps its own time. */
    get updatedAt() {
      return updatedAt;
    },
    /** The last forecast refresh's failure, cleared by the next success. */
    get error() {
      return error;
    },
    get loading() {
      return loading;
    },
    /** There is a forecast, but the last refresh failed: what's shown may be out of date. */
    get stale() {
      return forecast !== undefined && error !== null;
    },
    show,
    refresh,
    /** `$effect(() => weatherStore.subscribe())` polls while the page is mounted and visible. */
    subscribe(): () => void {
      return untrack(() => {
        if (consumers++ === 0) {
          document.addEventListener("visibilitychange", onVisibility);
          window.addEventListener("online", onOnline);
          schedule();
        }
        return () => {
          if (--consumers > 0) return;
          document.removeEventListener("visibilitychange", onVisibility);
          window.removeEventListener("online", onOnline);
          if (timer !== null) clearTimeout(timer);
          timer = null;
          controller?.abort();
          controller = null;
          loading = false;
        };
      });
    },
  };
}

export const weatherStore = createWeatherStore();
```

- [ ] **Step 4: Check and test**

Run: `cd frontend && bun run check && bun run test`

Expected: 0 errors, 0 warnings, all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add -A frontend/src/lib
git commit -m "Add the clock, location and weather stores

The location and the last forecast are kept on the device; the forecast,
air quality and warnings load in parallel and fail independently.

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 12: Location UI: sky, picker, sheet, welcome screen and header

**Files:**
- Create, ported from dashboard: `frontend/src/lib/components/sky-backdrop.svelte`, `frontend/src/lib/components/precipitation-layer.svelte`, `frontend/src/lib/components/weather-icon.svelte`
- Create: `frontend/src/lib/components/{location-picker,location-sheet,welcome,app-header}.svelte`
- Modify: `frontend/src/routes/+page.svelte`

**Interfaces:**
- Consumes: the stores (Task 11), the helpers (Tasks 9 and 10) and `cities`.
- Produces these components:
  - `<WeatherIcon name size? class?>`;
  - `<SkyBackdrop sky precipitation transitionMs?>`;
  - `<LocationPicker onchosen?>`;
  - `<LocationSheet bind:open>`;
  - `<Welcome>`;
  - `<AppHeader location updatedAt loading now onlocation onrefresh>`.

- [ ] **Step 1: Port the sky and the icons**

```bash
mkdir -p frontend/src/lib/components
for f in sky-backdrop precipitation-layer weather-icon; do
  cp ~/dev/go/husak-dashboard/frontend/src/lib/components/$f.svelte frontend/src/lib/components/
done
```

In `frontend/src/lib/components/weather-icon.svelte`, keep the `<script lang="ts" module>` block as it is. Replace everything after it with:

```svelte
<script lang="ts">
  import { prefersReducedMotion } from "svelte/motion";
  import { stillSvgUrl } from "$lib/helpers/svg";

  type Props = { name: string; size?: number; class?: string };
  let { name, size = 64, class: className }: Props = $props();

  const animated = $derived(icons[name] ?? icons["not-available"]);

  // Readers who asked their system for reduced motion get a still icon: an <img> plays the
  // icons' SMIL animation whatever CSS says. Until the still is ready the animated one shows.
  let still = $state<{ of: string; url: string }>();
  $effect(() => {
    if (!prefersReducedMotion.current) return;
    const of = animated;
    let live = true;
    stillSvgUrl(of)
      .then((url) => {
        if (live) still = { of, url };
      })
      .catch(() => {});
    return () => {
      live = false;
    };
  });
  const src = $derived(prefersReducedMotion.current && still?.of === animated ? still.url : animated);
</script>

<!-- Decorative: the text beside every icon says the same thing. -->
<img {src} width={size} height={size} alt="" draggable="false" class={className} />
```

- [ ] **Step 2: Write the location picker**

`frontend/src/lib/components/location-picker.svelte`:

```svelte
<script lang="ts">
  import History from "@lucide/svelte/icons/history";
  import LoaderCircle from "@lucide/svelte/icons/loader-circle";
  import LocateFixed from "@lucide/svelte/icons/locate-fixed";
  import Search from "@lucide/svelte/icons/search";
  import { isAbortError } from "$lib/api/client";
  import { cities } from "$lib/data/cities";
  import { searchPlaces } from "$lib/services/places";
  import { locationStore } from "$lib/stores/location.svelte";
  import type { SavedLocation } from "$lib/types/location";
  import type { Place } from "$lib/types/place";

  type Props = { onchosen?: () => void };
  let { onchosen }: Props = $props();
  const id = $props.id();

  let query = $state("");
  let results = $state<Place[]>([]);
  /** The query the results (or the failure) answer: shown only while it matches the input. */
  let answered = $state("");
  let failed = $state(false);
  let searching = $state(false);
  let input = $state<HTMLInputElement>();
  let list = $state<HTMLUListElement>();

  let timer: ReturnType<typeof setTimeout> | undefined;
  let controller: AbortController | undefined;

  const trimmed = $derived(query.trim());
  const current = $derived(answered !== "" && answered === trimmed);
  const status = $derived(
    !current ? "" : failed ? "Pretraživanje ne radi" : results.length === 1 ? "1 rezultat" : `${results.length} rezultata`,
  );

  /** Searches 250 ms after the last keystroke, cancelling the request before it. */
  function search() {
    clearTimeout(timer);
    controller?.abort();
    const q = query.trim();
    if (q.length < 2) {
      results = [];
      answered = "";
      searching = false;
      return;
    }
    searching = true;
    timer = setTimeout(async () => {
      const mine = (controller = new AbortController());
      try {
        results = await searchPlaces(q, mine.signal);
        failed = false;
      } catch (e) {
        if (isAbortError(e)) return;
        results = [];
        failed = true;
      }
      answered = q;
      searching = false;
    }, 250);
  }

  // A pending search dies with the picker (when the sheet closes).
  $effect(() => () => {
    clearTimeout(timer);
    controller?.abort();
  });

  function pick(location: SavedLocation) {
    locationStore.choose(location);
    query = "";
    results = [];
    answered = "";
    onchosen?.();
  }

  const fromSearch = (p: Place): SavedLocation => ({
    kind: "place",
    name: p.name,
    county: p.county || null,
    latitude: p.latitude,
    longitude: p.longitude,
  });

  async function useMyLocation() {
    if (await locationStore.locate()) onchosen?.();
  }

  /** Enter takes the best match; the down arrow moves into the results. */
  function onInputKeydown(e: KeyboardEvent) {
    if (e.key === "Enter" && current && results.length) {
      e.preventDefault();
      pick(fromSearch(results[0]));
    } else if (e.key === "ArrowDown") {
      e.preventDefault();
      list?.querySelector("button")?.focus();
    }
  }

  /** The arrows move between results; up from the first goes back to the input. */
  function onResultKeydown(e: KeyboardEvent) {
    const buttons = [...(list?.querySelectorAll("button") ?? [])];
    const i = buttons.indexOf(e.currentTarget as HTMLButtonElement);
    if (e.key === "ArrowDown") {
      e.preventDefault();
      buttons[Math.min(i + 1, buttons.length - 1)]?.focus();
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      (i <= 0 ? input : buttons[i - 1])?.focus();
    }
  }
</script>

<div class="flex flex-col gap-6">
  <div>
    <button
      type="button"
      onclick={useMyLocation}
      disabled={locationStore.locating}
      class="text-navy flex w-full items-center justify-center gap-2 rounded-2xl bg-white px-5 py-3.5 text-base font-semibold shadow-lg transition hover:bg-white/90 disabled:opacity-70"
    >
      {#if locationStore.locating}
        <LoaderCircle class="size-5 animate-spin" aria-hidden="true" /> Tražim vašu lokaciju…
      {:else}
        <LocateFixed class="size-5" aria-hidden="true" /> Koristi moju lokaciju
      {/if}
    </button>
    {#if locationStore.error}
      <p role="alert" class="text-warn-yellow mt-2 text-sm">{locationStore.error}</p>
    {/if}
  </div>

  <div>
    <label for="search-{id}" class="sr-only">Pretraži mjesto</label>
    <div class="relative">
      <Search class="text-muted-foreground pointer-events-none absolute top-1/2 left-3.5 size-5 -translate-y-1/2" aria-hidden="true" />
      <input
        bind:this={input}
        bind:value={query}
        oninput={search}
        onkeydown={onInputKeydown}
        id="search-{id}"
        type="search"
        placeholder="Pretraži mjesto…"
        autocomplete="off"
        spellcheck="false"
        enterkeyhint="search"
        aria-describedby="search-status-{id}"
        class="placeholder:text-muted-foreground w-full rounded-2xl bg-white/10 py-3 pr-11 pl-11 text-base focus:bg-white/15"
      />
      {#if searching}
        <LoaderCircle class="text-muted-foreground pointer-events-none absolute top-1/2 right-3.5 size-5 -translate-y-1/2 animate-spin" aria-hidden="true" />
      {/if}
    </div>
    <p id="search-status-{id}" class="sr-only" aria-live="polite">{status}</p>
    {#if current && results.length}
      <ul bind:this={list} aria-label="Rezultati pretrage" class="mt-2 flex flex-col gap-1">
        {#each results as place (`${place.name}|${place.latitude}|${place.longitude}`)}
          <li>
            <button
              type="button"
              onclick={() => pick(fromSearch(place))}
              onkeydown={onResultKeydown}
              class="flex w-full flex-col rounded-xl px-3 py-2 text-left hover:bg-white/10 focus:bg-white/15"
            >
              <span class="font-medium">{place.name}</span>
              {#if place.county}<span class="text-muted-foreground text-sm">{place.county}</span>{/if}
            </button>
          </li>
        {/each}
      </ul>
    {:else if current && !searching}
      <p class="text-muted-foreground mt-2 px-1 text-sm">
        {failed ? "Pretraživanje trenutno ne radi. Odaberite grad s popisa." : `Nema mjesta „${answered}”.`}
      </p>
    {/if}
  </div>

  {#if locationStore.recent.length}
    <section aria-labelledby="recent-{id}">
      <h2 id="recent-{id}" class="text-muted-foreground mb-2 text-xs font-medium tracking-wider uppercase">Nedavno</h2>
      <ul class="flex flex-wrap gap-2">
        {#each locationStore.recent as place (`${place.latitude},${place.longitude}`)}
          <li>
            <button type="button" onclick={() => pick(place)} class="flex items-center gap-1.5 rounded-full bg-white/10 px-3.5 py-1.5 text-sm hover:bg-white/20">
              <History class="size-4 opacity-70" aria-hidden="true" />{place.name}
            </button>
          </li>
        {/each}
      </ul>
    </section>
  {/if}

  <section aria-labelledby="cities-{id}">
    <h2 id="cities-{id}" class="text-muted-foreground mb-2 text-xs font-medium tracking-wider uppercase">Gradovi</h2>
    <ul class="grid grid-cols-2 gap-2 sm:grid-cols-3">
      {#each cities as city (city.name)}
        <li>
          <button type="button" onclick={() => pick(city)} class="w-full rounded-xl bg-white/8 px-3 py-2.5 text-left hover:bg-white/15">
            {city.name}
          </button>
        </li>
      {/each}
    </ul>
  </section>
</div>
```

- [ ] **Step 3: Write the sheet, the welcome screen and the header**

`frontend/src/lib/components/location-sheet.svelte`:

```svelte
<script lang="ts">
  import X from "@lucide/svelte/icons/x";
  import LocationPicker from "./location-picker.svelte";

  type Props = { open: boolean };
  let { open = $bindable() }: Props = $props();
  let dialog = $state<HTMLDialogElement>();

  // The <dialog> owns modality: the focus trap, Escape, and the inert page behind it. open mirrors it.
  $effect(() => {
    if (!dialog) return;
    if (open && !dialog.open) dialog.showModal();
    else if (!open && dialog.open) dialog.close();
  });
</script>

<!-- A click on the backdrop lands on the dialog itself and closes it; Escape is the keyboard way. -->
<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_noninteractive_element_interactions -->
<dialog
  bind:this={dialog}
  onclose={() => (open = false)}
  onclick={(e) => {
    if (e.target === dialog) dialog?.close();
  }}
  aria-labelledby="sheet-title"
  class="bg-sheet text-foreground border-0 p-0 shadow-2xl"
>
  <div class="flex flex-col gap-5 p-5 sm:p-6">
    <div class="flex items-center justify-between gap-4">
      <h2 id="sheet-title" class="text-lg font-semibold">Odaberite mjesto</h2>
      <button type="button" aria-label="Zatvori" onclick={() => dialog?.close()} class="-mr-2 rounded-full p-2 hover:bg-white/10">
        <X class="size-5" aria-hidden="true" />
      </button>
    </div>
    {#if open}<LocationPicker onchosen={() => dialog?.close()} />{/if}
  </div>
</dialog>

<style>
  /* A bottom sheet on phones… */
  dialog {
    margin: auto 0 0;
    width: 100%;
    max-width: 100%;
    max-height: 88dvh;
    overflow-y: auto;
    border-radius: 1.5rem 1.5rem 0 0;
    padding-bottom: env(safe-area-inset-bottom);
  }
  /* …and a centred card from 640px. */
  @media (min-width: 640px) {
    dialog {
      margin: auto;
      max-width: 32rem;
      border-radius: 1.5rem;
    }
  }
  dialog::backdrop {
    background: oklch(0 0 0 / 55%);
  }
</style>
```

`frontend/src/lib/components/welcome.svelte`:

```svelte
<script lang="ts">
  import LocationPicker from "./location-picker.svelte";
</script>

<main class="mx-auto flex min-h-dvh w-full max-w-xl flex-col justify-center gap-8 px-4 py-10">
  <header class="text-lift text-center">
    <h1 class="text-5xl font-extralight tracking-tight sm:text-6xl">Vrijeme</h1>
    <p class="text-muted-foreground mt-3 text-lg">Prognoza za mjesto u kojem jeste: sada, po satima i za idućih 7 dana.</p>
  </header>
  <section aria-label="Odabir mjesta" class="panel p-5 sm:p-6">
    <LocationPicker />
  </section>
  <p class="text-muted-foreground text-center text-sm">
    Lokaciju pamti samo vaš uređaj, a prognozu tražimo za točku zaokruženu na kilometar.
  </p>
</main>
```

`frontend/src/lib/components/app-header.svelte`:

```svelte
<script lang="ts">
  import ChevronDown from "@lucide/svelte/icons/chevron-down";
  import Navigation from "@lucide/svelte/icons/navigation";
  import RotateCw from "@lucide/svelte/icons/rotate-cw";
  import { formatAgo } from "$lib/helpers/format";
  import { locationLabel } from "$lib/helpers/location";
  import type { SavedLocation } from "$lib/types/location";

  type Props = {
    location: SavedLocation;
    updatedAt: Date | null;
    loading: boolean;
    now: Date;
    onlocation: () => void;
    onrefresh: () => void;
  };
  let { location, updatedAt, loading, now, onlocation, onrefresh }: Props = $props();
</script>

<header class="flex items-start justify-between gap-3">
  <button
    type="button"
    onclick={onlocation}
    aria-label="Promijeni mjesto, sada {locationLabel(location)}"
    class="group -ml-2 min-w-0 rounded-2xl px-2 py-1 text-left hover:bg-white/10"
  >
    <span class="text-lift flex items-center gap-1.5 text-2xl font-semibold">
      {#if location.kind === "gps"}<Navigation class="size-4 shrink-0 fill-current" aria-hidden="true" />{/if}
      <span class="truncate">{locationLabel(location)}</span>
      <ChevronDown class="size-5 shrink-0 opacity-70 transition group-hover:translate-y-0.5" aria-hidden="true" />
    </span>
    {#if location.county}<span class="text-muted-foreground block truncate text-sm">{location.county}</span>{/if}
  </button>
  <div class="text-muted-foreground flex shrink-0 items-center gap-1 text-sm">
    {#if updatedAt}<span>Ažurirano {formatAgo(updatedAt, now)}</span>{/if}
    <button type="button" onclick={onrefresh} disabled={loading} aria-label="Osvježi prognozu" class="rounded-full p-2 hover:bg-white/10 disabled:opacity-60">
      <RotateCw class="size-5 {loading ? 'animate-spin' : ''}" aria-hidden="true" />
    </button>
  </div>
</header>
```

- [ ] **Step 4: Put the page together (the forecast panels arrive in Task 13)**

`frontend/src/routes/+page.svelte`:

```svelte
<script lang="ts">
  import { onMount } from "svelte";
  import { prefersReducedMotion } from "svelte/motion";
  import AppHeader from "$lib/components/app-header.svelte";
  import LocationSheet from "$lib/components/location-sheet.svelte";
  import SkyBackdrop from "$lib/components/sky-backdrop.svelte";
  import Welcome from "$lib/components/welcome.svelte";
  import { todayOf } from "$lib/helpers/forecast";
  import { formatDegrees } from "$lib/helpers/format";
  import { locationLabel } from "$lib/helpers/location";
  import { skyAt } from "$lib/helpers/sky";
  import { weatherInfo } from "$lib/helpers/weather";
  import { clock } from "$lib/stores/clock.svelte";
  import { locationStore } from "$lib/stores/location.svelte";
  import { weatherStore } from "$lib/stores/weather.svelte";

  $effect(() => clock.subscribe());
  $effect(() => weatherStore.subscribe());
  // Follow the chosen location; the store ignores a repeat of the same ~1 km cell.
  $effect(() => {
    const location = locationStore.current;
    if (location) weatherStore.show(location);
  });
  // A returning GPS reader whose permission stands is followed without a prompt.
  onMount(() => void locationStore.relocateIfGranted());

  let sheetOpen = $state(false);

  const now = $derived(clock.now);
  const forecast = $derived(weatherStore.forecast);
  const today = $derived(forecast ? todayOf(forecast, now) : undefined);
  const sunrise = $derived(today ? new Date(today.sunrise) : undefined);
  const sunset = $derived(today ? new Date(today.sunset) : undefined);
  const isDay = $derived(sunrise && sunset ? now >= sunrise && now < sunset : true);
  const info = $derived(weatherInfo(forecast?.current.weatherCode ?? 0, isDay));
  // Before the first forecast (the welcome screen) the sky assumes an ordinary day.
  const sky = $derived.by(() => {
    const rise = sunrise ?? new Date(new Date(now).setHours(6, 30, 0, 0));
    const set = sunset ?? new Date(new Date(now).setHours(18, 30, 0, 0));
    return skyAt(now, rise, set, forecast?.current.cloudCover ?? 0, info.precipitation);
  });

  async function refresh() {
    if (locationStore.current?.kind === "gps") await locationStore.locate();
    await weatherStore.refresh();
  }

  const title = $derived(locationStore.current ? `${locationLabel(locationStore.current)} · Vrijeme` : "Vrijeme");
</script>

<svelte:head><title>{title}</title></svelte:head>

<SkyBackdrop {sky} precipitation={prefersReducedMotion.current ? "none" : info.precipitation} />

{#if !locationStore.current}
  <Welcome />
{:else}
  <main class="mx-auto w-full max-w-[1100px] px-4 pt-[max(1rem,env(safe-area-inset-top))] pb-[max(2.5rem,env(safe-area-inset-bottom))] sm:px-6">
    <AppHeader
      location={locationStore.current}
      updatedAt={weatherStore.updatedAt}
      loading={weatherStore.loading}
      {now}
      onlocation={() => (sheetOpen = true)}
      onrefresh={refresh}
    />
    {#if forecast}
      <p class="text-lift mt-8 text-8xl font-extralight">{formatDegrees(forecast.current.temperature)}</p>
    {/if}
  </main>
  <LocationSheet bind:open={sheetOpen} />
{/if}
```

- [ ] **Step 5: Check, then try it in the browser**

Run: `cd frontend && bun run check && bun run test`

Expected: 0 errors and 0 warnings; tests PASS.

Start both servers:

```bash
cd backend && go run .     # API on :3000
cd frontend && bun run dev # SPA on :5173
```

Open http://localhost:5173 in Chrome with DevTools in a 390×844 touch viewport. Use Sensors → Location for the positions, or the claude-in-chrome skill. Each check below should hold:

1. **First visit.** A navy-to-teal sky and the welcome screen: "Vrijeme", "Koristi moju lokaciju", search and 20 cities. No permission prompt appears on load.
2. **Pick a city.** Clicking "Split" shows the header "Split · Splitsko-dalmatinska županija" and a temperature, and the tab title is "Split · Vrijeme".
3. **Search.** Opening the header button shows a bottom sheet. Typing "Daru" lists "Daruvar · Bjelovarsko-bilogorska županija"; Enter picks it, and the sheet closes. "Nedavno" now lists Daruvar and Split.
4. **My location.** With Sensors set to 45.59/17.22 (Daruvar), the button prompts for permission; once allowed, the header shows a navigation arrow and "Daruvar".
   - Set to Berlin (52.52/13.40): the header reads "Moja lokacija" with no county.
   - With the location blocked in site settings: a Croatian error under the button, and the picker stays usable.
5. **Reload.** The same location and the temperature appear at once from the device cache. "Ažurirano prije …" shows the cache's age, then it updates.
6. **Keyboard.** Tab reaches the header button and the refresh button. In the sheet, Escape closes it; ↓ from the search field moves through results; Enter picks.

- [ ] **Step 6: Commit**

```bash
git add -A frontend
git commit -m "Add the location picker, sheet, welcome screen and header over the living sky

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 13: Forecast UI: current conditions, hourly, 7 days, detail tiles, air quality and warnings

**Files:**
- Create: `frontend/src/lib/components/{current-conditions,hourly-chart,hourly-panel,daily-panel,day-details,detail-tile,detail-tiles,air-quality-card,warning-banner,forecast-skeleton,app-footer}.svelte`
- Modify: `frontend/src/routes/+page.svelte`

**Interfaces:**
- Consumes: everything above.
- Produces these components:
  - `<CurrentConditions forecast today info hint>`;
  - `<HourlyPanel hours days timeZone>`, which wraps `<HourlyChart hours timeZone sunEvents>` and exports `type SunEvent` from its module script;
  - `<DailyPanel forecast now>`, with `<DayDetails day hours timeZone>` inside;
  - `<DetailTiles forecast today>` and `<DetailTile label value detail? icon?>`;
  - `<AirQualityCard airQuality>` and `<WarningBanner warnings now timeZone>`;
  - `<ForecastSkeleton>` and `<AppFooter>`.

- [ ] **Step 1: Write the current conditions and the hourly panel**

`frontend/src/lib/components/current-conditions.svelte`:

```svelte
<script lang="ts">
  import Umbrella from "@lucide/svelte/icons/umbrella";
  import { formatDegrees } from "$lib/helpers/format";
  import type { WeatherInfo } from "$lib/helpers/weather";
  import type { DailyForecast, Forecast } from "$lib/types/forecast";
  import WeatherIcon from "./weather-icon.svelte";

  type Props = { forecast: Forecast; today: DailyForecast | undefined; info: WeatherInfo; hint: string };
  let { forecast, today, info, hint }: Props = $props();

  const wet = $derived(hint !== "" && !hint.startsWith("Bez"));
</script>

<section aria-label="Trenutno vrijeme" class="text-lift flex flex-col items-center py-2 text-center sm:flex-row sm:gap-6 sm:text-left">
  <WeatherIcon name={info.icon} size={148} class="-my-3 shrink-0" />
  <div>
    <p class="text-[5.5rem] leading-none font-extralight tabular-nums sm:text-[6.5rem]">{formatDegrees(forecast.current.temperature)}</p>
    <p class="mt-1 text-2xl font-light">{info.label}</p>
    <p class="text-muted-foreground mt-1">
      Osjećaj {formatDegrees(forecast.current.apparentTemperature)}
      {#if today}
        · <span class="sr-only">najviša</span><span aria-hidden="true">↑</span>{formatDegrees(today.temperatureMax)}
        <span class="sr-only">najniža</span><span aria-hidden="true">↓</span>{formatDegrees(today.temperatureMin)}
      {/if}
    </p>
    {#if hint}
      <p class="mt-2 flex items-center justify-center gap-1.5 sm:justify-start">
        {#if wet}<Umbrella class="size-5" aria-hidden="true" />{/if}{hint}
      </p>
    {/if}
  </div>
</section>
```

`frontend/src/lib/components/hourly-chart.svelte`:

```svelte
<script lang="ts" module>
  export interface SunEvent {
    time: Date;
    kind: "sunrise" | "sunset";
  }
</script>

<script lang="ts">
  import { formatClock, formatDegrees, formatHour } from "$lib/helpers/format";
  import { areaPath, extent, linePath, scale } from "$lib/helpers/sparkline";
  import { weatherInfo } from "$lib/helpers/weather";
  import type { HourlyForecast } from "$lib/types/forecast";
  import WeatherIcon from "./weather-icon.svelte";

  type Props = { hours: HourlyForecast[]; timeZone: string; sunEvents?: SunEvent[] };
  let { hours, timeZone, sunEvents = [] }: Props = $props();

  // One column per hour; the strip scrolls sideways inside its panel. Two charts share the time
  // axis, never two y-scales: temperature above, the chance of rain below.
  const SLOT = 56;
  const TEMP_H = 92;
  const PAD = 24; // room for the labels above the curve
  const RAIN_H = 34;
  const BAR = 20;

  const width = $derived(hours.length * SLOT);
  const cx = (i: number) => (i + 0.5) * SLOT;
  const temps = $derived(hours.map((h) => h.temperature));
  const domain = $derived.by((): [number, number] => {
    const [lo, hi] = extent(temps) ?? [0, 1];
    const grow = Math.max(0, 4 - (hi - lo)) / 2; // a flat day stays flat-looking
    return [lo - grow, hi + grow];
  });
  const y = $derived(scale(domain, [TEMP_H - PAD, PAD]));
  const wet = $derived(hours.some((h) => h.precipitationProbability >= 5));
  const rainY = scale([0, 100], [RAIN_H, 12]);

  // A bar with a rounded top and a square foot on the baseline.
  function bar(x: number, top: number) {
    const r = Math.min(4, (RAIN_H - top) / 2, BAR / 2);
    const l = x - BAR / 2;
    return `M${l},${RAIN_H}V${top + r}Q${l},${top} ${l + r},${top}H${l + BAR - r}Q${l + BAR},${top} ${l + BAR},${top + r}V${RAIN_H}Z`;
  }

  // Sunrise and sunset inside the strip, placed by time between the hour columns.
  const start = $derived(hours.length ? Date.parse(hours[0].time) : 0);
  const sunMarks = $derived(
    sunEvents
      .map((e) => ({ ...e, x: ((e.time.getTime() - start) / 3_600_000) * SLOT }))
      .filter((m) => m.x > 0 && m.x < width),
  );
</script>

<div class="relative pb-6" style:width="{width}px">
  <div class="flex">
    {#each hours as h, i (h.time)}
      <div class="flex flex-col items-center" style:width="{SLOT}px">
        <span class="text-muted-foreground text-sm">{i === 0 ? "Sada" : formatHour(new Date(h.time), timeZone)}</span>
        <WeatherIcon name={weatherInfo(h.weatherCode, h.isDay).icon} size={40} />
      </div>
    {/each}
  </div>

  <svg {width} height={TEMP_H} viewBox="0 0 {width} {TEMP_H}" class="block overflow-visible">
    <g transform="translate({SLOT / 2} 0)">
      <path d={areaPath(temps, width - SLOT, TEMP_H, domain, PAD)} fill="white" opacity="0.08" />
      <path
        d={linePath(temps, width - SLOT, TEMP_H, domain, PAD)}
        fill="none"
        stroke="white"
        stroke-opacity="0.9"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </g>
    {#each temps as t, i (hours[i].time)}
      <text x={cx(i)} y={y(t) - 9} text-anchor="middle" class="fill-white text-[15px] font-medium">{formatDegrees(t)}</text>
    {/each}
  </svg>

  {#if wet}
    <svg {width} height={RAIN_H} viewBox="0 0 {width} {RAIN_H}" class="block overflow-visible">
      <line x1="0" x2={width} y1={RAIN_H - 0.5} y2={RAIN_H - 0.5} stroke="white" stroke-opacity="0.15" />
      {#each hours as h, i (h.time)}
        {#if h.precipitationProbability >= 5}
          <path d={bar(cx(i), rainY(h.precipitationProbability))} fill="var(--color-rain)" opacity="0.85" />
        {/if}
        {#if h.precipitationProbability >= 20}
          <text x={cx(i)} y={rainY(h.precipitationProbability) - 3} text-anchor="middle" class="fill-white/85 text-[12px]">
            {Math.round(h.precipitationProbability)} %
          </text>
        {/if}
      {/each}
    </svg>
  {/if}

  {#each sunMarks as m (m.time.getTime())}
    <div class="pointer-events-none absolute inset-y-0 border-l border-dashed border-amber-200/50" style:left="{m.x}px">
      <span class="absolute bottom-0 left-1 text-xs whitespace-nowrap text-amber-100">
        {m.kind === "sunrise" ? "Izlazak" : "Zalazak"} {formatClock(m.time, timeZone)}
      </span>
    </div>
  {/each}
</div>
```

`frontend/src/lib/components/hourly-panel.svelte`:

```svelte
<script lang="ts">
  import { formatDegrees, formatHour, formatPercent } from "$lib/helpers/format";
  import { weatherInfo } from "$lib/helpers/weather";
  import type { DailyForecast, HourlyForecast } from "$lib/types/forecast";
  import HourlyChart, { type SunEvent } from "./hourly-chart.svelte";

  type Props = { hours: HourlyForecast[]; days: DailyForecast[]; timeZone: string };
  let { hours, days, timeZone }: Props = $props();

  const sunEvents = $derived(
    days
      .flatMap((d): SunEvent[] => [
        { time: new Date(d.sunrise), kind: "sunrise" },
        { time: new Date(d.sunset), kind: "sunset" },
      ])
      .filter((e) => e.time.getUTCFullYear() > 1), // no sunrise (polar day): the zero time
  );
</script>

<section aria-labelledby="hourly-title" class="panel py-4">
  <h2 id="hourly-title" class="text-muted-foreground mb-2 px-5 text-xs font-medium tracking-wider uppercase">Idućih 24 sata</h2>
  <!-- The strip scrolls, so it takes focus for the keyboard; the chart is drawn for the eye, and the
       table below says the same to screen readers. -->
  <!-- svelte-ignore a11y_no_noninteractive_tabindex -->
  <div tabindex="0" role="region" aria-label="Prognoza po satima" class="scrollbar-thin overflow-x-auto px-2">
    <div aria-hidden="true"><HourlyChart {hours} {timeZone} {sunEvents} /></div>
  </div>
  <table class="sr-only">
    <caption>Prognoza po satima</caption>
    <thead>
      <tr><th scope="col">Sat</th><th scope="col">Vrijeme</th><th scope="col">Temperatura</th><th scope="col">Vjerojatnost oborine</th></tr>
    </thead>
    <tbody>
      {#each hours as h (h.time)}
        <tr>
          <th scope="row">{formatHour(new Date(h.time), timeZone)}</th>
          <td>{weatherInfo(h.weatherCode, h.isDay).label}</td>
          <td>{formatDegrees(h.temperature)}</td>
          <td>{formatPercent(h.precipitationProbability)}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</section>
```

- [ ] **Step 2: Write the 7-day panel**

`frontend/src/lib/components/day-details.svelte`:

```svelte
<script lang="ts">
  import {
    formatClock,
    formatDegrees,
    formatDuration,
    formatHour,
    formatMillimetres,
    formatPercent,
    formatSpeed,
    minutesOfDay,
  } from "$lib/helpers/format";
  import { uvBand } from "$lib/helpers/uv";
  import { weatherInfo } from "$lib/helpers/weather";
  import { compassName } from "$lib/helpers/wind";
  import type { DailyForecast, HourlyForecast } from "$lib/types/forecast";
  import WeatherIcon from "./weather-icon.svelte";

  type Props = { day: DailyForecast; hours: HourlyForecast[]; timeZone: string };
  let { day, hours, timeZone }: Props = $props();

  // Every third hour keeps the strip short enough for a phone.
  const shown = $derived(hours.filter((h) => (minutesOfDay(new Date(h.time), timeZone) / 60) % 3 === 0));
</script>

<div class="px-2 pt-1 pb-3">
  {#if shown.length}
    <ol class="scrollbar-thin flex gap-1 overflow-x-auto pb-2">
      {#each shown as h (h.time)}
        <li class="flex min-w-14 flex-col items-center rounded-xl bg-white/6 px-1 py-2 text-sm">
          <span class="text-muted-foreground">{formatHour(new Date(h.time), timeZone)}</span>
          <WeatherIcon name={weatherInfo(h.weatherCode, h.isDay).icon} size={32} />
          <span class="font-medium">{formatDegrees(h.temperature)}</span>
          {#if h.precipitationProbability >= 20}<span class="text-rain text-xs">{formatPercent(h.precipitationProbability)}</span>{/if}
        </li>
      {/each}
    </ol>
  {/if}
  <dl class="mt-2 grid grid-cols-1 gap-x-4 gap-y-1.5 text-sm sm:grid-cols-2">
    <div><dt class="text-muted-foreground inline">Oborine</dt> <dd class="inline">{formatMillimetres(day.precipitationSum)} · {formatPercent(day.precipitationProbabilityMax)}</dd></div>
    <div><dt class="text-muted-foreground inline">Vjetar</dt> <dd class="inline">do {formatSpeed(day.windSpeedMax)}, {compassName(day.windDirectionDominant)} · udari {formatSpeed(day.windGustsMax)}</dd></div>
    <div><dt class="text-muted-foreground inline">UV indeks</dt> <dd class="inline">{Math.round(day.uvIndexMax)} · {uvBand(day.uvIndexMax).label}</dd></div>
    <div><dt class="text-muted-foreground inline">Sunce</dt> <dd class="inline">{formatClock(new Date(day.sunrise), timeZone)} – {formatClock(new Date(day.sunset), timeZone)} ({formatDuration(day.daylightSeconds)})</dd></div>
  </dl>
</div>
```

`frontend/src/lib/components/daily-panel.svelte`:

```svelte
<script lang="ts">
  import ChevronDown from "@lucide/svelte/icons/chevron-down";
  import Umbrella from "@lucide/svelte/icons/umbrella";
  import { hoursOfDay, weekRange } from "$lib/helpers/forecast";
  import { dayLabel, formatDegrees, formatPercent } from "$lib/helpers/format";
  import { weatherInfo } from "$lib/helpers/weather";
  import type { Forecast } from "$lib/types/forecast";
  import DayDetails from "./day-details.svelte";
  import WeatherIcon from "./weather-icon.svelte";

  type Props = { forecast: Forecast; now: Date };
  let { forecast, now }: Props = $props();

  // Every range bar sits on the week's scale, so the days compare at a glance.
  const range = $derived(weekRange(forecast.daily));
  const pct = (t: number) => Math.min(100, Math.max(0, ((t - range[0]) / Math.max(1, range[1] - range[0])) * 100));
</script>

<section aria-labelledby="daily-title" class="panel px-3 py-4 sm:px-4">
  <h2 id="daily-title" class="text-muted-foreground mb-1 px-2 text-xs font-medium tracking-wider uppercase">7 dana</h2>
  <ul>
    {#each forecast.daily as day (day.date)}
      {@const label = dayLabel(day.date, now, forecast.timezone)}
      {@const info = weatherInfo(day.weatherCode)}
      <li class="border-t border-white/8 first:border-t-0">
        <details class="group">
          <summary class="flex cursor-pointer list-none items-center gap-2.5 rounded-xl px-2 py-2.5 hover:bg-white/8 sm:gap-3 [&::-webkit-details-marker]:hidden">
            <span class="w-[4.5rem] shrink-0 font-medium">{label}</span>
            <WeatherIcon name={info.icon} size={40} class="-my-1 shrink-0" />
            <span class="sr-only">{info.label},</span>
            <span class="text-rain flex w-12 shrink-0 items-center gap-0.5 text-sm">
              {#if day.precipitationProbabilityMax >= 20}<Umbrella class="size-3.5" aria-hidden="true" />{formatPercent(day.precipitationProbabilityMax)}{/if}
            </span>
            <span class="text-muted-foreground w-9 shrink-0 text-right tabular-nums"><span class="sr-only">najniža</span>{formatDegrees(day.temperatureMin)}</span>
            <span class="relative h-1.5 flex-1 rounded-full bg-white/12" aria-hidden="true">
              <span
                class="absolute inset-y-0 rounded-full bg-linear-to-r from-sky-300 to-amber-300"
                style:left="{pct(day.temperatureMin)}%"
                style:right="{100 - pct(day.temperatureMax)}%"
              ></span>
              {#if label === "Danas"}
                <span class="absolute top-1/2 size-2.5 -translate-x-1/2 -translate-y-1/2 rounded-full bg-white ring-2 ring-black/20" style:left="{pct(forecast.current.temperature)}%"></span>
              {/if}
            </span>
            <span class="w-9 shrink-0 text-right font-medium tabular-nums"><span class="sr-only">najviša</span>{formatDegrees(day.temperatureMax)}</span>
            <ChevronDown class="size-4 shrink-0 opacity-60 transition group-open:rotate-180" aria-hidden="true" />
          </summary>
          <DayDetails {day} hours={hoursOfDay(forecast, day.date)} timeZone={forecast.timezone} />
        </details>
      </li>
    {/each}
  </ul>
</section>
```

- [ ] **Step 3: Write the detail tiles, the air card and the warnings**

`frontend/src/lib/components/detail-tile.svelte`:

```svelte
<script lang="ts">
  import type { Snippet } from "svelte";

  type Props = { label: string; value: string; detail?: string; icon?: Snippet };
  let { label, value, detail, icon }: Props = $props();
</script>

<div class="panel flex flex-col gap-1 p-4">
  <dt class="text-muted-foreground text-xs font-medium tracking-wider uppercase">{label}</dt>
  <dd class="flex items-center gap-2 text-2xl font-light tabular-nums">{@render icon?.()}{value}</dd>
  {#if detail}<dd class="text-muted-foreground text-sm">{detail}</dd>{/if}
</div>
```

`frontend/src/lib/components/detail-tiles.svelte`:

```svelte
<script lang="ts">
  import Navigation2 from "@lucide/svelte/icons/navigation-2";
  import {
    formatClock,
    formatDegrees,
    formatDistance,
    formatDuration,
    formatPercent,
    formatPressure,
    formatSpeed,
  } from "$lib/helpers/format";
  import { uvBand } from "$lib/helpers/uv";
  import { compassName, windArrowDegrees, windStrength } from "$lib/helpers/wind";
  import type { DailyForecast, Forecast } from "$lib/types/forecast";
  import DetailTile from "./detail-tile.svelte";

  type Props = { forecast: Forecast; today: DailyForecast | undefined };
  let { forecast, today }: Props = $props();

  const c = $derived(forecast.current);
  const tz = $derived(forecast.timezone);
  const uv = $derived(uvBand(c.uvIndex));
</script>

<section aria-labelledby="details-title">
  <h2 id="details-title" class="sr-only">Pojedinosti</h2>
  <dl class="grid grid-cols-2 gap-3 sm:grid-cols-3">
    <DetailTile
      label="Vjetar"
      value={formatSpeed(c.windSpeed)}
      detail="{windStrength(c.windSpeed)}, {compassName(c.windDirection)} · udari {formatSpeed(c.windGusts)}"
    >
      {#snippet icon()}
        <Navigation2 class="size-5 shrink-0 fill-current" style="transform: rotate({windArrowDegrees(c.windDirection)}deg)" aria-hidden="true" />
      {/snippet}
    </DetailTile>
    <DetailTile
      label="UV indeks"
      value={String(Math.round(c.uvIndex))}
      detail={today ? `${uv.label} · danas do ${Math.round(today.uvIndexMax)}` : uv.label}
    />
    <DetailTile label="Vlažnost" value={formatPercent(c.humidity)} detail="Rosište {formatDegrees(c.dewPoint)}" />
    <DetailTile label="Tlak" value={formatPressure(c.pressure)} detail="Na razini mora" />
    <DetailTile label="Vidljivost" value={formatDistance(c.visibility)} />
    {#if today}
      <DetailTile
        label="Sunce"
        value="{formatClock(new Date(today.sunrise), tz)} – {formatClock(new Date(today.sunset), tz)}"
        detail="Dan traje {formatDuration(today.daylightSeconds)}"
      />
    {/if}
  </dl>
</section>
```

`frontend/src/lib/components/air-quality-card.svelte`:

```svelte
<script lang="ts">
  import { activePollen, aqiLabels, pollenLevelLabels, pollenNames } from "$lib/helpers/air";
  import type { AirQuality } from "$lib/types/airquality";

  type Props = { airQuality: AirQuality };
  let { airQuality }: Props = $props();

  const pollen = $derived(activePollen(airQuality));
  // The European AQI's own colour scale, always next to the label that says the same.
  const aqiDot = {
    good: "bg-teal-300",
    fair: "bg-emerald-300",
    moderate: "bg-yellow-300",
    poor: "bg-orange-400",
    very_poor: "bg-red-500",
    extremely_poor: "bg-purple-500",
  };
  const pollenDot = { none: "", low: "bg-emerald-300", moderate: "bg-yellow-300", high: "bg-orange-400", very_high: "bg-red-500" };
</script>

<section aria-labelledby="air-title" class="panel p-5">
  <h2 id="air-title" class="text-muted-foreground text-xs font-medium tracking-wider uppercase">Kvaliteta zraka</h2>
  <p class="mt-2 flex items-center gap-3">
    <span class="size-3 shrink-0 rounded-full {aqiDot[airQuality.aqi.level]}" aria-hidden="true"></span>
    <span class="text-2xl font-light">{aqiLabels[airQuality.aqi.level]}</span>
    <span class="text-muted-foreground">EAQI {Math.round(airQuality.aqi.value)}</span>
  </p>
  <h3 class="text-muted-foreground mt-4 text-xs font-medium tracking-wider uppercase">Pelud danas</h3>
  {#if pollen.length}
    <ul class="mt-2 flex flex-wrap gap-2">
      {#each pollen as p (p.type)}
        <li class="flex items-center gap-2 rounded-full bg-white/10 px-3 py-1 text-sm">
          <span class="size-2 rounded-full {pollenDot[p.level]}" aria-hidden="true"></span>
          {pollenNames[p.type]} <span class="text-muted-foreground">{pollenLevelLabels[p.level]}</span>
        </li>
      {/each}
    </ul>
  {:else}
    <p class="text-muted-foreground mt-1">Nema peludi u zraku.</p>
  {/if}
</section>
```

`frontend/src/lib/components/warning-banner.svelte`:

```svelte
<script lang="ts">
  import TriangleAlert from "@lucide/svelte/icons/triangle-alert";
  import { warningWhen } from "$lib/helpers/alerts";
  import type { Warning } from "$lib/types/warning";

  type Props = { warnings: Warning[]; now: Date; timeZone: string };
  let { warnings, now, timeZone }: Props = $props();

  // Status colours, always with the icon and DHMZ's own words ("Žuto upozorenje za vjetar").
  const fill = { yellow: "bg-warn-yellow text-black", orange: "bg-warn-orange text-black", red: "bg-warn-red text-white" };
</script>

<section aria-label="Upozorenja DHMZ-a" class="flex flex-col gap-2">
  {#each warnings as w (`${w.event}|${w.onset}`)}
    <details class="rounded-2xl {fill[w.level]}">
      <summary class="flex cursor-pointer list-none items-center gap-3 px-4 py-2.5 font-medium [&::-webkit-details-marker]:hidden">
        <TriangleAlert class="size-5 shrink-0" aria-hidden="true" />
        <span class="min-w-0 flex-1">{w.event}</span>
        <span class="shrink-0 text-sm opacity-80">{warningWhen(w, now, timeZone)}</span>
      </summary>
      {#if w.description}<p class="px-4 pb-3 text-sm">{w.description}</p>{/if}
    </details>
  {/each}
</section>
```

- [ ] **Step 4: Write the skeleton and the footer**

`frontend/src/lib/components/forecast-skeleton.svelte`:

```svelte
<div role="status" class="mt-2 grid animate-pulse gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,26rem)]">
  <span class="sr-only">Učitavam prognozu…</span>
  <div class="flex flex-col gap-4">
    <div class="mx-auto my-4 h-36 w-64 rounded-3xl bg-white/10"></div>
    <div class="panel h-56"></div>
    <div class="grid grid-cols-2 gap-3 sm:grid-cols-3">
      {#each [1, 2, 3, 4, 5, 6] as n (n)}<div class="panel h-28"></div>{/each}
    </div>
  </div>
  <div class="panel h-[30rem]"></div>
</div>
```

`frontend/src/lib/components/app-footer.svelte`:

```svelte
<footer class="text-muted-foreground mt-10 text-center text-xs leading-relaxed [&_a]:underline [&_a]:underline-offset-2">
  <p>
    Prognoza i kvaliteta zraka: <a href="https://open-meteo.com/">Open-Meteo</a> (CC BY 4.0) · Upozorenja:
    <a href="https://meteo.hr/">DHMZ</a> putem <a href="https://meteoalarm.org/">Meteoalarm</a> · Mjesta:
    <a href="https://www.geonames.org/">GeoNames</a> (CC BY 4.0)
  </p>
</footer>
```

- [ ] **Step 5: Lay out the page**

Replace `frontend/src/routes/+page.svelte` with:

```svelte
<script lang="ts">
  import { onMount } from "svelte";
  import { prefersReducedMotion } from "svelte/motion";
  import AirQualityCard from "$lib/components/air-quality-card.svelte";
  import AppFooter from "$lib/components/app-footer.svelte";
  import AppHeader from "$lib/components/app-header.svelte";
  import CurrentConditions from "$lib/components/current-conditions.svelte";
  import DailyPanel from "$lib/components/daily-panel.svelte";
  import DetailTiles from "$lib/components/detail-tiles.svelte";
  import ForecastSkeleton from "$lib/components/forecast-skeleton.svelte";
  import HourlyPanel from "$lib/components/hourly-panel.svelte";
  import LocationSheet from "$lib/components/location-sheet.svelte";
  import SkyBackdrop from "$lib/components/sky-backdrop.svelte";
  import WarningBanner from "$lib/components/warning-banner.svelte";
  import Welcome from "$lib/components/welcome.svelte";
  import { todayOf, upcomingHours } from "$lib/helpers/forecast";
  import { locationLabel } from "$lib/helpers/location";
  import { skyAt } from "$lib/helpers/sky";
  import { rainHint, weatherInfo } from "$lib/helpers/weather";
  import { clock } from "$lib/stores/clock.svelte";
  import { locationStore } from "$lib/stores/location.svelte";
  import { weatherStore } from "$lib/stores/weather.svelte";

  $effect(() => clock.subscribe());
  $effect(() => weatherStore.subscribe());
  // Follow the chosen location; the store ignores a repeat of the same ~1 km cell.
  $effect(() => {
    const location = locationStore.current;
    if (location) weatherStore.show(location);
  });
  // A returning GPS reader whose permission stands is followed without a prompt.
  onMount(() => void locationStore.relocateIfGranted());

  let sheetOpen = $state(false);

  const now = $derived(clock.now);
  const forecast = $derived(weatherStore.forecast);
  const today = $derived(forecast ? todayOf(forecast, now) : undefined);
  const sunrise = $derived(today ? new Date(today.sunrise) : undefined);
  const sunset = $derived(today ? new Date(today.sunset) : undefined);
  const isDay = $derived(sunrise && sunset ? now >= sunrise && now < sunset : true);
  const info = $derived(weatherInfo(forecast?.current.weatherCode ?? 0, isDay));
  const hint = $derived(forecast ? rainHint(forecast.hourly, now, forecast.timezone) : "");
  // Before the first forecast (the welcome screen) the sky assumes an ordinary day.
  const sky = $derived.by(() => {
    const rise = sunrise ?? new Date(new Date(now).setHours(6, 30, 0, 0));
    const set = sunset ?? new Date(new Date(now).setHours(18, 30, 0, 0));
    return skyAt(now, rise, set, forecast?.current.cloudCover ?? 0, info.precipitation);
  });

  async function refresh() {
    if (locationStore.current?.kind === "gps") await locationStore.locate();
    await weatherStore.refresh();
  }

  const title = $derived(locationStore.current ? `${locationLabel(locationStore.current)} · Vrijeme` : "Vrijeme");
</script>

<svelte:head><title>{title}</title></svelte:head>

<SkyBackdrop {sky} precipitation={prefersReducedMotion.current ? "none" : info.precipitation} />

{#if !locationStore.current}
  <Welcome />
{:else}
  <main class="mx-auto w-full max-w-[1100px] px-4 pt-[max(1rem,env(safe-area-inset-top))] pb-[max(2.5rem,env(safe-area-inset-bottom))] sm:px-6">
    <AppHeader
      location={locationStore.current}
      updatedAt={weatherStore.updatedAt}
      loading={weatherStore.loading}
      {now}
      onlocation={() => (sheetOpen = true)}
      onrefresh={refresh}
    />

    {#if forecast && weatherStore.warnings?.length}
      <div class="mt-4"><WarningBanner warnings={weatherStore.warnings} {now} timeZone={forecast.timezone} /></div>
    {/if}

    {#if forecast}
      <div class="mt-2 grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,26rem)] lg:items-start">
        <div class="flex min-w-0 flex-col gap-4">
          <CurrentConditions {forecast} {today} {info} {hint} />
          <HourlyPanel hours={upcomingHours(forecast.hourly, now, 24)} days={forecast.daily} timeZone={forecast.timezone} />
          <DetailTiles {forecast} {today} />
        </div>
        <div class="flex min-w-0 flex-col gap-4">
          <DailyPanel {forecast} {now} />
          {#if weatherStore.airQuality}<AirQualityCard airQuality={weatherStore.airQuality} />{/if}
        </div>
      </div>
    {:else if weatherStore.error}
      <div role="alert" class="panel mt-6 p-6 text-center">
        <p class="text-lg">Prognoza trenutno nije dostupna.</p>
        <button type="button" onclick={refresh} class="mt-4 rounded-full bg-white/15 px-5 py-2 font-medium hover:bg-white/25">
          Pokušaj ponovno
        </button>
      </div>
    {:else}
      <ForecastSkeleton />
    {/if}

    <AppFooter />
  </main>
  <LocationSheet bind:open={sheetOpen} />
{/if}
```

- [ ] **Step 6: Check, then look at it**

Run: `cd frontend && bun run check && bun run test`

Expected: 0 errors and 0 warnings; tests PASS.

With both servers running (Task 12, Step 5), check at 390×844 and at 1280×800:

1. **Phone layout** (390 px, single column), top to bottom:
   - the header;
   - the hero, with an animated icon, a large temperature, the Croatian condition, "Osjećaj … · ↑… ↓…" and the rain hint;
   - "Idućih 24 sata", a strip that scrolls sideways with hour labels ("Sada", "15 h", …), icons, a labelled temperature curve, rain bars (only when wet) and dashed sunrise/sunset markers;
   - six tiles (Vjetar with a rotated arrow, UV, Vlažnost, Tlak, Vidljivost, Sunce);
   - "7 dana", with range bars on a shared scale and a dot for the current temperature on "Danas";
   - "Kvaliteta zraka" with pollen chips, or "Nema peludi u zraku.";
   - the footer credits.
2. **Desktop layout** (1280 px): two columns, days and air quality on the right, at most 1100 px wide.
3. **Day details.** Tapping a day expands its every-third-hour strip and the facts. Tapping again collapses it, and Enter and Space work on the focused row.
4. **Warnings.** If DHMZ has a warning for the location's county, a coloured banner appears under the header and expands to DHMZ's description. Otherwise, for a coastal city, cross-check against https://meteoalarm.org.
5. **Forecast down.** Stop the backend and reload: the cached forecast still shows. With a fresh profile (no cache), "Prognoza trenutno nije dostupna." appears with "Pokušaj ponovno".
6. **Reduced motion.** With DevTools → Rendering → `prefers-reduced-motion: reduce`, the icons stand still and no rain or snow falls.

- [ ] **Step 7: Commit**

```bash
git add -A frontend
git commit -m "Redesign the forecast: hero, 24-hour chart, 7 days, detail tiles, air and DHMZ warnings

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 14: Refresh UX: pull-to-refresh, offline and stale status, theme colour and design preview

**Files:**
- Create: `frontend/src/lib/components/pull-to-refresh.svelte`, `frontend/src/lib/components/status-banner.svelte`
- Modify: `frontend/src/routes/+page.svelte`

**Interfaces:**
- Consumes: `parsePreview` and `previewOffset` (Task 9), `oklchToHex` (Task 9) and the stores.
- Produces `<PullToRefresh onrefresh children>` and `<StatusBanner offline stale>`.

- [ ] **Step 1: Write the components**

`frontend/src/lib/components/pull-to-refresh.svelte`:

```svelte
<script lang="ts">
  import RotateCw from "@lucide/svelte/icons/rotate-cw";
  import type { Snippet } from "svelte";

  type Props = { onrefresh: () => Promise<void>; children: Snippet };
  let { onrefresh, children }: Props = $props();

  // Pulling THRESHOLD px (after 1:2 resistance) at the top of the page and letting go refreshes.
  const THRESHOLD = 72;
  const MAX = 112;

  let pull = $state(0);
  let dragging = $state(false);
  let refreshing = $state(false);
  let startY = 0;

  function ontouchstart(e: TouchEvent) {
    if (refreshing || window.scrollY > 0 || e.touches.length !== 1) return;
    startY = e.touches[0].clientY;
    dragging = true;
  }

  function ontouchmove(e: TouchEvent) {
    if (!dragging) return;
    const dy = e.touches[0].clientY - startY;
    pull = dy > 0 && window.scrollY <= 0 ? Math.min(MAX, dy / 2) : 0;
  }

  async function ontouchend() {
    if (!dragging) return;
    dragging = false;
    if (pull < THRESHOLD) {
      pull = 0;
      return;
    }
    refreshing = true;
    pull = THRESHOLD;
    try {
      await onrefresh();
    } finally {
      refreshing = false;
      pull = 0;
    }
  }
</script>

<svelte:window {ontouchstart} {ontouchmove} {ontouchend} ontouchcancel={ontouchend} />

<!-- Touch only, and decorative: the header's refresh button does the same for everyone. -->
<div
  aria-hidden="true"
  class="pointer-events-none fixed inset-x-0 top-0 z-20 flex justify-center"
  style:transform="translateY({pull - 44}px)"
  style:opacity={Math.min(1, pull / THRESHOLD)}
  style:transition={dragging ? "none" : "transform 200ms ease, opacity 200ms ease"}
>
  <div class="mt-[max(0.5rem,env(safe-area-inset-top))] grid size-10 place-items-center rounded-full bg-white/20 shadow-lg backdrop-blur">
    <RotateCw
      class="size-5 {refreshing ? 'animate-spin' : ''}"
      style={refreshing ? undefined : `transform: rotate(${(pull / THRESHOLD) * 270}deg)`}
    />
  </div>
</div>

<!-- Transformed only while pulled: a transform would otherwise trap fixed-position descendants. -->
<div style:transform={pull ? `translateY(${pull / 2}px)` : undefined} style:transition={dragging ? "none" : "transform 200ms ease"}>
  {@render children()}
</div>
```

`frontend/src/lib/components/status-banner.svelte`:

```svelte
<script lang="ts">
  import CloudOff from "@lucide/svelte/icons/cloud-off";
  import WifiOff from "@lucide/svelte/icons/wifi-off";

  type Props = { offline: boolean; stale: boolean };
  let { offline, stale }: Props = $props();
</script>

{#if offline || stale}
  <p role="status" class="panel mt-3 flex items-center gap-2 px-4 py-2 text-sm">
    {#if offline}
      <WifiOff class="size-4 shrink-0" aria-hidden="true" />Niste povezani s internetom. Prikazani su zadnji spremljeni podaci.
    {:else}
      <CloudOff class="size-4 shrink-0" aria-hidden="true" />Prognozu trenutno ne možemo osvježiti. Prikazani su zadnji podaci.
    {/if}
  </p>
{/if}
```

- [ ] **Step 2: Wire them into the page**

In `frontend/src/routes/+page.svelte`:

1. Add the imports:

```ts
  import { page } from "$app/state";
  import { online } from "svelte/reactivity/window";
  import PullToRefresh from "$lib/components/pull-to-refresh.svelte";
  import StatusBanner from "$lib/components/status-banner.svelte";
  import { parsePreview, previewOffset } from "$lib/helpers/preview";
```

   and change the sky import to `import { oklchToHex, skyAt } from "$lib/helpers/sky";`.

2. Replace the `onMount(...)` line and the `now`, `forecast`, `info` and `sky` declarations with this block. Everything else in the script stays:

```ts
  // Design preview (?at=19:40&code=61&cloud=90), read once: the page at another time or in other weather.
  const preview = parsePreview(page.url.searchParams);
  const previewing = preview.at !== undefined || preview.code !== undefined || preview.cloud !== undefined;

  // The sky fades over a minute as the day moves on, but appears at once on the first paint and in
  // a preview, instead of fading in from navy.
  let painted = $state(false);
  onMount(() => {
    // A returning GPS reader whose permission stands is followed without a prompt.
    void locationStore.relocateIfGranted();
    requestAnimationFrame(() => requestAnimationFrame(() => (painted = true)));
  });

  const forecast = $derived(weatherStore.forecast);
  const now = $derived(
    preview.at !== undefined && forecast
      ? new Date(clock.now.getTime() + previewOffset(preview.at, clock.now, forecast.timezone))
      : clock.now,
  );
```

   Keep `today`, `sunrise`, `sunset` and `isDay` as they are, then:

```ts
  const info = $derived(weatherInfo(preview.code ?? forecast?.current.weatherCode ?? 0, isDay));
```

   Keep `hint`, then:

```ts
  const sky = $derived.by(() => {
    const rise = sunrise ?? new Date(new Date(now).setHours(6, 30, 0, 0));
    const set = sunset ?? new Date(new Date(now).setHours(18, 30, 0, 0));
    return skyAt(now, rise, set, preview.cloud ?? forecast?.current.cloudCover ?? 0, info.precipitation);
  });

  // The browser's own bars take the top of the sky.
  $effect(() => {
    document.querySelector('meta[name="theme-color"]')?.setAttribute("content", oklchToHex(sky.top));
  });
```

3. In the markup, give `SkyBackdrop` the transition:

```svelte
<SkyBackdrop {sky} precipitation={prefersReducedMotion.current ? "none" : info.precipitation} transitionMs={painted && !previewing ? 60_000 : 0} />
```

   Wrap `<main>…</main>` in `<PullToRefresh onrefresh={refresh}>…</PullToRefresh>`, and add the status line right after `<AppHeader … />`:

```svelte
    <StatusBanner offline={online.current === false} stale={weatherStore.stale} />
```

- [ ] **Step 3: Check, then try it**

Run: `cd frontend && bun run check && bun run test`

Expected: 0 errors and 0 warnings; tests PASS.

In the browser, at 390×844 with touch emulation:

1. **Pull to refresh.** Dragging down from the top of the page shows a circular arrow that turns as you pull; letting go past the threshold spins it and refreshes ("Ažurirano upravo"). With a GPS location it re-locates first. A short pull springs back without refreshing. Pulling while scrolled down scrolls normally.
2. **Offline.** DevTools → Network → Offline: the banner reads "Niste povezani s internetom…". Back online, the page refreshes by itself and the banner goes away.
3. **Theme colour.** In Chrome's mobile emulation the meta tag carries the sky's top colour. Check `document.querySelector('meta[name=theme-color]').content` in the console.
4. **Design preview.** Each of these shows the matching sky with no fade:
   - `?at=06:35` (dawn: indigo over a muted rose);
   - `?at=12:00` (day blue);
   - `?at=18:40` (afterglow);
   - `?at=23:00` (navy);
   - `?code=63&cloud=95` (greyer, with rain);
   - `?code=73&cloud=100` (snow);
   - `?code=95` (storm, with a faint flash).
5. **Background polling.** Leave the tab for over 10 minutes (or set `STALE_MS` lower temporarily) and come back: one refresh, not one per missed poll.

- [ ] **Step 4: Commit**

```bash
git add -A frontend
git commit -m "Add pull-to-refresh, offline and stale status, a sky-coloured theme and design previews

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 15: Installable app (PWA)

**Files:**
- Create: `frontend/static/favicon.svg`, `frontend/icons/maskable.svg`, `frontend/icons/render.sh`, `frontend/static/manifest.json`, `frontend/src/service-worker.ts`
- Generate: `frontend/static/favicon.png`, which replaces the old one; `frontend/static/apple-touch-icon.png`; `frontend/static/icons/icon-192.png`, `icon-512.png` and `icon-maskable-512.png`
- Modify: `frontend/src/app.html`

- [ ] **Step 1: Draw the icons and render the PNGs**

`frontend/static/favicon.svg`:

```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
  <defs>
    <linearGradient id="sky" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#012a4a"/>
      <stop offset="1" stop-color="#2c7da0"/>
    </linearGradient>
  </defs>
  <rect width="512" height="512" rx="112" fill="url(#sky)"/>
  <circle cx="318" cy="190" r="78" fill="#ffc94d"/>
  <path d="M146 384h220a66 66 0 0 0 0-132a104 104 0 0 0-196 20a56 56 0 0 0-24 112z" fill="#fff"/>
</svg>
```

`frontend/icons/maskable.svg` is full-bleed: the platform crops it to its own shape, so the mark stays inside the 80% safe zone.

```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
  <defs>
    <linearGradient id="sky" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#012a4a"/>
      <stop offset="1" stop-color="#2c7da0"/>
    </linearGradient>
  </defs>
  <rect width="512" height="512" fill="url(#sky)"/>
  <g transform="translate(71.68 71.68) scale(0.72)">
    <circle cx="318" cy="190" r="78" fill="#ffc94d"/>
    <path d="M146 384h220a66 66 0 0 0 0-132a104 104 0 0 0-196 20a56 56 0 0 0-24 112z" fill="#fff"/>
  </g>
</svg>
```

`frontend/icons/render.sh`:

```sh
#!/bin/sh
# Renders the PNG icons in static/ from static/favicon.svg and icons/maskable.svg (needs rsvg-convert).
# Rerun after changing either SVG.
set -e
cd "$(dirname "$0")/.."
mkdir -p static/icons
rsvg-convert -w 48 -h 48 static/favicon.svg -o static/favicon.png
rsvg-convert -w 192 -h 192 static/favicon.svg -o static/icons/icon-192.png
rsvg-convert -w 512 -h 512 static/favicon.svg -o static/icons/icon-512.png
rsvg-convert -w 512 -h 512 icons/maskable.svg -o static/icons/icon-maskable-512.png
rsvg-convert -w 180 -h 180 icons/maskable.svg -o static/apple-touch-icon.png # iOS rounds the corners itself
```

Run: `chmod +x frontend/icons/render.sh && frontend/icons/render.sh && ls -la frontend/static frontend/static/icons`

Expected: the five PNGs exist. Open `frontend/static/icons/icon-512.png` to check it: a white cloud in front of an amber sun on the navy-to-teal gradient.

- [ ] **Step 2: Write the manifest and link the icons**

`frontend/static/manifest.json`. It's `.json`, not `.webmanifest`, because Go's built-in MIME table knows JSON:

```json
{
  "name": "Vrijeme",
  "short_name": "Vrijeme",
  "description": "Vremenska prognoza za vaše mjesto u Hrvatskoj.",
  "lang": "hr",
  "start_url": "/",
  "scope": "/",
  "display": "standalone",
  "background_color": "#012a4a",
  "theme_color": "#012a4a",
  "icons": [
    { "src": "/icons/icon-192.png", "sizes": "192x192", "type": "image/png" },
    { "src": "/icons/icon-512.png", "sizes": "512x512", "type": "image/png" },
    { "src": "/icons/icon-maskable-512.png", "sizes": "512x512", "type": "image/png", "purpose": "maskable" }
  ]
}
```

In `frontend/src/app.html`, replace the `<link rel="icon" …favicon.png />` line with:

```html
		<link rel="icon" href="%sveltekit.assets%/favicon.svg" type="image/svg+xml" />
		<link rel="icon" href="%sveltekit.assets%/favicon.png" sizes="48x48" type="image/png" />
		<link rel="apple-touch-icon" href="%sveltekit.assets%/apple-touch-icon.png" />
		<link rel="manifest" href="%sveltekit.assets%/manifest.json" />
		<meta name="apple-mobile-web-app-title" content="Vrijeme" />
		<meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
```

- [ ] **Step 3: Write the service worker**

`frontend/src/service-worker.ts`:

```ts
/// <reference types="@sveltejs/kit" />
/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />
import { build, files, version } from "$service-worker";

// Precaches the app shell so the installed app opens offline; the last forecast comes from the
// device's storage (stores/weather.svelte.ts). The API is never cached here.
const sw = self as unknown as ServiceWorkerGlobalScope;
const CACHE = `vrijeme-${version}`;
const ASSETS = new Set([...build, ...files]);
const SHELL = "/";

sw.addEventListener("install", (event) => {
  event.waitUntil(
    (async () => {
      const cache = await caches.open(CACHE);
      await cache.addAll([...ASSETS]);
      // spa/v2 answers the shell only to navigations, which it recognizes by Accept: text/html.
      await cache.add(new Request(SHELL, { headers: { Accept: "text/html" } }));
      await sw.skipWaiting();
    })(),
  );
});

sw.addEventListener("activate", (event) => {
  event.waitUntil(
    (async () => {
      for (const key of await caches.keys()) if (key !== CACHE) await caches.delete(key);
      await sw.clients.claim();
    })(),
  );
});

sw.addEventListener("fetch", (event) => {
  const request = event.request;
  if (request.method !== "GET") return;
  const url = new URL(request.url);
  if (url.origin !== sw.location.origin || url.pathname.startsWith("/api/")) return;

  if (ASSETS.has(url.pathname)) {
    event.respondWith(caches.match(url.pathname).then((cached) => cached ?? fetch(request)));
  } else if (request.mode === "navigate") {
    // Network first, so a deploy shows at once; the cached shell when offline.
    event.respondWith(fetch(request).catch(async () => (await caches.match(SHELL)) ?? Response.error()));
  }
});
```

- [ ] **Step 4: Build and serve it from Go, as production does**

```bash
cd frontend && bun run check && bun run build
rm -rf ../backend/public && cp -r build ../backend/public
cd ../backend && go run .
```

In another terminal:

```bash
curl -s localhost:3000/manifest.json | head -3
curl -sI localhost:3000/service-worker.js | grep -i content-type
curl -s -H 'Accept: text/html' localhost:3000/ | grep -o '<html lang="hr">'
curl -s localhost:3000/api/v1/nope
```

Expected:
- the manifest JSON;
- `text/javascript` (charset optional);
- `<html lang="hr">`;
- `{"code":404,"message":"Not Found"}` (the exact message may differ).

Then open http://localhost:3000 in Chrome:
- **Application → Manifest:** the name is "Vrijeme", the icons render, and there are no installability errors.
- **Service Workers:** activated.
- **Offline:** tick Offline and reload. The page loads from the shell and shows the cached forecast with the offline banner.
- **Lighthouse** (mobile): Accessibility at least 95, Best Practices with no "Requests the geolocation permission on page load", and an installable PWA.

Stop the server and `rm -rf backend/public` (it's git-ignored).

- [ ] **Step 5: Commit**

```bash
git add -A frontend
git commit -m "Make the app installable: manifest, icons and an offline app shell

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

## Task 16: README, Docker and the final verification

**Files:**
- Create: `README.md`
- Delete: `frontend/README.md`

- [ ] **Step 1: Write the README**

`README.md`:

~~~~markdown
# Vrijeme

[vrijeme.app](https://vrijeme.app): the weather where you are, in Croatian. It shows now, every hour and the next seven days, plus DHMZ warnings for your county and today's air quality and pollen. The page sits over a sky that follows the time of day and the weather.

You pick a location in one of three ways:
- your device's position;
- a search for any Croatian settlement;
- one of 20 cities, one per county.

The location is remembered on the device only, rounded to about a kilometre. Pull down, or press ↻, to refresh; with your own position, that also re-locates. The app installs to the home screen (PWA).

A Go/[Raptor](https://github.com/go-raptor/raptor) backend serves the JSON API and the SvelteKit SPA from one origin.

## API

| Route | |
|---|---|
| `GET /api/v1/forecast?lat&lon` | Now, hourly to the end of the week, and 7 days (Open-Meteo; cached 10 min per ~1 km cell) |
| `GET /api/v1/air-quality?lat&lon` | Today's European AQI and pollen (cached 30 min) |
| `GET /api/v1/warnings?lat&lon` | DHMZ warnings for that county via Meteoalarm, in force or within 48 h (cached 15 min) |
| `GET /api/v1/places?q=` | Croatian settlements by name (Open-Meteo geocoding; cached 24 h) |
| `GET /api/v1/places/nearest?lat&lon` | The settlement at a point, from the embedded GeoNames dataset (404 abroad) |
| `GET /healthz`, `GET /readyz` | Liveness and readiness probes |

Each upstream is served stale while it fails. The API is limited to 20 requests a second per client IP.

## Configure

Nothing is required: there are no API keys. The `Dockerfile` sets production up for the HTTPS reverse proxy:
- `SERVER_ADDRESS=0.0.0.0`;
- `SERVER_IP_EXTRACTOR=x-forwarded-for`, so rate limits and logs see the client rather than the proxy;
- `APP_SECURE_HSTS_MAX_AGE=300`; raise it to `31536000` once HTTPS is settled.

`APP_OPENMETEO_URL`, `APP_AIRQUALITY_URL`, `APP_GEOCODING_URL` and `APP_METEOALARM_URL` override the upstreams; the tests use them.

## Develop

```bash
# terminal 1: the API on :3000 (reads backend/.raptor.dev.yaml)
cd backend && raptor dev          # or: go run .

# terminal 2: the SPA on :5173, proxying /api to :3000
cd frontend && bun install && bun run dev
```

Open http://localhost:5173. In Chrome DevTools, Sensors → Location simulates a position.

### Design preview

Query parameters show the page at another time or in other weather:
- `?at=06:35` (dawn) and `?at=18:40` (afterglow);
- `?code=63&cloud=95` (rain), `?code=73&cloud=100` (snow) and `?code=95` (storm).

`at` is a time on the location's clock, `code` a WMO weather code, and `cloud` a cloud cover in percent.

### Generated files

- **Places:** `backend/app/services/data/hr-places.tsv` comes from GeoNames. Regenerate it with `cd backend && go run ./cmd/genplaces`.
- **Icons:** the PNG icons are rendered from `frontend/static/favicon.svg` and `frontend/icons/maskable.svg` by `frontend/icons/render.sh`, which needs `rsvg-convert`.

## Check

```bash
cd backend && go vet ./... && go test ./...
cd frontend && bun run check && bun run test && bun run build
```

## Deploy

The `Dockerfile` builds one image: the Go binary plus the SPA in `public/`. Pushing to the `production` branch builds and publishes it (see `.github/workflows`).

Probes: `GET /healthz` (liveness) and `GET /readyz` (readiness, which answers 503 from the moment shutdown begins), with a timeout of at least 3 s.

## Data

- Forecasts, air quality and geocoding: [Open-Meteo](https://open-meteo.com/), CC BY 4.0
- Warnings: [DHMZ](https://meteo.hr/) via [Meteoalarm](https://meteoalarm.org/)
- Places: [GeoNames](https://www.geonames.org/), CC BY 4.0
- Weather icons: [Meteocons](https://github.com/basmilius/weather-icons) by Bas Milius, MIT
~~~~

Then: `git rm -q frontend/README.md`

- [ ] **Step 2: Run every check**

```bash
cd backend && go vet ./... && go test ./...
cd ../frontend && bun run check && bun run test && bun run build
```

Expected: everything PASSES, and `svelte-check` reports 0 errors and 0 warnings.

- [ ] **Step 3: Build and smoke-test the image**

```bash
docker build -t vrijeme:redesign .    # or: podman build -t vrijeme:redesign .
docker run --rm -d -p 3000:3000 --name vrijeme-test vrijeme:redesign
curl -s localhost:3000/healthz
curl -s -H 'Accept: text/html' localhost:3000/ | grep -o '<html lang="hr">'
curl -s 'localhost:3000/api/v1/places/nearest?lat=45.59&lon=17.22'
curl -s 'localhost:3000/api/v1/forecast?lat=45.59&lon=17.22' | head -c 200
docker stop vrijeme-test
```

Expected:
- `{"status":"ok"}`;
- `<html lang="hr">`;
- `{"name":"Daruvar","county":"Bjelovarsko-bilogorska županija",…}`;
- a forecast from the live Open-Meteo API.

If neither docker nor podman is installed, say so and skip this step; don't claim it passed.

- [ ] **Step 4: Commit**

```bash
git add -A README.md frontend/README.md
git commit -m "Document the app: API, configuration, development, checks and data sources

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```
