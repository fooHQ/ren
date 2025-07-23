//go:build module_exec

package exec

import (
	modexec "github.com/risor-io/risor/modules/exec"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modexec.Module()
}
