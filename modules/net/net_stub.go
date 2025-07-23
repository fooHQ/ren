//go:build !module_net

package net

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
