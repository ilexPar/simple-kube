package namespaced

import "github.com/ilexPar/simple-kube/pkg/errors"

type NamespacedGet[T NamespacedResources] struct {
	Action[T]
	Id       string
	callback func(interface{}) error
}

func (g *NamespacedGet[T]) Run() (T, error) {
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

	err = g.resource.Load(obj, res)
	return *res, err
}

func (g *NamespacedGet[T]) DataHandler(handler func(interface{}) error) NamespacedGetInterface[T] {
	g.callback = handler
	return g
}
