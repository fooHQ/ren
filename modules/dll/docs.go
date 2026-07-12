package dll

import "github.com/deepnoodle-ai/risor/v2/pkg/object"

// ModuleDoc returns the module-level documentation for "dll".
func ModuleDoc() string {
	return "Load dynamic-link libraries and call their exported procedures (Windows only)."
}

// Docs returns documentation for every name exposed by the "dll" module.
func Docs() []object.FuncSpec {
	return docs
}

var docs = []object.FuncSpec{
	{Name: "load", Doc: "Open the dynamic-link library at the given path and return a handle (Windows only; fails on other platforms)", Args: []string{"path"}, Returns: "handle"},
}
