package namespaced

type NamespacedCreate[T NamespacedResources] struct {
	Action[T]
	Resource T
	callback func(interface{}) error
}

func (c *NamespacedCreate[T]) Run() error {
	obj, err := c.Resource.Dump(c.Resource)
	if err != nil {
		return err
	}

	if c.callback != nil {
		if err = c.callback(obj); err != nil {
			return err
		}
	}

	err = c.api.Create(c.namespace, obj)
	return err
}

func (c *NamespacedCreate[T]) DataHandler(
	handler func(interface{}) error,
) NamespacedPutInterface[T] {
	c.callback = handler
	return c
}
