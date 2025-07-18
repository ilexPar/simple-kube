package resources

import (
	"github.com/ilexPar/simple-kube/pkg/base"

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

type CronJobAPI[R base.KubernetesResources] struct {
	base.KubeAPI
}

func (cj *CronJobAPI[R]) Get(name, namespace string) (*R, error) {
	res, err := cj.Client.BatchV1().
		CronJobs(namespace).
		Get(cj.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	result, err := base.Cast[*R](res)
	return result, err
}

func (cj *CronJobAPI[R]) Create(namespace string, obj *R) error {
	res, err := base.Cast[*batch.CronJob](obj)
	if err != nil {
		return err
	}

	_, err = cj.Client.BatchV1().
		CronJobs(namespace).
		Create(cj.Context, res, metav1.CreateOptions{})
	return err
}

func (cj *CronJobAPI[R]) Update(namespace string, obj *R) error {
	res, err := base.Cast[*batch.CronJob](obj)
	if err != nil {
		return err
	}

	_, err = cj.Client.BatchV1().
		CronJobs(namespace).
		Update(cj.Context, res, metav1.UpdateOptions{})
	return err
}

func (cj *CronJobAPI[R]) List(namespace string) ([]R, error) {
	res := []R{}
	list, err := cj.Client.BatchV1().
		CronJobs(namespace).
		List(cj.Context, cj.Opts.List)
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

func (cj *CronJobAPI[R]) Delete(name, namespace string) error {
	return cj.Client.BatchV1().
		CronJobs(namespace).
		Delete(cj.Context, name, metav1.DeleteOptions{})
}
