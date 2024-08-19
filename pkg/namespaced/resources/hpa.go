package resources

import (
	"github.com/ilexPar/simple-kube/pkg/base"
	sm "github.com/ilexPar/struct-marshal/pkg"
	scaling "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (h HPA) API() NamespacedResourceAPI {
	return &HPAapi{}
}

func (h HPA) Dump(from interface{}) (interface{}, error) {
	res := &scaling.HorizontalPodAutoscaler{}
	err := sm.Marshal(from, res)
	return res, err
}

func (h HPA) Load(from, into interface{}) error {
	return sm.Unmarshal(from, into)
}

type HPAapi struct {
	base.KubeAPI
}

func (h *HPAapi) Get(name, namespace string) (interface{}, error) {
	res, err := h.Client.AutoscalingV2().
		HorizontalPodAutoscalers(namespace).
		Get(h.Context, name, metav1.GetOptions{})
	return res, err
}

func (h *HPAapi) Create(namespace string, obj interface{}) error {
	res := obj.(*scaling.HorizontalPodAutoscaler)
	_, err := h.Client.AutoscalingV2().
		HorizontalPodAutoscalers(namespace).
		Create(h.Context, res, metav1.CreateOptions{})
	return err
}

func (h *HPAapi) Update(namespace string, obj interface{}) error {
	res := obj.(*scaling.HorizontalPodAutoscaler)
	_, err := h.Client.AutoscalingV2().
		HorizontalPodAutoscalers(namespace).
		Update(h.Context, res, metav1.UpdateOptions{})
	return err
}

func (h *HPAapi) List(namespace string) ([]interface{}, error) {
	var res []interface{}
	list, err := h.Client.AutoscalingV2().
		HorizontalPodAutoscalers(namespace).
		List(h.Context, h.Opts.List)
	for _, v := range list.Items {
		res = append(res, v)
	}
	return res, err
}

func (h *HPAapi) Delete(name, namespace string) error {
	return h.Client.AutoscalingV2().
		HorizontalPodAutoscalers(namespace).
		Delete(h.Context, name, metav1.DeleteOptions{})
}
