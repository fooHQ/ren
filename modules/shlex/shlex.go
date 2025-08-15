//go:build !shlex_module_stub

package shlex

import (
	modshlex "github.com/risor-io/risor/modules/shlex"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modshlex.Module()
}
