//go:build !module_strings

package strings

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
