//go:build !math_module_stub

package math

import (
	modmath "github.com/risor-io/risor/modules/math"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modmath.Module()
}
