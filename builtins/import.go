package builtins

import (
	"context"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"

	"github.com/foohq/ren"
)

// Import loads a module and returns it. The single argument is the module
// reference: a package-relative path (e.g. "utils/log") or a built-in module
// URL (e.g. "builtin://os").
func Import(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("import", 1, len(args))
	}

	name, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}

	mod, err := ren.GetImporter(ctx).Import(ctx, name)
	if err != nil {
		return nil, err
	}

	return mod, nil
}
