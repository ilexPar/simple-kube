package cluster

import (
	"context"

	"github.com/ilexPar/simple-kube/pkg/base"
	"github.com/ilexPar/simple-kube/pkg/cluster/resources"

	"k8s.io/client-go/kubernetes"
)

type Action[T ClusterResources] struct {
	resource T
	api      resources.ClusterResourceAPI
	opts     base.QueryOpts
}

func NewClusterAction[T ClusterResources](
	resource T,
	api resources.ClusterResourceAPI,
) *Action[T] {
	return &Action[T]{
		resource: resource,
		api:      api,
	}
}

func (ca *Action[T]) Get(name string) ClusterGetInterface[T] {
	return &ClusterGet[T]{
		*ca,
		name,
		nil,
	}
}

func (ca *Action[T]) Create(resource T) ClusterPutInterface[T] {
	return &ClusterCreate[T]{
		*ca,
		resource,
		nil,
	}
}

func (ca *Action[T]) Update(resource T) ClusterPutInterface[T] {
	return &ClusterUpdate[T]{
		*ca,
		resource,
		nil,
	}
}

func (ca *Action[T]) List() ClusterListInterface[T] {
	return &ClusterList[T]{
		*ca,
	}
}

func (ca *Action[T]) Delete(resource string) ClusterDeleteInterface[T] {
	return &ClusterDelete[T]{
		*ca,
		resource,
	}
}

func NewQuery(ctx context.Context, client kubernetes.Interface) *Query {
	return &Query{
		ctx:    ctx,
		client: client,
	}
}

type Query struct {
	ctx    context.Context
	client kubernetes.Interface
}

func (c *Query) getResourceAPI(
	res ClusterResourceMethods,
) resources.ClusterResourceAPI {
	api := res.API()
	api.Config(c.ctx, c.client)
	return api
}

func (c *Query) Namespace() ClusterAction[resources.Namespace] {
	res := resources.Namespace{}
	return NewClusterAction(
		res,
		c.getResourceAPI(res),
	)
}
