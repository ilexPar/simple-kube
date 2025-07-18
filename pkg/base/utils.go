package base

import "fmt"

func FlattenLabels(labels map[string]string) string {
	var flatLabels string
	labelLength := len(labels)
	idx := 0
	for i, v := range labels {
		label := fmt.Sprintf("%s=%s", i, v)
		flatLabels += label
		if idx < (labelLength - 1) {
			flatLabels += ","
		}
		idx++
	}
	return flatLabels
}

func Cast[T any](obj interface{}) (T, error) {
	if obj == nil {
		var zero T
		return zero, fmt.Errorf("cannot cast nil to %T", zero)
	}
	val, ok := obj.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("cannot cast %T to %T", obj, zero)
	}
	return val, nil
}
