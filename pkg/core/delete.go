package core

import "context"

type Delete[T any, R any] struct {
	Action[T, R]
	Id string
	// Ctx overrides the client-configured context for this query; nil falls back
	// to the context the client was configured with.
	Ctx context.Context
}

func (d *Delete[T, R]) WithContext(ctx context.Context) DeleteInterface[T, R] {
	d.Ctx = ctx
	return d
}

func (d *Delete[T, R]) Run() error {
	return d.Backend.Delete(d.Ctx, d.Id)
}
