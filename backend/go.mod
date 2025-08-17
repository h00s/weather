module github.com/h00s/weather

go 1.25

// replace github.com/go-raptor/raptor/v4 => ../../go-raptor/raptor/v4

require (
	github.com/go-logistiq/logistiq v1.0.17
	github.com/go-raptor/controllers/spa v1.0.3
	github.com/go-raptor/middlewares/cors v1.0.3
	github.com/go-raptor/middlewares/logger v1.0.0
	github.com/go-raptor/raptor/v4 v4.0.6
	github.com/h00s/goopenmeteo v1.0.1
)

require (
	github.com/go-raptor/connectors v1.0.2 // indirect
	github.com/h00s/tint v1.1.2 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/nats-io/nats.go v1.44.0 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
