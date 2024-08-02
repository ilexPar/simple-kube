package resources

import (
	"github.com/ilexPar/simple-kube/pkg/base"

	sm "github.com/ilexPar/struct-marshal/pkg"
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

func (j Job) API() NamespacedResourceAPI {
	return &JobAPI{}
}

func (j Job) Dump(from interface{}) (interface{}, error) {
	res := &batch.Job{}
	err := sm.Marshal(from, res)
	return res, err
}

func (j Job) Load(from, into interface{}) error {
	return sm.Unmarshal(from, into)
}

type JobAPI struct {
	base.KubeAPI
}

func (j *JobAPI) Get(name, namespace string) (interface{}, error) {
	res, err := j.Client.BatchV1().
		Jobs(namespace).
		Get(j.Context, name, metav1.GetOptions{})
	return res, err
}

func (j *JobAPI) Create(namespace string, obj interface{}) error {
	res := obj.(*batch.Job)
	_, err := j.Client.BatchV1().
		Jobs(namespace).
		Create(j.Context, res, metav1.CreateOptions{})
	return err
}

func (j *JobAPI) Update(namespace string, obj interface{}) error {
	res := obj.(*batch.Job)
	_, err := j.Client.BatchV1().
		Jobs(namespace).
		Update(j.Context, res, metav1.UpdateOptions{})
	return err
}

func (j *JobAPI) List(namespace string) ([]interface{}, error) {
	var res []interface{}
	list, err := j.Client.BatchV1().
		Jobs(namespace).
		List(j.Context, j.Opts.List)
	for _, v := range list.Items {
		res = append(res, v)
	}
	return res, err
}

func (j *JobAPI) Delete(name, namespace string) error {
	return j.Client.BatchV1().
		Jobs(namespace).
		Delete(j.Context, name, metav1.DeleteOptions{})
}
