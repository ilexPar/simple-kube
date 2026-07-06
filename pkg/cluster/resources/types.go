package resources

import (
	"context"

	"k8s.io/client-go/kubernetes"

	"github.com/ilexPar/simple-kube/pkg/base"
)

type ClusterResourceAPI[R base.KubernetesResources] interface {
	Config(ctx context.Context, k8s kubernetes.Interface)
	SetOpts(opts base.QueryOpts)
	SetContext(ctx context.Context)

	Get(name string) (*R, error)
	Create(obj *R) error
	Update(obj *R) error
	List() ([]R, error)
	Delete(name string) error
}
