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
