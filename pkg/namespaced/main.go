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
	"github.com/ilexPar/simple-kube/pkg/core"
	"github.com/ilexPar/simple-kube/pkg/namespaced/resources"
)

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

// newAction wires a namespaced resource API into the shared generic core
// action: the concrete API type is known by each Query method, so no runtime
// type dispatch is needed.
func newAction[R base.KubernetesResources, T any](
	q *Query,
	api resources.NamespacedResourceAPI[R],
) core.ScopeAction[T, R] {
	api.Config(q.ctx, q.client)
	backend := nsBackend[R]{api: api, namespace: q.namespace}
	return core.NewAction(core.SMCodec[T, R]{}, backend)
}

func (n *Query) Deployment() core.ScopeAction[resources.Deployment, apps.Deployment] {
	return newAction[apps.Deployment, resources.Deployment](n, &resources.DeploymentAPI[apps.Deployment]{})
}

func (n *Query) Service() core.ScopeAction[resources.Service, api.Service] {
	return newAction[api.Service, resources.Service](n, &resources.ServiceAPI[api.Service]{})
}

func (n *Query) Job() core.ScopeAction[resources.Job, batch.Job] {
	return newAction[batch.Job, resources.Job](n, &resources.JobAPI[batch.Job]{})
}

func (n *Query) CronJob() core.ScopeAction[resources.CronJob, batch.CronJob] {
	return newAction[batch.CronJob, resources.CronJob](n, &resources.CronJobAPI[batch.CronJob]{})
}

func (n *Query) ConfigMap() core.ScopeAction[resources.ConfigMap, api.ConfigMap] {
	return newAction[api.ConfigMap, resources.ConfigMap](n, &resources.ConfigMapAPI[api.ConfigMap]{})
}

func (n *Query) Ingress() core.ScopeAction[resources.Ingress, net.Ingress] {
	return newAction[net.Ingress, resources.Ingress](n, &resources.IngressAPI[net.Ingress]{})
}

func (n *Query) HPA() core.ScopeAction[resources.HPA, scaling.HorizontalPodAutoscaler] {
	return newAction[scaling.HorizontalPodAutoscaler, resources.HPA](
		n,
		&resources.HPAapi[scaling.HorizontalPodAutoscaler]{},
	)
}
