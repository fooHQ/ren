//go:build !filepath_module_stub

package filepath

import (
	modfilepath "github.com/risor-io/risor/modules/filepath"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modfilepath.Module()
}
