package cluster

import "github.com/ilexPar/simple-kube/pkg/errors"

type ClusterGet[T ClusterResources] struct {
	Action[T]
	Id       string
	callback func(interface{}) error
}

func (g *ClusterGet[T]) Run() (T, error) {
	res := new(T)

	obj, err := g.api.Get(g.Id)
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

func (g *ClusterGet[T]) DataHandler(handler func(interface{}) error) ClusterGetInterface[T] {
	g.callback = handler
	return g
}
