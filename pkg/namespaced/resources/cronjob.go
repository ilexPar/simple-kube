package resources

import (
	"github.com/ilexPar/simple-kube/pkg/base"

	sm "github.com/ilexPar/struct-marshal/pkg"
	batch "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CronJob struct {
	Name           string            `sm:"metadata.name"`
	Schedule       string            `sm:"spec.schedule"`
	Behaviour      CronJobBehaviour  `sm:"->"`
	ServiceAccount string            `sm:"spec.jobTemplate.spec.template.spec.serviceAccountName"`
	Containers     []Container       `sm:"spec.jobTemplate.spec.template.spec.containers"`
	Labels         map[string]string `sm:"metadata.labels"`
	NodeSelector   map[string]string `sm:"spec.jobTemplate.spec.template.spec.nodeSelector"`
	TemplateLabels map[string]string `sm:"spec.jobTemplate.spec.template.metadata.labels"`
}

type CronJobBehaviour struct {
	RestartPolicy    v1.RestartPolicy `sm:"spec.jobTemplate.spec.template.spec.restartPolicy"`
	SuccessHistory   int32            `sm:"spec.successfulJobsHistoryLimit"`
	FailedHistory    int32            `sm:"spec.failedJobsHistoryLimit"`
	StartingDeadline int64            `sm:"spec.startingDeadlineSeconds"`
}

func (cj CronJob) API() NamespacedResourceAPI {
	return &CronJobAPI{}
}

func (cj CronJob) Dump(from interface{}) (interface{}, error) {
	res := &batch.CronJob{}
	err := sm.Marshal(from, res)
	return res, err
}

func (cj CronJob) Load(from, into interface{}) error {
	return sm.Unmarshal(from, into)
}

type CronJobAPI struct {
	base.KubeAPI
}

func (cj *CronJobAPI) Get(name, namespace string) (interface{}, error) {
	res, err := cj.Client.BatchV1().
		CronJobs(namespace).
		Get(cj.Context, name, metav1.GetOptions{})
	return res, err
}

func (cj *CronJobAPI) Create(namespace string, obj interface{}) error {
	res := obj.(*batch.CronJob)
	_, err := cj.Client.BatchV1().
		CronJobs(namespace).
		Create(cj.Context, res, metav1.CreateOptions{})
	return err
}

func (cj *CronJobAPI) Update(namespace string, obj interface{}) error {
	res := obj.(*batch.CronJob)
	_, err := cj.Client.BatchV1().
		CronJobs(namespace).
		Update(cj.Context, res, metav1.UpdateOptions{})
	return err
}

func (cj *CronJobAPI) List(namespace string) ([]interface{}, error) {
	var res []interface{}
	list, err := cj.Client.BatchV1().
		CronJobs(namespace).
		List(cj.Context, cj.Opts.List)
	for _, v := range list.Items {
		res = append(res, v)
	}
	return res, err
}

func (cj *CronJobAPI) Delete(name, namespace string) error {
	return cj.Client.BatchV1().
		CronJobs(namespace).
		Delete(cj.Context, name, metav1.DeleteOptions{})
}
