// Package shlex implements the Ren "shlex" module, which splits command-line
// strings into arguments using shell-style quoting rules.
package shlex

import (
	"context"

	"github.com/u-root/u-root/pkg/shlex"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
)

// Argv splits a command-line string into a list of arguments following
// shell-style quoting and escaping. It takes a single string argument.
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

// Module returns the "shlex" module with its functions registered.
func Module() *object.Module {
	return object.NewBuiltinsModule("shlex", map[string]object.Object{
		"argv": object.NewBuiltin("argv", Argv),
	})
}
