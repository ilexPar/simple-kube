package resources

import (
	"github.com/ilexPar/simple-kube/pkg/base"

	sm "github.com/ilexPar/struct-marshal/pkg"
	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Deployment struct {
	Name            string            `sm:"metadata.name"`
	ServiceAccount  string            `sm:"spec.template.spec.serviceAccountName"`
	Containers      []Container       `sm:"spec.template.spec.containers"`
	Labels          map[string]string `sm:"metadata.labels"`
	TemplateLabels  map[string]string `sm:"spec.template.metadata.labels"`
	ServiceSelector map[string]string `sm:"spec.selector.matchLabels"`
	NodeSelector    map[string]string `sm:"spec.template.spec.nodeSelector"`
}

func (d Deployment) API() NamespacedResourceAPI {
	return &DeploymentAPI{}
}

func (d Deployment) Dump(from interface{}) (interface{}, error) {
	res := &apps.Deployment{}
	err := sm.Marshal(from, res)
	return res, err
}

func (d Deployment) Load(from, into interface{}) error {
	return sm.Unmarshal(from, into)
}

type DeploymentAPI struct {
	base.KubeAPI
}

func (d *DeploymentAPI) Get(name, namespace string) (interface{}, error) {
	res, err := d.Client.AppsV1().
		Deployments(namespace).
		Get(d.Context, name, metav1.GetOptions{})
	return res, err
}

func (d *DeploymentAPI) Create(namespace string, obj interface{}) error {
	res := obj.(*apps.Deployment)
	_, err := d.Client.AppsV1().
		Deployments(namespace).
		Create(d.Context, res, metav1.CreateOptions{})
	return err
}

func (d *DeploymentAPI) Update(namespace string, obj interface{}) error {
	res := obj.(*apps.Deployment)
	_, err := d.Client.AppsV1().
		Deployments(namespace).
		Update(d.Context, res, metav1.UpdateOptions{})
	return err
}

func (d *DeploymentAPI) List(namespace string) ([]interface{}, error) {
	var res []interface{}
	list, err := d.Client.AppsV1().
		Deployments(namespace).
		List(d.Context, d.Opts.List)
	for _, v := range list.Items {
		res = append(res, v)
	}
	return res, err
}

func (d *DeploymentAPI) Delete(name, namespace string) error {
	return d.Client.AppsV1().
		Deployments(namespace).
		Delete(d.Context, name, metav1.DeleteOptions{})
}
