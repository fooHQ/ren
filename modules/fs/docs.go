package fs

import "github.com/deepnoodle-ai/risor/v2/pkg/object"

// ModuleDoc returns the module-level documentation for "fs".
func ModuleDoc() string {
	return "Filesystem operations over the scheme-based virtual filesystems registered with the runtime."
}

// Docs returns documentation for every name exposed by the "fs" module,
// including its error sentinels.
func Docs() []object.FuncSpec {
	return docs
}

var docs = []object.FuncSpec{
	{Name: "open_file", Doc: "Open a file and return a file object; mode is a fopen-style string such as \"r\", \"w\", or \"a+\"", Args: []string{"path", "mode", "perm"}, Returns: "file"},
	{Name: "read_file", Doc: "Read a file and return its contents", Args: []string{"path"}, Returns: "bytes"},
	{Name: "write_file", Doc: "Write data to a file, creating it as needed", Args: []string{"path", "data", "perm"}, Returns: "nil"},
	{Name: "read_dir", Doc: "List a directory and return its entries", Args: []string{"path"}, Returns: "list"},
	{Name: "stat", Doc: "Return a file_info object describing a file", Args: []string{"path"}, Returns: "file_info"},
	{Name: "mkdir", Doc: "Create a single directory", Args: []string{"path", "perm"}, Returns: "nil"},
	{Name: "mkdir_all", Doc: "Create a directory along with any missing parents", Args: []string{"path", "perm"}, Returns: "nil"},
	{Name: "mkdir_temp", Doc: "Create a new temporary directory and return its path", Args: []string{"dir", "pattern"}, Returns: "string"},
	{Name: "remove", Doc: "Delete a file or empty directory", Args: []string{"path"}, Returns: "nil"},
	{Name: "remove_all", Doc: "Delete a path and any children it contains", Args: []string{"path"}, Returns: "nil"},
	{Name: "rename", Doc: "Move a file or directory (cannot cross filesystem boundaries)", Args: []string{"oldpath", "newpath"}, Returns: "nil"},
	{Name: "symlink", Doc: "Create a symbolic link (cannot cross filesystem boundaries)", Args: []string{"oldname", "newname"}, Returns: "nil"},
	{Name: "err_not_exist", Doc: "Error sentinel: the file does not exist", Returns: "error"},
	{Name: "err_exist", Doc: "Error sentinel: the file already exists", Returns: "error"},
	{Name: "err_permission", Doc: "Error sentinel: permission denied", Returns: "error"},
	{Name: "err_closed", Doc: "Error sentinel: the file is already closed", Returns: "error"},
	{Name: "err_invalid", Doc: "Error sentinel: invalid argument", Returns: "error"},
}
