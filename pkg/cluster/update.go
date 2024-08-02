package cluster

type ClusterUpdate[T ClusterResources] struct {
	Action[T]
	Resource T
	callback func(interface{}) error
}

func (u *ClusterUpdate[T]) Run() error {
	obj, err := u.Resource.Dump(u.Resource)
	if err != nil {
		return err
	}

	if u.callback != nil {
		if err = u.callback(obj); err != nil {
			return err
		}
	}

	err = u.api.Update(obj)
	return err
}

func (u *ClusterUpdate[T]) DataHandler(
	handler func(interface{}) error,
) ClusterPutInterface[T] {
	u.callback = handler
	return u
}
