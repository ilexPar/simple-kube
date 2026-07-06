package resources

import (
	scaling "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ilexPar/simple-kube/pkg/base"
)

type HPA struct {
	Name    string      `sm:"metadata.name"`
	Min     int         `sm:"spec.minReplicas"`
	Max     int         `sm:"spec.maxReplicas"`
	Target  HPATarget   `sm:"spec.scaleTargetRef"`
	Metrics []HPAMetric `sm:"spec.metrics"`
}

type HPATarget struct {
	Name       string `sm:"name"`
	Kind       string `sm:"kind"`
	APIVersion string `sm:"apiVersion"`
}

type HPAMetric struct {
	Type     scaling.MetricSourceType `sm:"type"`
	Resource HPAResourceMetric        `sm:"resource"`
}

type HPAResourceMetric struct {
	Name        string                   `sm:"name"`
	Type        scaling.MetricTargetType `sm:"target.type"`
	Utilization int                      `sm:"target.averageUtilization"`
}

type HPAapi[R base.KubernetesResources] struct {
	base.KubeAPI
}

func (h *HPAapi[R]) Get(name, namespace string) (*R, error) {
	res, err := h.Client.AutoscalingV2().
		HorizontalPodAutoscalers(namespace).
		Get(h.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	result, err := base.Cast[*R](res)
	return result, err
}

func (h *HPAapi[R]) Create(namespace string, obj *R) error {
	res, err := base.Cast[*scaling.HorizontalPodAutoscaler](obj)
	if err != nil {
		return err
	}

	_, err = h.Client.AutoscalingV2().
		HorizontalPodAutoscalers(namespace).
		Create(h.Context, res, metav1.CreateOptions{})
	return err
}

func (h *HPAapi[R]) Update(namespace string, obj *R) error {
	res, err := base.Cast[*scaling.HorizontalPodAutoscaler](obj)
	if err != nil {
		return err
	}

	_, err = h.Client.AutoscalingV2().
		HorizontalPodAutoscalers(namespace).
		Update(h.Context, res, metav1.UpdateOptions{})
	return err
}

func (h *HPAapi[R]) List(namespace string) ([]R, error) {
	list, err := h.Client.AutoscalingV2().
		HorizontalPodAutoscalers(namespace).
		List(h.Context, h.Opts.List)
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

func (h *HPAapi[R]) Delete(name, namespace string) error {
	return h.Client.AutoscalingV2().
		HorizontalPodAutoscalers(namespace).
		Delete(h.Context, name, metav1.DeleteOptions{})
}
