package shlex

import "github.com/deepnoodle-ai/risor/v2/pkg/object"

// ModuleDoc returns the module-level documentation for "shlex".
func ModuleDoc() string {
	return "Split command-line strings into arguments using shell-style quoting."
}

// Docs returns documentation for every name exposed by the "shlex" module.
func Docs() []object.FuncSpec {
	return docs
}

var docs = []object.FuncSpec{
	{Name: "argv", Doc: "Split a command-line string into a list of arguments following shell-style quoting and escaping", Args: []string{"line"}, Returns: "list"},
}
