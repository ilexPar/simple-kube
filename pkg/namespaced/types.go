package namespaced

import (
	"github.com/ilexPar/simple-kube/pkg/base"
	"github.com/ilexPar/simple-kube/pkg/namespaced/resources"

	apps "k8s.io/api/apps/v1"
	scaling "k8s.io/api/autoscaling/v2"
	batch "k8s.io/api/batch/v1"
	api "k8s.io/api/core/v1"
	net "k8s.io/api/networking/v1"
)

type NamespacedResources interface {
	resources.Deployment | resources.Service | resources.Job | resources.CronJob | resources.ConfigMap | resources.Ingress | resources.HPA
}

type QueryNamespace interface {
	Deployment() NamespacedAction[resources.Deployment, apps.Deployment]
	Service() NamespacedAction[resources.Service, api.Service]
	Job() NamespacedAction[resources.Job, batch.Job]
	CronJob() NamespacedAction[resources.CronJob, batch.CronJob]
	ConfigMap() NamespacedAction[resources.ConfigMap, api.ConfigMap]
	Ingress() NamespacedAction[resources.Ingress, net.Ingress]
	HPA() NamespacedAction[resources.HPA, scaling.HorizontalPodAutoscaler]
}

type NamespacedAction[T NamespacedResources, R base.KubernetesResources] interface {
	Get(string) NamespacedGetInterface[T, R]
	List() NamespacedListInterface[T, R]
	Create(T) NamespacedPutInterface[T, R]
	Update(T) NamespacedPutInterface[T, R]
	Delete(string) NamespacedDeleteInterface[T, R]
}

type NamespacedGetInterface[T NamespacedResources, R base.KubernetesResources] interface {
	Run() (T, error)
	DataHandler(func(*R) error) NamespacedGetInterface[T, R]
}

type NamespacedPutInterface[T NamespacedResources, R base.KubernetesResources] interface {
	Run() error
	DataHandler(func(*R) error) NamespacedPutInterface[T, R]
}

type NamespacedListInterface[T NamespacedResources, R base.KubernetesResources] interface {
	Run() ([]T, error)
	FilterByLabels(labels map[string]string) NamespacedListInterface[T, R]
}

type NamespacedDeleteInterface[T NamespacedResources, R base.KubernetesResources] interface {
	Run() error
}
