package cluster

import (
	api "k8s.io/api/core/v1"

	"github.com/ilexPar/simple-kube/pkg/cluster/resources"
	"github.com/ilexPar/simple-kube/pkg/core"
)

type QueryCluster interface {
	Namespace() core.ScopeAction[resources.Namespace, api.Namespace]
}
