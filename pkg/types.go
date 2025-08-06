package simplekube

import (
	"context"

	"k8s.io/client-go/kubernetes"

	"github.com/ilexPar/simple-kube/pkg/cluster"
	"github.com/ilexPar/simple-kube/pkg/namespaced"
)

type ClientInterface interface {
	Config(ctx context.Context, client kubernetes.Interface) ClientInterface

	cluster.QueryCluster
	InNamespace(namespace string) namespaced.QueryNamespace
}
