package core

// ScopeAction is the entrypoint returned by a scope (e.g. Deployment(),
// Namespace()) exposing the available actions for a resource.
type ScopeAction[T any, R any] interface {
	Get(name string) GetInterface[T, R]
	List() ListInterface[T, R]
	Create(resource T) PutInterface[T, R]
	Update(resource T) PutInterface[T, R]
	Delete(name string) DeleteInterface[T, R]
}

type GetInterface[T any, R any] interface {
	Run() (T, error)
	DataHandler(func(*R) error) GetInterface[T, R]
}

type PutInterface[T any, R any] interface {
	Run() error
	DataHandler(func(*R) error) PutInterface[T, R]
}

type ListInterface[T any, R any] interface {
	Run() ([]T, error)
	FilterByLabels(labels map[string]string) ListInterface[T, R]
}

type DeleteInterface[T any, R any] interface {
	Run() error
}
