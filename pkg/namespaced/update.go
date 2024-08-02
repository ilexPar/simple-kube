package namespaced

type NamespacedUpdate[T NamespacedResources] struct {
	Action[T]
	Resource T
	callback func(interface{}) error
}

func (u *NamespacedUpdate[T]) Run() error {
	obj, err := u.Resource.Dump(u.Resource)
	if err != nil {
		return err
	}

	if u.callback != nil {
		if err = u.callback(obj); err != nil {
			return err
		}
	}

	err = u.api.Update(u.namespace, obj)
	return err
}

func (u *NamespacedUpdate[T]) DataHandler(
	handler func(interface{}) error,
) NamespacedPutInterface[T] {
	u.callback = handler
	return u
}
