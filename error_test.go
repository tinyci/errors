package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestErrors(t *testing.T) {
	errorFuncs := map[string]func() error{
		"plain": func() error { return New("error") },
		"errorf": func() error {
			return Errorf("error: %w", errors.New("some other error"))
		},
		"grpc": func() error {
			return GRPC(codes.Aborted, "whoops: %v", 420)
		},
	}

	for name, fun := range errorFuncs {
		WithCaller = false
		WithCallerVerbose = false
		WithStack = false

		err := fun()
		assert.Error(t, err, name)
		assert.Contains(t, err.Error(), "error", name)
		assert.IsType(t, err, WithLocation{}, name)
		assert.Error(t, err.(WithLocation).Err)

		WithCaller = true
		err = fun()
		assert.Error(t, err, name)
		assert.Contains(t, err.Error(), "tinyci/errors.TestErrors", name)
		assert.IsType(t, err, WithLocation{}, name)

		WithCallerVerbose = true
		err = fun()
		assert.Error(t, err, name)
		assert.Contains(t, err.Error(), "error_test.go", name)
		assert.IsType(t, err, WithLocation{}, name)

		WithStack = true

		err = fun()
		assert.Error(t, err, name)
		assert.Contains(t, err.Error(), "STACK TRACE", name)
		assert.Contains(t, err.Error(), "WithLocation.Error", name)
		assert.IsType(t, err, WithLocation{}, name)
	}
}
