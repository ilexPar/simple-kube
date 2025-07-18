package resources

import (
	"github.com/ilexPar/simple-kube/pkg/base"

	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Service struct {
	Name     string            `sm:"metadata.name"`
	Port     int               `sm:"spec.ports[0].port"`
	Selector map[string]string `sm:"spec.selector"`
	Labels   map[string]string `sm:"metadata.labels"`
}

type ServiceAPI[R base.KubernetesResources] struct {
	base.KubeAPI
}

func (s *ServiceAPI[R]) Get(name, namespace string) (*R, error) {
	res, err := s.Client.CoreV1().
		Services(namespace).
		Get(s.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	result, err := base.Cast[*R](res)
	return result, err
}

func (s *ServiceAPI[R]) Create(namespace string, obj *R) error {
	res, err := base.Cast[*api.Service](obj)
	if err != nil {
		return err
	}

	_, err = s.Client.CoreV1().
		Services(namespace).
		Create(s.Context, res, metav1.CreateOptions{})
	return err
}

func (s *ServiceAPI[R]) Update(namespace string, obj *R) error {
	res, err := base.Cast[*api.Service](obj)
	if err != nil {
		return err
	}

	_, err = s.Client.CoreV1().
		Services(namespace).
		Update(s.Context, res, metav1.UpdateOptions{})
	return err
}

func (s *ServiceAPI[R]) List(namespace string) ([]R, error) {
	res := []R{}
	list, err := s.Client.CoreV1().
		Services(namespace).
		List(s.Context, s.Opts.List)
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

func (s *ServiceAPI[R]) Delete(name, namespace string) error {
	return s.Client.CoreV1().
		Services(namespace).
		Delete(s.Context, name, metav1.DeleteOptions{})
}
