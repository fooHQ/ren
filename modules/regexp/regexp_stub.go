//go:build !module_regexp

package regexp

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
