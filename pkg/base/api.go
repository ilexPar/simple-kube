package base

import (
	"context"

	"k8s.io/client-go/kubernetes"
)

type KubeAPI struct {
	Client  kubernetes.Interface
	Context context.Context
	Opts    QueryOpts
}

func (api *KubeAPI) Config(ctx context.Context, client kubernetes.Interface) {
	api.Client = client
	api.Context = ctx
}

func (api *KubeAPI) SetOpts(opts QueryOpts) {
	api.Opts = opts
}
