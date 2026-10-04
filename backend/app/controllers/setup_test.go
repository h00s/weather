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
