package resources

import (
	"github.com/ilexPar/simple-kube/pkg/base"

	sm "github.com/ilexPar/struct-marshal/pkg"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMap struct {
	Name   string            `sm:"metadata.name"`
	Labels map[string]string `sm:"metadata.labels"`
	Data   map[string]string `sm:"data"`
}

func (cm ConfigMap) API() NamespacedResourceAPI {
	return &ConfigMapAPI{}
}

func (cm ConfigMap) Dump(from interface{}) (interface{}, error) {
	res := &api.ConfigMap{}
	err := sm.Marshal(from, res)
	return res, err
}

func (cm ConfigMap) Load(from, into interface{}) error {
	return sm.Unmarshal(from, into)
}

type ConfigMapAPI struct {
	base.KubeAPI
}

func (cm *ConfigMapAPI) Get(name, namespace string) (interface{}, error) {
	res, err := cm.Client.CoreV1().
		ConfigMaps(namespace).
		Get(cm.Context, name, metav1.GetOptions{})
	return res, err
}

func (cm *ConfigMapAPI) Create(namespace string, obj interface{}) error {
	res := obj.(*api.ConfigMap)
	_, err := cm.Client.CoreV1().
		ConfigMaps(namespace).
		Create(cm.Context, res, metav1.CreateOptions{})
	return err
}

func (cm *ConfigMapAPI) Update(namespace string, obj interface{}) error {
	res := obj.(*api.ConfigMap)
	_, err := cm.Client.CoreV1().
		ConfigMaps(namespace).
		Update(cm.Context, res, metav1.UpdateOptions{})
	return err
}

func (cm *ConfigMapAPI) List(namespace string) ([]interface{}, error) {
	var res []interface{}
	list, err := cm.Client.CoreV1().
		ConfigMaps(namespace).
		List(cm.Context, cm.Opts.List)
	for _, v := range list.Items {
		res = append(res, v)
	}
	return res, err
}

func (cm *ConfigMapAPI) Delete(name, namespace string) error {
	return cm.Client.CoreV1().
		ConfigMaps(namespace).
		Delete(cm.Context, name, metav1.DeleteOptions{})
}
