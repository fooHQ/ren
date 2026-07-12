package os

import "github.com/deepnoodle-ai/risor/v2/pkg/object"

// ModuleDoc returns the module-level documentation for "os".
func ModuleDoc() string {
	return "Process, environment, and user/group information."
}

// Docs returns documentation for every name exposed by the "os" module,
// including its dynamic stream attributes.
func Docs() []object.FuncSpec {
	return docs
}

var docs = []object.FuncSpec{
	{Name: "args", Doc: "Return the script's command-line arguments", Returns: "list"},
	{Name: "exit", Doc: "Exit the script with the given status code", Args: []string{"code"}, Returns: "nil"},
	{Name: "chdir", Doc: "Change the script's working directory", Args: []string{"dir"}, Returns: "nil"},
	{Name: "getwd", Doc: "Return the script's current working directory", Returns: "string"},
	{Name: "temp_dir", Doc: "Return the default directory for temporary files", Returns: "string"},
	{Name: "getenv", Doc: "Return the value of an environment variable, or an empty string if unset", Args: []string{"key"}, Returns: "string"},
	{Name: "setenv", Doc: "Set an environment variable", Args: []string{"key", "value"}, Returns: "nil"},
	{Name: "unsetenv", Doc: "Remove an environment variable", Args: []string{"key"}, Returns: "nil"},
	{Name: "environ", Doc: "Return the environment as a list of \"key=value\" strings", Returns: "list"},
	{Name: "getpid", Doc: "Return the process ID of the caller", Returns: "int"},
	{Name: "getuid", Doc: "Return the numeric user ID of the caller", Returns: "int"},
	{Name: "hostname", Doc: "Return the host name reported by the kernel", Returns: "string"},
	{Name: "current_user", Doc: "Return the current user as a map of its fields", Returns: "map"},
	{Name: "lookup_user", Doc: "Look up a user by username", Args: []string{"username"}, Returns: "map"},
	{Name: "lookup_uid", Doc: "Look up a user by numeric ID", Args: []string{"uid"}, Returns: "map"},
	{Name: "lookup_group", Doc: "Look up a group by name", Args: []string{"name"}, Returns: "map"},
	{Name: "lookup_gid", Doc: "Look up a group by numeric ID", Args: []string{"gid"}, Returns: "map"},
	{Name: "user_home_dir", Doc: "Return the current user's home directory", Returns: "string"},
	{Name: "user_cache_dir", Doc: "Return the default root directory for user-specific cached data", Returns: "string"},
	{Name: "user_config_dir", Doc: "Return the default root directory for user-specific configuration", Returns: "string"},
	{Name: "stdin", Doc: "The standard input stream as a file object", Returns: "file"},
	{Name: "stdout", Doc: "The standard output stream as a file object", Returns: "file"},
}
