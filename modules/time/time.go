//go:build !time_module_stub

package time

import (
	modtime "github.com/risor-io/risor/modules/time"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modtime.Module()
}
