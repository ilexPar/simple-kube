package namespaced

import (
	"github.com/ilexPar/simple-kube/pkg/base"
)

type NamespacedList[T NamespacedResources] struct {
	Action[T]
}

func (l *NamespacedList[T]) Run() ([]T, error) {
	res := []T{}
	l.api.SetOpts(l.opts)
	objs, err := l.api.List(l.namespace)
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

func (l *NamespacedList[T]) FilterByLabels(labels map[string]string) NamespacedListInterface[T] {
	l.opts.List.LabelSelector = base.FlattenLabels(labels)
	return l
}
