package namespaced

import (
	"context"

	apps "k8s.io/api/apps/v1"
	scaling "k8s.io/api/autoscaling/v2"
	batch "k8s.io/api/batch/v1"
	api "k8s.io/api/core/v1"
	net "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/ilexPar/simple-kube/pkg/base"
	"github.com/ilexPar/simple-kube/pkg/namespaced/resources"
	skns "github.com/ilexPar/simple-kube/pkg/namespaced/resources"
)

type Action[T NamespacedResources, R base.KubernetesResources] struct {
	namespace string
	resource  T
	api       resources.NamespacedResourceAPI[R]
	opts      base.QueryOpts
}

func NewAction[R base.KubernetesResources, T NamespacedResources](
	namespace string,
	resource T,
	api resources.NamespacedResourceAPI[R],
) *Action[T, R] {
	return &Action[T, R]{
		namespace: namespace,
		resource:  resource,
		api:       api,
	}
}

func (ns *Action[T, R]) Get(name string) NamespacedGetInterface[T, R] {
	return &NamespacedGet[T, R]{
		*ns,
		name,
		nil,
	}
}

func (ns *Action[T, R]) Create(resource T) NamespacedPutInterface[T, R] {
	return &NamespacedCreate[T, R]{
		*ns,
		resource,
		nil,
	}
}

func (ns *Action[T, R]) Update(resource T) NamespacedPutInterface[T, R] {
	return &NamespacedUpdate[T, R]{
		*ns,
		resource,
		nil,
	}
}

func (ns *Action[T, R]) List() NamespacedListInterface[T, R] {
	return &NamespacedList[T, R]{
		*ns,
	}
}

func (ns *Action[T, R]) Delete(resource string) NamespacedDeleteInterface[T, R] {
	return &NamespacedDelete[T, R]{
		*ns,
		resource,
	}
}

type Query struct {
	ctx       context.Context
	client    kubernetes.Interface
	namespace string
}

func (c *Query) Config(
	ctx context.Context,
	client kubernetes.Interface,
	namespace string,
) *Query {
	c.ctx = ctx
	c.client = client
	c.namespace = namespace
	return c
}

func GetNamespacedAPI[R base.KubernetesResources, T NamespacedResources](
	ctx context.Context,
	client kubernetes.Interface,
	res T,
) skns.NamespacedResourceAPI[R] {
	var api skns.NamespacedResourceAPI[R]

	switch any(res).(type) {
	case skns.Deployment:
		api = &skns.DeploymentAPI[R]{}
	case skns.Service:
		api = &skns.ServiceAPI[R]{}
	case skns.Job:
		api = &skns.JobAPI[R]{}
	case skns.CronJob:
		api = &skns.CronJobAPI[R]{}
	case skns.ConfigMap:
		api = &skns.ConfigMapAPI[R]{}
	case skns.Ingress:
		api = &skns.IngressAPI[R]{}
	case skns.HPA:
		api = &skns.HPAapi[R]{}
	default:
		panic("cannot resolve Kube API")
	}

	api.Config(ctx, client)
	return api
}

func (n *Query) Deployment() NamespacedAction[skns.Deployment, apps.Deployment] {
	res := skns.Deployment{}
	return NewAction[apps.Deployment](
		n.namespace,
		res,
		GetNamespacedAPI[apps.Deployment](n.ctx, n.client, res),
	)
}

func (n *Query) Service() NamespacedAction[skns.Service, api.Service] {
	res := skns.Service{}
	return NewAction[api.Service](
		n.namespace,
		res,
		GetNamespacedAPI[api.Service](n.ctx, n.client, res),
	)
}

func (n *Query) Job() NamespacedAction[skns.Job, batch.Job] {
	res := skns.Job{}
	return NewAction[batch.Job](
		n.namespace,
		res,
		GetNamespacedAPI[batch.Job](n.ctx, n.client, res),
	)
}

func (n *Query) CronJob() NamespacedAction[skns.CronJob, batch.CronJob] {
	res := skns.CronJob{}
	return NewAction[batch.CronJob](
		n.namespace,
		res,
		GetNamespacedAPI[batch.CronJob](n.ctx, n.client, res),
	)
}

func (n *Query) ConfigMap() NamespacedAction[skns.ConfigMap, api.ConfigMap] {
	res := skns.ConfigMap{}
	return NewAction[api.ConfigMap](
		n.namespace,
		res,
		GetNamespacedAPI[api.ConfigMap](n.ctx, n.client, res),
	)
}

func (n *Query) Ingress() NamespacedAction[skns.Ingress, net.Ingress] {
	res := skns.Ingress{}
	return NewAction[net.Ingress](
		n.namespace,
		res,
		GetNamespacedAPI[net.Ingress](n.ctx, n.client, res),
	)
}

func (n *Query) HPA() NamespacedAction[skns.HPA, scaling.HorizontalPodAutoscaler] {
	res := skns.HPA{}
	return NewAction[scaling.HorizontalPodAutoscaler](
		n.namespace,
		res,
		GetNamespacedAPI[scaling.HorizontalPodAutoscaler](n.ctx, n.client, res),
	)
}
