package namespaced

import (
	"context"

	"github.com/ilexPar/simple-kube/pkg/base"
	"github.com/ilexPar/simple-kube/pkg/namespaced/resources"
	skns "github.com/ilexPar/simple-kube/pkg/namespaced/resources"

	"k8s.io/client-go/kubernetes"
)

type Action[T NamespacedResources] struct {
	namespace string
	resource  T
	api       resources.NamespacedResourceAPI
	opts      base.QueryOpts
}

func NewAction[T NamespacedResources](
	namespace string,
	resource T,
	api resources.NamespacedResourceAPI,
) *Action[T] {
	return &Action[T]{
		namespace: namespace,
		resource:  resource,
		api:       api,
	}
}

func (ns *Action[T]) Get(name string) NamespacedGetInterface[T] {
	return &NamespacedGet[T]{
		*ns,
		name,
		nil,
	}
}

func (ns *Action[T]) Create(resource T) NamespacedPutInterface[T] {
	return &NamespacedCreate[T]{
		*ns,
		resource,
		nil,
	}
}

func (ns *Action[T]) Update(resource T) NamespacedPutInterface[T] {
	return &NamespacedUpdate[T]{
		*ns,
		resource,
		nil,
	}
}

func (ns *Action[T]) List() NamespacedListInterface[T] {
	return &NamespacedList[T]{
		*ns,
	}
}

func (ns *Action[T]) Delete(resource string) NamespacedDeleteInterface[T] {
	return &NamespacedDelete[T]{
		*ns,
		resource,
	}
}

func NewQuery(
	namespace string,
	ctx context.Context,
	client kubernetes.Interface,
) *Query {
	return &Query{
		namespace: namespace,
		ctx:       ctx,
		client:    client,
	}
}

type Query struct {
	namespace string
	ctx       context.Context
	client    kubernetes.Interface
}

func (n *Query) getResourceAPI(
	res NamespacedResourceMethods,
) skns.NamespacedResourceAPI {
	api := res.API()
	api.Config(n.ctx, n.client)
	return api
}

func (n *Query) Deployment() NamespacedAction[skns.Deployment] {
	res := skns.Deployment{}
	return NewAction(
		n.namespace,
		res,
		n.getResourceAPI(res),
	)
}

func (n *Query) Service() NamespacedAction[skns.Service] {
	res := skns.Service{}
	return NewAction(
		n.namespace,
		res,
		n.getResourceAPI(res),
	)
}

func (n *Query) Job() NamespacedAction[skns.Job] {
	res := skns.Job{}
	return NewAction(
		n.namespace,
		res,
		n.getResourceAPI(res),
	)
}

func (n *Query) CronJob() NamespacedAction[skns.CronJob] {
	res := skns.CronJob{}
	return NewAction(
		n.namespace,
		res,
		n.getResourceAPI(res),
	)
}

func (n *Query) ConfigMap() NamespacedAction[skns.ConfigMap] {
	res := skns.ConfigMap{}
	return NewAction(
		n.namespace,
		res,
		n.getResourceAPI(res),
	)
}

func (n *Query) Ingress() NamespacedAction[skns.Ingress] {
	res := skns.Ingress{}
	return NewAction(
		n.namespace,
		res,
		n.getResourceAPI(res),
	)
}
