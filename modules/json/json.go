//go:build !json_module_stub

package json

import (
	modjson "github.com/risor-io/risor/modules/json"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modjson.Module()
}
