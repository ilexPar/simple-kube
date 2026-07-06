package namespaced

import (
	apps "k8s.io/api/apps/v1"
	scaling "k8s.io/api/autoscaling/v2"
	batch "k8s.io/api/batch/v1"
	api "k8s.io/api/core/v1"
	net "k8s.io/api/networking/v1"

	"github.com/ilexPar/simple-kube/pkg/core"
	"github.com/ilexPar/simple-kube/pkg/namespaced/resources"
)

type QueryNamespace interface {
	Deployment() core.ScopeAction[resources.Deployment, apps.Deployment]
	Service() core.ScopeAction[resources.Service, api.Service]
	Job() core.ScopeAction[resources.Job, batch.Job]
	CronJob() core.ScopeAction[resources.CronJob, batch.CronJob]
	ConfigMap() core.ScopeAction[resources.ConfigMap, api.ConfigMap]
	Ingress() core.ScopeAction[resources.Ingress, net.Ingress]
	HPA() core.ScopeAction[resources.HPA, scaling.HorizontalPodAutoscaler]
}
