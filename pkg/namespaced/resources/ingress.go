package resources

import (
	"github.com/ilexPar/simple-kube/pkg/base"

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

type IngressAPI[R base.KubernetesResources] struct {
	base.KubeAPI
}

func (i *IngressAPI[R]) Get(name, namespace string) (*R, error) {
	res, err := i.Client.NetworkingV1().
		Ingresses(namespace).
		Get(i.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	result, err := base.Cast[*R](res)
	return result, err
}

func (i *IngressAPI[R]) Create(namespace string, obj *R) error {
	res, err := base.Cast[*net.Ingress](obj)
	if err != nil {
		return err
	}

	_, err = i.Client.NetworkingV1().
		Ingresses(namespace).
		Create(i.Context, res, metav1.CreateOptions{})
	return err
}

func (i *IngressAPI[R]) Update(namespace string, obj *R) error {
	res, err := base.Cast[*net.Ingress](obj)
	if err != nil {
		return err
	}

	_, err = i.Client.NetworkingV1().
		Ingresses(namespace).
		Update(i.Context, res, metav1.UpdateOptions{})
	return err
}

func (i *IngressAPI[R]) List(namespace string) ([]R, error) {
	res := []R{}
	list, err := i.Client.NetworkingV1().
		Ingresses(namespace).
		List(i.Context, i.Opts.List)
	if err != nil {
		return nil, err
	}

	for _, v := range list.Items {
		item, err := base.Cast[R](v)
		if err == nil {
			res = append(res, item)
		}
	}
	return res, err
}

func (i *IngressAPI[R]) Delete(name, namespace string) error {
	return i.Client.NetworkingV1().
		Ingresses(namespace).
		Delete(i.Context, name, metav1.DeleteOptions{})
}
