package errors

import (
	"errors"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

const (
	ERROR_NOT_FOUND = "not found"
)

func Format(err error) error {
	if kerrors.IsNotFound(err) {
		return errors.New(ERROR_NOT_FOUND)
	}

	return err
}
