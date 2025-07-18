package resources

import (
	"github.com/ilexPar/simple-kube/pkg/base"

	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMap struct {
	Name   string            `sm:"metadata.name"`
	Labels map[string]string `sm:"metadata.labels"`
	Data   map[string]string `sm:"data"`
}

type ConfigMapAPI[R base.KubernetesResources] struct {
	base.KubeAPI
}

func (cm *ConfigMapAPI[R]) Get(name, namespace string) (*R, error) {
	res, err := cm.Client.CoreV1().
		ConfigMaps(namespace).
		Get(cm.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	result, err := base.Cast[*R](res)
	return result, err
}

func (cm *ConfigMapAPI[R]) Create(namespace string, obj *R) error {
	res, err := base.Cast[*api.ConfigMap](obj)
	if err != nil {
		return err
	}

	_, err = cm.Client.CoreV1().
		ConfigMaps(namespace).
		Create(cm.Context, res, metav1.CreateOptions{})
	return err
}

func (cm *ConfigMapAPI[R]) Update(namespace string, obj *R) error {
	res, err := base.Cast[*api.ConfigMap](obj)
	if err != nil {
		return err
	}

	_, err = cm.Client.CoreV1().
		ConfigMaps(namespace).
		Update(cm.Context, res, metav1.UpdateOptions{})
	return err
}

func (cm *ConfigMapAPI[R]) List(namespace string) ([]R, error) {
	res := []R{}
	list, err := cm.Client.CoreV1().
		ConfigMaps(namespace).
		List(cm.Context, cm.Opts.List)
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

func (cm *ConfigMapAPI[R]) Delete(name, namespace string) error {
	return cm.Client.CoreV1().
		ConfigMaps(namespace).
		Delete(cm.Context, name, metav1.DeleteOptions{})
}
