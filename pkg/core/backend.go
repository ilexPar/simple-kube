package core

import (
	"github.com/ilexPar/simple-kube/pkg/base"
)

// Backend performs the actual Kubernetes calls for a single resource type R.
// It is the only place where the namespaced vs cluster difference lives: the
// namespaced adapter closes over its namespace, the cluster adapter does not.
type Backend[R any] interface {
	Get(id string) (*R, error)
	List(opts base.QueryOpts) ([]R, error)
	Create(obj *R) error
	Update(obj *R) error
	Delete(id string) error
}
