//go:build !module_rand

package rand

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
