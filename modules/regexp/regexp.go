//go:build !regexp_module_stub

package regexp

import (
	modregexp "github.com/risor-io/risor/modules/regexp"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modregexp.Module()
}
