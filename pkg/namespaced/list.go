package namespaced

import (
	"github.com/ilexPar/simple-kube/pkg/base"
	sm "github.com/ilexPar/struct-marshal/pkg"
)

type NamespacedList[T NamespacedResources, R base.KubernetesResources] struct {
	Action[T, R]
}

func (l *NamespacedList[T, R]) Run() ([]T, error) {
	res := []T{}
	l.api.SetOpts(l.opts)
	objs, err := l.api.List(l.namespace)
	if err != nil {
		return res, err
	}
	for _, obj := range objs {
		resource := new(T)
		if err = sm.Unmarshal(obj, resource); err != nil {
			return res, err
		}
		res = append(res, *resource)
	}

	return res, err
}

func (l *NamespacedList[T, R]) FilterByLabels(labels map[string]string) NamespacedListInterface[T, R] {
	l.opts.List.LabelSelector = base.FlattenLabels(labels)
	return l
}
