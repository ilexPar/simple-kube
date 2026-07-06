package resources

import (
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ilexPar/simple-kube/pkg/base"
)

type Namespace struct {
	Name   string            `sm:"metadata.name"`
	Labels map[string]string `sm:"metadata.labels"`
}

type NamespaceAPI[R base.KubernetesResources] struct {
	base.KubeAPI
}

func (n *NamespaceAPI[R]) Get(name string) (*R, error) {
	res, err := n.Client.CoreV1().
		Namespaces().
		Get(n.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return base.Cast[*R](res)
}

func (n *NamespaceAPI[R]) Create(obj *R) error {
	res, err := base.Cast[*api.Namespace](obj)
	if err != nil {
		return err
	}
	_, err = n.Client.CoreV1().
		Namespaces().
		Create(n.Context, res, metav1.CreateOptions{})
	return err
}

func (n *NamespaceAPI[R]) Update(obj *R) error {
	res, err := base.Cast[*api.Namespace](obj)
	if err != nil {
		return err
	}
	_, err = n.Client.CoreV1().
		Namespaces().
		Update(n.Context, res, metav1.UpdateOptions{})
	return err
}

func (n *NamespaceAPI[R]) List() ([]R, error) {
	list, err := n.Client.CoreV1().
		Namespaces().
		List(n.Context, n.Opts.List)
	if err != nil {
		return nil, err
	}

	res := make([]R, 0, len(list.Items))
	for i := range list.Items {
		item, err := base.Cast[R](list.Items[i])
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func (n *NamespaceAPI[R]) Delete(name string) error {
	return n.Client.CoreV1().
		Namespaces().
		Delete(n.Context, name, metav1.DeleteOptions{})
}
