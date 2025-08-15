//go:build !cli_module_stub

package cli

import (
	modcli "github.com/risor-io/risor/modules/cli"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modcli.Module()
}
