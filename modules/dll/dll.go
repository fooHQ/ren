//go:build !windows

package dll

import (
	"context"
	"fmt"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
)

func Load(ctx context.Context, args ...object.Object) (object.Object, error) {
	return nil, fmt.Errorf("dll.load: not supported on this platform")
}

func Module() *object.Module {
	return object.NewBuiltinsModule("dll", map[string]object.Object{
		"load": object.NewBuiltin("load", Load),
	})
}
