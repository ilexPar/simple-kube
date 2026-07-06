package core

import (
	"context"

	"github.com/ilexPar/simple-kube/pkg/base"
)

type List[T any, R any] struct {
	Action[T, R]
	Opts base.QueryOpts
	// Ctx overrides the client-configured context for this query; nil falls back
	// to the context the client was configured with.
	Ctx context.Context
}

// FilterByLabels mutates only this List's own Opts field, so it is safe to build
// queries independently.
func (l *List[T, R]) FilterByLabels(labels map[string]string) ListInterface[T, R] {
	l.Opts.List.LabelSelector = base.FlattenLabels(labels)
	return l
}

func (l *List[T, R]) WithContext(ctx context.Context) ListInterface[T, R] {
	l.Ctx = ctx
	return l
}

func (l *List[T, R]) Run() ([]T, error) {
	res := []T{}

	objs, err := l.Backend.List(l.Ctx, l.Opts)
	if err != nil {
		return res, err
	}

	for i := range objs {
		var out T
		if err = l.Codec.Decode(&objs[i], &out); err != nil {
			return res, err
		}
		res = append(res, out)
	}

	return res, nil
}
