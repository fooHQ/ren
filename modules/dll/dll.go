//go:build !windows

// Package dll implements the Ren "dll" module for loading dynamic-link
// libraries and calling their exported procedures. It is only functional on
// Windows; on other platforms its operations return an unsupported-platform
// error.
package dll

import (
	"context"
	"fmt"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
)

// Load reports that dynamic library loading is not supported on this platform.
func Load(ctx context.Context, args ...object.Object) (object.Object, error) {
	return nil, fmt.Errorf("dll.load: not supported on this platform")
}

// Module returns the "dll" module. On non-Windows platforms its load function
// always fails.
func Module() *object.Module {
	return object.NewBuiltinsModule("dll", map[string]object.Object{
		"load": object.NewBuiltin("load", Load),
	})
}
