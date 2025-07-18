package namespaced

import (
	"github.com/ilexPar/simple-kube/pkg/base"
	sm "github.com/ilexPar/struct-marshal/pkg"
)

type NamespacedUpdate[T NamespacedResources, R base.KubernetesResources] struct {
	Action[T, R]
	Resource T
	callback func(*R) error
}

func (u *NamespacedUpdate[T, R]) Run() error {
	obj := new(R)
	err := sm.Marshal(u.Resource, obj)
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

func (u *NamespacedUpdate[T, R]) DataHandler(
	handler func(*R) error,
) NamespacedPutInterface[T, R] {
	u.callback = handler
	return u
}
