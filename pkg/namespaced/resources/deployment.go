package resources

import (
	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ilexPar/simple-kube/pkg/base"
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

type DeploymentAPI[R base.KubernetesResources] struct {
	base.KubeAPI
}

func (d *DeploymentAPI[R]) Get(name, namespace string) (*R, error) {
	res, err := d.Client.AppsV1().
		Deployments(namespace).
		Get(d.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	result, err := base.Cast[*R](res)
	return result, err
}

func (d *DeploymentAPI[R]) Create(namespace string, obj *R) error {
	res, err := base.Cast[*apps.Deployment](obj)
	if err != nil {
		return err
	}

	_, err = d.Client.AppsV1().
		Deployments(namespace).
		Create(d.Context, res, metav1.CreateOptions{})
	return err
}

func (d *DeploymentAPI[R]) Update(namespace string, obj *R) error {
	res, err := base.Cast[*apps.Deployment](obj)
	if err != nil {
		return err
	}

	_, err = d.Client.AppsV1().
		Deployments(namespace).
		Update(d.Context, res, metav1.UpdateOptions{})
	return err
}

func (d *DeploymentAPI[R]) List(namespace string) ([]R, error) {
	list, err := d.Client.AppsV1().
		Deployments(namespace).
		List(d.Context, d.Opts.List)
	if err != nil {
		return nil, err
	}

	res := make([]R, 0, len(list.Items))
	for _, v := range list.Items {
		item, err := base.Cast[R](v)
		if err == nil {
			res = append(res, item)
		}
	}
	return res, err
}

func (d *DeploymentAPI[R]) Delete(name, namespace string) error {
	return d.Client.AppsV1().
		Deployments(namespace).
		Delete(d.Context, name, metav1.DeleteOptions{})
}
