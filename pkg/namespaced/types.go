package namespaced

import (
	"github.com/ilexPar/simple-kube/pkg/base"
	"github.com/ilexPar/simple-kube/pkg/namespaced/resources"
)

type NamespacedResourceMethods interface {
	base.ResourceInterface

	API() resources.NamespacedResourceAPI
}

type NamespacedResourcesConstrain interface {
	resources.Deployment | resources.Service | resources.Job | resources.CronJob | resources.ConfigMap | resources.Ingress
}

type NamespacedResources interface {
	NamespacedResourcesConstrain
	NamespacedResourceMethods
}

type QueryNamespace interface {
	Deployment() NamespacedAction[resources.Deployment]
	Service() NamespacedAction[resources.Service]
	Job() NamespacedAction[resources.Job]
	CronJob() NamespacedAction[resources.CronJob]
	ConfigMap() NamespacedAction[resources.ConfigMap]
	Ingress() NamespacedAction[resources.Ingress]
}

type NamespacedAction[T NamespacedResources] interface {
	Get(string) NamespacedGetInterface[T]
	List() NamespacedListInterface[T]
	Create(T) NamespacedPutInterface[T]
	Update(T) NamespacedPutInterface[T]
	Delete(string) NamespacedDeleteInterface[T]
}

type NamespacedGetInterface[T NamespacedResources] interface {
	Run() (T, error)
	DataHandler(func(interface{}) error) NamespacedGetInterface[T]
}

type NamespacedPutInterface[T NamespacedResources] interface {
	Run() error
	DataHandler(func(interface{}) error) NamespacedPutInterface[T]
}

type NamespacedListInterface[T NamespacedResources] interface {
	Run() ([]T, error)
	FilterByLabels(labels map[string]string) NamespacedListInterface[T]
}

type NamespacedDeleteInterface[T NamespacedResources] interface {
	Run() error
}
