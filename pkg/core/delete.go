package core

type Delete[T any, R any] struct {
	Action[T, R]
	Id string
}

func (d *Delete[T, R]) Run() error {
	return d.Backend.Delete(d.Id)
}
