package base

import (
	apps "k8s.io/api/apps/v1"
	scaling "k8s.io/api/autoscaling/v2"
	batch "k8s.io/api/batch/v1"
	api "k8s.io/api/core/v1"
	net "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type QueryOpts struct {
	List metav1.ListOptions
}

type ResourceInterface interface {
	Load(from, into interface{}) error
	Dump(from interface{}) (interface{}, error)
}

type KubernetesResources interface {
	apps.Deployment | api.ConfigMap | api.Service | batch.Job | batch.CronJob | net.Ingress | scaling.HorizontalPodAutoscaler
}
