package cluster

import (
	"context"

	"github.com/ilexPar/simple-kube/pkg/base"
	"github.com/ilexPar/simple-kube/pkg/cluster/resources"
)

// clusterBackend adapts a cluster resource API to core.Backend. Cluster-scoped
// resources take no namespace, so this is a straight delegation.
type clusterBackend[R base.KubernetesResources] struct {
	api resources.ClusterResourceAPI[R]
}

func (b clusterBackend[R]) Get(ctx context.Context, id string) (*R, error) {
	b.api.SetContext(ctx)
	return b.api.Get(id)
}

func (b clusterBackend[R]) Create(ctx context.Context, obj *R) error {
	b.api.SetContext(ctx)
	return b.api.Create(obj)
}

func (b clusterBackend[R]) Update(ctx context.Context, obj *R) error {
	b.api.SetContext(ctx)
	return b.api.Update(obj)
}

func (b clusterBackend[R]) Delete(ctx context.Context, id string) error {
	b.api.SetContext(ctx)
	return b.api.Delete(id)
}

func (b clusterBackend[R]) List(ctx context.Context, opts base.QueryOpts) ([]R, error) {
	b.api.SetContext(ctx)
	b.api.SetOpts(opts)
	return b.api.List()
}
