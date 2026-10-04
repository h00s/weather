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
