//go:build !module_urlpath

package urlpath

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
