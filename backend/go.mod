module github.com/h00s/weather

go 1.24.3

// replace github.com/go-raptor/raptor/v4 => ../../go-raptor/raptor/v4

require (
	github.com/go-logistiq/logistiq v1.0.14
	github.com/go-raptor/controllers/spa v1.0.3
	github.com/go-raptor/middlewares/cors v1.0.2
	github.com/go-raptor/raptor/v4 v4.0.0-alpha.14
	github.com/h00s/goopenmeteo v1.0.1
)

require (
	github.com/go-raptor/connectors v1.0.2 // indirect
	github.com/h00s/tint v1.0.7 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/nats-io/nats.go v1.42.0 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
