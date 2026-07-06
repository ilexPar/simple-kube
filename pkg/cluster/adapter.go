package cluster

import (
	"github.com/ilexPar/simple-kube/pkg/base"
	"github.com/ilexPar/simple-kube/pkg/cluster/resources"
)

// clusterBackend adapts a cluster resource API to core.Backend. Cluster-scoped
// resources take no namespace, so this is a straight delegation.
type clusterBackend[R base.KubernetesResources] struct {
	api resources.ClusterResourceAPI[R]
}

func (b clusterBackend[R]) Get(id string) (*R, error) {
	return b.api.Get(id)
}

func (b clusterBackend[R]) Create(obj *R) error {
	return b.api.Create(obj)
}

func (b clusterBackend[R]) Update(obj *R) error {
	return b.api.Update(obj)
}

func (b clusterBackend[R]) Delete(id string) error {
	return b.api.Delete(id)
}

func (b clusterBackend[R]) List(opts base.QueryOpts) ([]R, error) {
	b.api.SetOpts(opts)
	return b.api.List()
}
