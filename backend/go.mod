module github.com/h00s/weather

go 1.26

// replace github.com/go-raptor/raptor/v4 => ../../go-raptor/raptor/v4

require (
	github.com/go-raptor/controllers/spa v1.0.3
	github.com/go-raptor/middlewares/cors v1.0.3
	github.com/go-raptor/middlewares/logger v1.0.1
	github.com/go-raptor/raptor/v4 v4.1.5
	github.com/h00s/goopenmeteo v1.0.1
)

require (
	github.com/go-raptor/connectors v1.0.2 // indirect
	github.com/h00s/tint v1.1.4 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
