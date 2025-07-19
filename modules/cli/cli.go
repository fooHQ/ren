//go:build module_cli

package cli

import (
	modcli "github.com/risor-io/risor/modules/cli"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modcli.Module()
}

func Builtins() map[string]object.Object {
	return nil
}
