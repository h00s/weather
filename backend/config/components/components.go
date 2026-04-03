package components

import (
	"github.com/go-raptor/raptor/v4"
)

func New() *raptor.Components {
	return &raptor.Components{
		Controllers: Controllers(),
		Services:    Services(),
		Middlewares: Middlewares(),
	}
}
