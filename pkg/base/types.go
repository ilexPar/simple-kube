package base

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type QueryOpts struct {
	List metav1.ListOptions
}

type ResourceInterface interface {
	Load(from, into interface{}) error
	Dump(from interface{}) (interface{}, error)
}
