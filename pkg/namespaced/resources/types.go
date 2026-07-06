package resources

import (
	"context"

	"k8s.io/client-go/kubernetes"

	"github.com/ilexPar/simple-kube/pkg/base"
)

type NamespacedResourceAPI[R base.KubernetesResources] interface {
	Config(ctx context.Context, k8s kubernetes.Interface)
	SetOpts(opts base.QueryOpts)
	SetContext(ctx context.Context)

	Get(name, namespace string) (*R, error)
	Create(namespace string, obj *R) error
	Update(namespace string, obj *R) error
	List(namespace string) ([]R, error)
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
	Cpu    string `sm:"cpu"`
	Memory string `sm:"memory"`
}

type EnvVar struct {
	Name  string `sm:"name"`
	Value string `sm:"value"`
}
