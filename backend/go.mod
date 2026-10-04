module github.com/h00s/weather

go 1.27

// replace github.com/go-raptor/raptor/v4 => ../../go-raptor/raptor/v4

require (
	github.com/go-raptor/controllers/spa/v2 v2.1.0
	github.com/go-raptor/middlewares/csrf v1.1.0
	github.com/go-raptor/middlewares/limiter v1.1.0
	github.com/go-raptor/middlewares/logger v1.4.0
	github.com/go-raptor/middlewares/requestid v1.0.0
	github.com/go-raptor/middlewares/secure v1.0.0
	github.com/go-raptor/raptor/v4 v4.6.1
	github.com/h00s/goopenmeteo v1.1.0
	github.com/lmittmann/tint v1.2.0
)

require (
	github.com/andybalholm/brotli v1.2.2 // indirect
	github.com/go-raptor/connectors v1.1.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/time v0.16.0 // indirect
)
