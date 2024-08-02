package simplekube

import (
	"github.com/ilexPar/simple-kube/pkg/cluster"
	"github.com/ilexPar/simple-kube/pkg/namespaced"
)

type ClientInterface interface {
	NamespacedQuery(namespace string) namespaced.QueryNamespace
	ClusterQuery() cluster.QueryCluster
}
