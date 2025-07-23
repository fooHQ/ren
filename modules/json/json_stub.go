//go:build !module_json

package json

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
