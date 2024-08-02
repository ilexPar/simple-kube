package cluster

type ClusterDelete[T ClusterResources] struct {
	Action[T]
	Id string
}

func (d *ClusterDelete[T]) Run() error {
	return d.api.Delete(d.Id)
}
