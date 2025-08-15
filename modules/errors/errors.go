//go:build !errors_module_stub

package errors

import (
	moderrors "github.com/risor-io/risor/modules/errors"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return moderrors.Module()
}
