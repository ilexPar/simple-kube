package cluster

import (
	"github.com/ilexPar/simple-kube/pkg/base"
)

type ClusterList[T ClusterResources] struct {
	Action[T]
}

func (l *ClusterList[T]) Run() ([]T, error) {
	res := []T{}
	l.api.SetOpts(l.opts)
	objs, err := l.api.List()
	if err != nil {
		return res, err
	}
	for _, obj := range objs {
		resource := new(T)
		if err = l.resource.Load(obj, resource); err != nil {
			return res, err
		}
		res = append(res, *resource)
	}

	return res, err
}

func (l *ClusterList[T]) FilterByLabels(labels map[string]string) ClusterListInterface[T] {
	l.opts.List.LabelSelector = base.FlattenLabels(labels)
	return l
}
