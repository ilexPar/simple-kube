package namespaced

type NamespacedDelete[T NamespacedResources] struct {
	Action[T]
	Id string
}

func (d *NamespacedDelete[T]) Run() error {
	return d.api.Delete(d.Id, d.namespace)
}
