//go:build !http_module_stub

package http

import (
	modhttp "github.com/risor-io/risor/modules/http"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modhttp.Module()
}
