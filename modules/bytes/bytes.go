//go:build !bytes_module_stub

package bytes

import (
	modbytes "github.com/risor-io/risor/modules/bytes"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modbytes.Module()
}
