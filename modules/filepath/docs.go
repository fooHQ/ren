package filepath

import "github.com/deepnoodle-ai/risor/v2/pkg/object"

// ModuleDoc returns the module-level documentation for "filepath".
func ModuleDoc() string {
	return "URL-aware path manipulation helpers."
}

// Docs returns documentation for every name exposed by the "filepath" module.
func Docs() []object.FuncSpec {
	return docs
}

var docs = []object.FuncSpec{
	{Name: "abs", Doc: "Return the absolute form of a path resolved against a working directory", Args: []string{"path", "wd"}, Returns: "string"},
	{Name: "base", Doc: "Return the last element of a path", Args: []string{"path"}, Returns: "string"},
	{Name: "clean", Doc: "Return the shortest equivalent form of a path", Args: []string{"path"}, Returns: "string"},
	{Name: "dir", Doc: "Return all but the last element of a path", Args: []string{"path"}, Returns: "string"},
	{Name: "ext", Doc: "Return the file name extension of a path", Args: []string{"path"}, Returns: "string"},
	{Name: "is_abs", Doc: "Report whether a path is absolute", Args: []string{"path"}, Returns: "bool"},
	{Name: "join", Doc: "Join any number of path elements into a single cleaned path", Args: []string{"elem..."}, Returns: "string"},
	{Name: "match", Doc: "Report whether a name matches a shell pattern", Args: []string{"pattern", "name"}, Returns: "bool"},
	{Name: "split", Doc: "Split a path into directory and file components", Args: []string{"path"}, Returns: "list"},
}
