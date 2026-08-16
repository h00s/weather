module github.com/h00s/weather

go 1.26

// replace github.com/go-raptor/raptor/v4 => ../../go-raptor/raptor/v4

require (
	github.com/go-raptor/controllers/spa/v2 v2.0.0
	github.com/go-raptor/middlewares/cors v1.0.10
	github.com/go-raptor/middlewares/logger v1.0.7
	github.com/go-raptor/raptor/v4 v4.3.1
	github.com/h00s/goopenmeteo v1.0.1
	github.com/lmittmann/tint v1.2.0
)

require (
	github.com/andybalholm/brotli v1.2.2 // indirect
	github.com/go-raptor/connectors v1.1.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)
