//go:build module_os

package os

import (
	modos "github.com/risor-io/risor/modules/os"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modos.Module()
}
