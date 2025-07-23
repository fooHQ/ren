//go:build !module_exec

package exec

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
