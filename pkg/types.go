package simplekube

import (
	"context"

	"github.com/ilexPar/simple-kube/pkg/cluster"
	"github.com/ilexPar/simple-kube/pkg/namespaced"
	"k8s.io/client-go/kubernetes"
)

type ClientInterface interface {
	Config(ctx context.Context, client kubernetes.Interface) ClientInterface

	cluster.QueryCluster
	InNamespace(namespace string) namespaced.QueryNamespace
}
