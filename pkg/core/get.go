package core

import (
	"github.com/ilexPar/simple-kube/pkg/errors"
)

type Get[T any, R any] struct {
	Action[T, R]
	Id       string
	callback func(*R) error
}

func (g *Get[T, R]) DataHandler(handler func(*R) error) GetInterface[T, R] {
	g.callback = handler
	return g
}

func (g *Get[T, R]) Run() (T, error) {
	var out T

	raw, err := g.Backend.Get(g.Id)
	if err != nil {
		return out, errors.Format(err)
	}

	if g.callback != nil {
		if err = g.callback(raw); err != nil {
			return out, err
		}
	}

	err = g.Codec.Decode(raw, &out)
	return out, err
}
