package core

import (
	"context"

	"github.com/ilexPar/simple-kube/pkg/base"
)

// Backend performs the actual Kubernetes calls for a single resource type R.
// It is the only place where the namespaced vs cluster difference lives: the
// namespaced adapter closes over its namespace, the cluster adapter does not.
// The ctx passed to each method may be nil, in which case the underlying API
// falls back to the client-configured context.
type Backend[R any] interface {
	Get(ctx context.Context, id string) (*R, error)
	List(ctx context.Context, opts base.QueryOpts) ([]R, error)
	Create(ctx context.Context, obj *R) error
	Update(ctx context.Context, obj *R) error
	Delete(ctx context.Context, id string) error
}
