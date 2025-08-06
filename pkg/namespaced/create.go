package namespaced

import (
	sm "github.com/ilexPar/struct-marshal/pkg"

	"github.com/ilexPar/simple-kube/pkg/base"
)

type NamespacedCreate[T NamespacedResources, R base.KubernetesResources] struct {
	Action[T, R]
	Resource T
	callback func(*R) error
}

func (c *NamespacedCreate[T, R]) Run() error {
	obj := new(R)
	err := sm.Marshal(c.Resource, obj)
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

func (c *NamespacedCreate[T, R]) DataHandler(
	handler func(*R) error,
) NamespacedPutInterface[T, R] {
	c.callback = handler
	return c
}
