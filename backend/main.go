package main

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/weather/config"
	"github.com/h00s/weather/config/components"
)

func main() {
	raptor.New(components.New(), config.Routes()).Run()
}
