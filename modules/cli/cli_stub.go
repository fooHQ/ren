//go:build !module_cli

package cli

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
