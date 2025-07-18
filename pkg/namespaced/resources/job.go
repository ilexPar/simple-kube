package resources

import (
	"github.com/ilexPar/simple-kube/pkg/base"

	batch "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Job struct {
	Name           string            `sm:"metadata.name"`
	ServiceAccount string            `sm:"spec.template.spec.serviceAccountName"`
	Behaviour      JobBehaviour      `sm:"->"`
	Containers     []Container       `sm:"spec.template.spec.containers"`
	Labels         map[string]string `sm:"metadata.labels"`
	NodeSelector   map[string]string `sm:"spec.template.spec.nodeSelector"`
	TemplateLabels map[string]string `sm:"spec.template.metadata.labels"`
}

type JobBehaviour struct {
	RestartPolicy v1.RestartPolicy `sm:"spec.template.spec.restartPolicy"`
	FinishedTTL   int32            `sm:"spec.ttlSecondsAfterFinished"`
}

type JobAPI[R base.KubernetesResources] struct {
	base.KubeAPI
}

func (j *JobAPI[R]) Get(name, namespace string) (*R, error) {
	res, err := j.Client.BatchV1().
		Jobs(namespace).
		Get(j.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	result, err := base.Cast[*R](res)
	return result, err
}

func (j *JobAPI[R]) Create(namespace string, obj *R) error {
	res, err := base.Cast[*batch.Job](obj)
	if err != nil {
		return err
	}

	_, err = j.Client.BatchV1().
		Jobs(namespace).
		Create(j.Context, res, metav1.CreateOptions{})
	return err
}

func (j *JobAPI[R]) Update(namespace string, obj *R) error {
	res, err := base.Cast[*batch.Job](obj)
	if err != nil {
		return err
	}

	_, err = j.Client.BatchV1().
		Jobs(namespace).
		Update(j.Context, res, metav1.UpdateOptions{})
	return err
}

func (j *JobAPI[R]) List(namespace string) ([]R, error) {
	res := []R{}
	list, err := j.Client.BatchV1().
		Jobs(namespace).
		List(j.Context, j.Opts.List)
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

func (j *JobAPI[R]) Delete(name, namespace string) error {
	return j.Client.BatchV1().
		Jobs(namespace).
		Delete(j.Context, name, metav1.DeleteOptions{})
}
