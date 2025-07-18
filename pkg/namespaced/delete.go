package namespaced

import "github.com/ilexPar/simple-kube/pkg/base"

type NamespacedDelete[T NamespacedResources, R base.KubernetesResources] struct {
	Action[T, R]
	Id string
}

func (d *NamespacedDelete[T, R]) Run() error {
	return d.api.Delete(d.Id, d.namespace)
}
