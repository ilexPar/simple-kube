package cluster

type ClusterCreate[T ClusterResources] struct {
	Action[T]
	Resource T
	callback func(interface{}) error
}

func (c *ClusterCreate[T]) Run() error {
	obj, err := c.Resource.Dump(c.Resource)
	if err != nil {
		return err
	}

	if c.callback != nil {
		if err = c.callback(obj); err != nil {
			return err
		}
	}

	err = c.api.Create(obj)
	return err
}

func (c *ClusterCreate[T]) DataHandler(
	handler func(interface{}) error,
) ClusterPutInterface[T] {
	c.callback = handler
	return c
}
