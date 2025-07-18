package namespaced

import (
	"github.com/ilexPar/simple-kube/pkg/base"
	"github.com/ilexPar/simple-kube/pkg/errors"
	sm "github.com/ilexPar/struct-marshal/pkg"
)

type NamespacedGet[T NamespacedResources, R base.KubernetesResources] struct {
	Action[T, R]
	Id       string
	callback func(*R) error
}

func (g *NamespacedGet[T, R]) Run() (T, error) {
	res := new(T)

	obj, err := g.api.Get(g.Id, g.namespace)
	if err != nil {
		return *res, errors.Format(err)
	}

	if g.callback != nil {
		if err = g.callback(obj); err != nil {
			return *res, err
		}
	}

	err = sm.Unmarshal(obj, res)
	return *res, err
}

func (g *NamespacedGet[T, R]) DataHandler(handler func(*R) error) NamespacedGetInterface[T, R] {
	g.callback = handler
	return g
}
