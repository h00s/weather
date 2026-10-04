package main

import (
	"log/slog"
	"os"
	_ "time/tzdata" // forecasts are placed in their IANA zone wherever the binary runs

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/config"
	"github.com/h00s/weather/config/components"
	"github.com/lmittmann/tint"
)

func main() {
	raptor.New(
		components.New(),
		config.Routes(),
		raptor.WithLogHandler(func(level *slog.LevelVar) slog.Handler {
			return tint.NewTextHandler(os.Stderr, &tint.Options{Level: level})
		}),
	).Run()
}
