//go:build !module_errors

package errors

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}

func Builtins() map[string]object.Object {
	return nil
}
