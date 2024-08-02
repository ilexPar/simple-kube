package resources

import (
	"context"

	"github.com/ilexPar/simple-kube/pkg/base"

	"k8s.io/client-go/kubernetes"
)

type NamespacedResourceAPI interface {
	Config(ctx context.Context, k8s kubernetes.Interface)
	SetOpts(opts base.QueryOpts)

	Get(name, namespace string) (interface{}, error)
	Create(namespace string, obj interface{}) error
	Update(namespace string, obj interface{}) error
	List(namespace string) ([]interface{}, error)
	Delete(name, namespace string) error
}

type Container struct {
	Name      string              `sm:"name"`
	Image     string              `sm:"image"`
	Ports     []ContainerPort     `sm:"ports"`
	Command   []string            `sm:"command"`
	Resources *ContainerResources `sm:"resources.limits"`
	Env       []EnvVar            `sm:"env"`
}

type ContainerPort struct {
	Port int `sm:"containerPort"`
}

type ContainerResources struct {
	Cpu    string `json:"cpu"`
	Memory string `json:"memory"`
}

type EnvVar struct {
	Name  string `sm:"name"`
	Value string `sm:"value"`
}
