package resources

import (
	"github.com/ilexPar/simple-kube/pkg/base"

	sm "github.com/ilexPar/struct-marshal/pkg"
	net "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Ingress struct {
	Name   string            `sm:"metadata.name"`
	Domain string            `sm:"spec.rules[0].host"`
	Paths  []IngressPathDef  `sm:"spec.rules[0].http.paths"`
	Labels map[string]string `sm:"metadata.labels"`
}

type IngressPathDef struct {
	Path    string       `sm:"path"`
	Type    net.PathType `sm:"pathType"`
	Service string       `sm:"backend.service.name"`
	Port    int          `sm:"backend.service.port.number"`
}

func (i Ingress) API() NamespacedResourceAPI {
	return &IngressAPI{}
}

func (i Ingress) Dump(from interface{}) (interface{}, error) {
	res := &net.Ingress{}
	err := sm.Marshal(from, res)
	return res, err
}

func (i Ingress) Load(from, into interface{}) error {
	return sm.Unmarshal(from, into)
}

type IngressAPI struct {
	base.KubeAPI
}

func (i *IngressAPI) Get(name, namespace string) (interface{}, error) {
	res, err := i.Client.NetworkingV1().
		Ingresses(namespace).
		Get(i.Context, name, metav1.GetOptions{})
	return res, err
}

func (i *IngressAPI) Create(namespace string, obj interface{}) error {
	res := obj.(*net.Ingress)
	_, err := i.Client.NetworkingV1().
		Ingresses(namespace).
		Create(i.Context, res, metav1.CreateOptions{})
	return err
}

func (i *IngressAPI) Update(namespace string, obj interface{}) error {
	res := obj.(*net.Ingress)
	_, err := i.Client.NetworkingV1().
		Ingresses(namespace).
		Update(i.Context, res, metav1.UpdateOptions{})
	return err
}

func (i *IngressAPI) List(namespace string) ([]interface{}, error) {
	var res []interface{}
	list, err := i.Client.NetworkingV1().
		Ingresses(namespace).
		List(i.Context, i.Opts.List)
	for _, v := range list.Items {
		res = append(res, v)
	}
	return res, err
}

func (i *IngressAPI) Delete(name, namespace string) error {
	return i.Client.NetworkingV1().
		Ingresses(namespace).
		Delete(i.Context, name, metav1.DeleteOptions{})
}
