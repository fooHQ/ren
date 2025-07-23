//go:build !module_math

package math

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
