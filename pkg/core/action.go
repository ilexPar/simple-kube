package core

// Action is the shared, fully typed implementation of ScopeAction reused by
// both the namespaced and cluster scopes. T is the simplified library type and
// R is the native Kubernetes type; a Codec bridges them and a Backend performs
// the API calls. It holds only the immutable wiring; per-query state (context,
// list options) lives on the individual action leaves.
type Action[T any, R any] struct {
	Codec   Codec[T, R]
	Backend Backend[R]
}

func NewAction[T any, R any](codec Codec[T, R], backend Backend[R]) *Action[T, R] {
	return &Action[T, R]{
		Codec:   codec,
		Backend: backend,
	}
}

func (a *Action[T, R]) Get(name string) GetInterface[T, R] {
	return &Get[T, R]{Action: *a, Id: name}
}

func (a *Action[T, R]) Create(resource T) PutInterface[T, R] {
	return &Put[T, R]{Action: *a, Resource: resource, send: a.Backend.Create}
}

func (a *Action[T, R]) Update(resource T) PutInterface[T, R] {
	return &Put[T, R]{Action: *a, Resource: resource, send: a.Backend.Update}
}

func (a *Action[T, R]) List() ListInterface[T, R] {
	return &List[T, R]{Action: *a}
}

func (a *Action[T, R]) Delete(name string) DeleteInterface[T, R] {
	return &Delete[T, R]{Action: *a, Id: name}
}
