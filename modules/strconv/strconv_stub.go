//go:build !module_strconv

package strconv

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
