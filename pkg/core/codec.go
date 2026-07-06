package core

import sm "github.com/ilexPar/struct-marshal/pkg"

// Codec converts between the simplified library type T and the native
// Kubernetes type R. Both scopes (namespaced and cluster) share the same
// struct-marshal based implementation.
type Codec[T any, R any] interface {
	Encode(in T) (*R, error)
	Decode(raw *R, out *T) error
}

// SMCodec bridges T and R using github.com/ilexPar/struct-marshal.
type SMCodec[T any, R any] struct{}

func (SMCodec[T, R]) Encode(in T) (*R, error) {
	out := new(R)
	err := sm.Marshal(in, out)
	return out, err
}

func (SMCodec[T, R]) Decode(raw *R, out *T) error {
	return sm.Unmarshal(raw, out)
}
