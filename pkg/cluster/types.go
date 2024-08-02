package cluster

import (
	"github.com/ilexPar/simple-kube/pkg/base"
	"github.com/ilexPar/simple-kube/pkg/cluster/resources"
)

type ClusterResourceMethods interface {
	base.ResourceInterface

	API() resources.ClusterResourceAPI
}

type ClusterResourcesConstrain interface {
	resources.Namespace
}

type ClusterResources interface {
	ClusterResourcesConstrain
	ClusterResourceMethods
}

type QueryCluster interface {
	Namespace() ClusterAction[resources.Namespace]
}

type ClusterAction[T ClusterResources] interface {
	Get(string) ClusterGetInterface[T]
	List() ClusterListInterface[T]
	Create(T) ClusterPutInterface[T]
	Update(T) ClusterPutInterface[T]
	Delete(string) ClusterDeleteInterface[T]
}

type ClusterGetInterface[T ClusterResources] interface {
	Run() (T, error)
	DataHandler(func(interface{}) error) ClusterGetInterface[T]
}

type ClusterPutInterface[T ClusterResources] interface {
	Run() error
	DataHandler(func(interface{}) error) ClusterPutInterface[T]
}

type ClusterListInterface[T ClusterResources] interface {
	Run() ([]T, error)
	FilterByLabels(labels map[string]string) ClusterListInterface[T]
}

type ClusterDeleteInterface[T ClusterResources] interface {
	Run() error
}
