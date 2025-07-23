//go:build !module_os

package os

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
