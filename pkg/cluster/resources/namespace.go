package resources

import (
	"github.com/ilexPar/simple-kube/pkg/base"

	sm "github.com/ilexPar/struct-marshal/pkg"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Namespace struct {
	Name   string            `sm:"metadata.name"`
	Labels map[string]string `sm:"metadata.labels"`
}

func (n Namespace) API() ClusterResourceAPI {
	return &NamespaceAPI{}
}

func (n Namespace) Dump(from interface{}) (interface{}, error) {
	res := &api.Namespace{}
	err := sm.Marshal(from, res)
	return res, err
}

func (n Namespace) Load(from, into interface{}) error {
	return sm.Unmarshal(from, into)
}

type NamespaceAPI struct {
	base.KubeAPI
}

func (n *NamespaceAPI) Get(name string) (interface{}, error) {
	res, err := n.Client.CoreV1().
		Namespaces().
		Get(n.Context, name, metav1.GetOptions{})
	return res, err
}

func (n *NamespaceAPI) Create(obj interface{}) error {
	res := obj.(*api.Namespace)
	_, err := n.Client.CoreV1().
		Namespaces().
		Create(n.Context, res, metav1.CreateOptions{})
	return err
}

func (n *NamespaceAPI) Update(obj interface{}) error {
	res := obj.(*api.Namespace)
	_, err := n.Client.CoreV1().
		Namespaces().
		Update(n.Context, res, metav1.UpdateOptions{})
	return err
}

func (n *NamespaceAPI) List() ([]interface{}, error) {
	var res []interface{}
	list, err := n.Client.CoreV1().
		Namespaces().
		List(n.Context, n.Opts.List)
	for _, v := range list.Items {
		res = append(res, v)
	}
	return res, err
}

func (n *NamespaceAPI) Delete(name string) error {
	return n.Client.CoreV1().
		Namespaces().
		Delete(n.Context, name, metav1.DeleteOptions{})
}
