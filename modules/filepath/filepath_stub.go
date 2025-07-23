//go:build !module_filepath

package filepath

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
