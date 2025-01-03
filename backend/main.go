package main

import (
	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/weather/config/initializers"
)

func main() {
	r := raptor.NewRaptor()

	r.Init(initializers.App(r.Utils.Config))
	r.Listen()
}
