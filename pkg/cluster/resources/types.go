package resources

import (
	"context"

	"k8s.io/client-go/kubernetes"

	"github.com/ilexPar/simple-kube/pkg/base"
)

type ClusterResourceAPI interface {
	Config(ctx context.Context, k8s kubernetes.Interface)
	SetOpts(opts base.QueryOpts)

	Get(name string) (interface{}, error)
	Create(obj interface{}) error
	Update(obj interface{}) error
	List() ([]interface{}, error)
	Delete(name string) error
}
