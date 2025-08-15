//go:build !strconv_module_stub

package strconv

import (
	modstrconv "github.com/risor-io/risor/modules/strconv"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modstrconv.Module()
}
