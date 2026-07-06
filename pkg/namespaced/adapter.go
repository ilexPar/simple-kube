package namespaced

import (
	"github.com/ilexPar/simple-kube/pkg/base"
	"github.com/ilexPar/simple-kube/pkg/namespaced/resources"
)

// nsBackend adapts a namespaced resource API to core.Backend by closing over
// the namespace, which is the only namespaced-specific piece of the CRUD flow.
type nsBackend[R base.KubernetesResources] struct {
	api       resources.NamespacedResourceAPI[R]
	namespace string
}

func (b nsBackend[R]) Get(id string) (*R, error) {
	return b.api.Get(id, b.namespace)
}

func (b nsBackend[R]) Create(obj *R) error {
	return b.api.Create(b.namespace, obj)
}

func (b nsBackend[R]) Update(obj *R) error {
	return b.api.Update(b.namespace, obj)
}

func (b nsBackend[R]) Delete(id string) error {
	return b.api.Delete(id, b.namespace)
}

func (b nsBackend[R]) List(opts base.QueryOpts) ([]R, error) {
	b.api.SetOpts(opts)
	return b.api.List(b.namespace)
}
