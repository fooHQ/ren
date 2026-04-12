package shlex

import (
	"context"

	"github.com/u-root/u-root/pkg/shlex"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
)

func Argv(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("shlex.argv", 1, len(args))
	}
	data, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	return object.NewStringList(shlex.Argv(data)), nil
}

func Module() *object.Module {
	return object.NewBuiltinsModule("shlex", map[string]object.Object{
		"argv": object.NewBuiltin("argv", Argv),
	})
}
