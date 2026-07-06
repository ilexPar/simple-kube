package namespaced

import (
	"context"

	"github.com/ilexPar/simple-kube/pkg/base"
	"github.com/ilexPar/simple-kube/pkg/namespaced/resources"
)

// nsBackend adapts a namespaced resource API to core.Backend by closing over
// the namespace, which is the only namespaced-specific piece of the CRUD flow.
type nsBackend[R base.KubernetesResources] struct {
	api       resources.NamespacedResourceAPI[R]
	namespace string
}

func (b nsBackend[R]) Get(ctx context.Context, id string) (*R, error) {
	b.api.SetContext(ctx)
	return b.api.Get(id, b.namespace)
}

func (b nsBackend[R]) Create(ctx context.Context, obj *R) error {
	b.api.SetContext(ctx)
	return b.api.Create(b.namespace, obj)
}

func (b nsBackend[R]) Update(ctx context.Context, obj *R) error {
	b.api.SetContext(ctx)
	return b.api.Update(b.namespace, obj)
}

func (b nsBackend[R]) Delete(ctx context.Context, id string) error {
	b.api.SetContext(ctx)
	return b.api.Delete(id, b.namespace)
}

func (b nsBackend[R]) List(ctx context.Context, opts base.QueryOpts) ([]R, error) {
	b.api.SetContext(ctx)
	b.api.SetOpts(opts)
	return b.api.List(b.namespace)
}
