package cluster

import (
	"context"

	api "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/ilexPar/simple-kube/pkg/base"
	"github.com/ilexPar/simple-kube/pkg/cluster/resources"
	"github.com/ilexPar/simple-kube/pkg/core"
)

type Query struct {
	ctx    context.Context
	client kubernetes.Interface
}

func (c *Query) Config(ctx context.Context, client kubernetes.Interface) *Query {
	c.ctx = ctx
	c.client = client
	return c
}

// newAction wires a cluster resource API into the shared generic core action.
func newAction[R base.KubernetesResources, T any](
	q *Query,
	api resources.ClusterResourceAPI[R],
) core.ScopeAction[T, R] {
	api.Config(q.ctx, q.client)
	backend := clusterBackend[R]{api: api}
	return core.NewAction(core.SMCodec[T, R]{}, backend)
}

func (c *Query) Namespace() core.ScopeAction[resources.Namespace, api.Namespace] {
	return newAction[api.Namespace, resources.Namespace](c, &resources.NamespaceAPI[api.Namespace]{})
}
