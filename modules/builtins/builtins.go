//go:build module_builtins

package builtins

import (
	"github.com/risor-io/risor/builtins"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}

func Builtins() map[string]object.Object {
	return builtins.Builtins()
}
