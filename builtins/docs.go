package builtins

import (
	modbuiltins "github.com/deepnoodle-ai/risor/v2/pkg/builtins"
	"github.com/deepnoodle-ai/risor/v2/pkg/object"
)

// ownDocs documents the builtins that Ren defines itself, as opposed to the
// ones it re-exports from Risor (whose documentation is taken from
// modbuiltins.Docs).
var ownDocs = []object.FuncSpec{
	{Name: "import", Doc: "Load a module and return it; the argument is a package path or a builtin:// URL", Args: []string{"url"}, Returns: "module", Example: `import("builtin://os")`},
	{Name: "print", Doc: "Write the arguments to standard output separated by spaces and followed by a newline", Args: []string{"value..."}, Returns: "nil"},
	{Name: "printf", Doc: "Write a formatted string to standard output", Args: []string{"format", "value..."}, Returns: "nil"},
	{Name: "pack", Doc: "Serialize a map into a little-endian byte buffer according to a schema", Args: []string{"schema", "data"}, Returns: "bytes"},
	{Name: "packsize", Doc: "Return the total byte size of a schema without packing any data", Args: []string{"schema"}, Returns: "int"},
	{Name: "unpack", Doc: "Deserialize a little-endian byte buffer into a map according to a schema", Args: []string{"schema", "buffer"}, Returns: "map"},
}

// Docs returns documentation for every global builtin, combining Ren's own
// builtins with the documentation Risor ships for the ones Ren re-exports.
func Docs() []object.FuncSpec {
	own := make(map[string]struct{}, len(ownDocs))
	for _, s := range ownDocs {
		own[s.Name] = struct{}{}
	}

	specs := make([]object.FuncSpec, len(ownDocs))
	copy(specs, ownDocs)

	for _, s := range modbuiltins.Docs() {
		if _, registered := builtins[s.Name]; !registered {
			continue
		}
		if _, isOwn := own[s.Name]; isOwn {
			continue
		}
		specs = append(specs, s)
	}
	return specs
}
