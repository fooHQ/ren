//go:build module_net

package net

import (
	modnet "github.com/risor-io/risor/modules/net"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return modnet.Module()
}
