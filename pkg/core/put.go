package core

import "context"

// Put implements both Create and Update; the concrete Backend call is bound to
// the send field at construction time.
type Put[T any, R any] struct {
	Action[T, R]
	Resource T
	// Ctx overrides the client-configured context for this query; nil falls back
	// to the context the client was configured with.
	Ctx      context.Context
	send     func(context.Context, *R) error
	callback func(*R) error
}

func (p *Put[T, R]) DataHandler(handler func(*R) error) PutInterface[T, R] {
	p.callback = handler
	return p
}

func (p *Put[T, R]) WithContext(ctx context.Context) PutInterface[T, R] {
	p.Ctx = ctx
	return p
}

func (p *Put[T, R]) Run() error {
	raw, err := p.Codec.Encode(p.Resource)
	if err != nil {
		return err
	}

	if p.callback != nil {
		if err = p.callback(raw); err != nil {
			return err
		}
	}

	return p.send(p.Ctx, raw)
}
