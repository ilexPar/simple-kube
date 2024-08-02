package resources

import (
	"github.com/ilexPar/simple-kube/pkg/base"

	sm "github.com/ilexPar/struct-marshal/pkg"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Service struct {
	Name     string            `sm:"metadata.name"`
	Port     int               `sm:"spec.ports[0].port"`
	Selector map[string]string `sm:"spec.selector"`
	Labels   map[string]string `sm:"metadata.labels"`
}

func (s Service) API() NamespacedResourceAPI {
	return &ServiceAPI{}
}

func (s Service) Dump(from interface{}) (interface{}, error) {
	res := &api.Service{}
	err := sm.Marshal(from, res)
	return res, err
}

func (s Service) Load(from, into interface{}) error {
	return sm.Unmarshal(from, into)
}

type ServiceAPI struct {
	base.KubeAPI
}

func (s *ServiceAPI) Get(name, namespace string) (interface{}, error) {
	res, err := s.Client.CoreV1().
		Services(namespace).
		Get(s.Context, name, metav1.GetOptions{})
	return res, err
}

func (s *ServiceAPI) Create(namespace string, obj interface{}) error {
	res := obj.(*api.Service)
	_, err := s.Client.CoreV1().
		Services(namespace).
		Create(s.Context, res, metav1.CreateOptions{})
	return err
}

func (s *ServiceAPI) Update(namespace string, obj interface{}) error {
	res := obj.(*api.Service)
	_, err := s.Client.CoreV1().
		Services(namespace).
		Update(s.Context, res, metav1.UpdateOptions{})
	return err
}

func (s *ServiceAPI) List(namespace string) ([]interface{}, error) {
	var res []interface{}
	list, err := s.Client.CoreV1().
		Services(namespace).
		List(s.Context, s.Opts.List)
	for _, v := range list.Items {
		res = append(res, v)
	}
	return res, err
}

func (s *ServiceAPI) Delete(name, namespace string) error {
	return s.Client.CoreV1().
		Services(namespace).
		Delete(s.Context, name, metav1.DeleteOptions{})
}
