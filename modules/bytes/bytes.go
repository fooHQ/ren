//go:build module_bytes

package bytes

import (
	modbytes "github.com/risor-io/risor/modules/bytes"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modbytes.Module()
}
