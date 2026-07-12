// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

// Package os implements the Ren "os" module, exposing process, environment,
// and user/group information to scripts. Every operation is dispatched through
// the OS abstraction stored on the context.
package os

import (
	"context"
	"fmt"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"

	"github.com/foohq/ren"
	"github.com/foohq/ren/objects"
)

// Args returns the script's command-line arguments as a list of strings. It
// takes no arguments.
func Args(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 0 {
		return nil, object.NewArgsError("os.args", 0, len(args))
	}
	argz := ren.GetOS(ctx).Args()
	items := make([]object.Object, len(argz))
	for i, arg := range argz {
		items[i] = object.NewString(arg)
	}
	return object.NewList(items), nil
}

// Exit invokes the configured exit handler with the given status code. A
// non-zero code is also surfaced as an error. It takes a single int argument.
func Exit(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("os.exit", 1, len(args))
	}
	code, ok := args[0].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("os.exit: expected int, got %s", args[0].Type())
	}
	ren.GetOS(ctx).Exit(int(code.Value()))
	if code.Value() != 0 {
		return nil, fmt.Errorf("os.exit: exited with code %d", code.Value())
	}
	return object.Nil, nil
}

// Chdir changes the script's working directory. It takes a single directory
// argument.
func Chdir(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("os.chdir", 1, len(args))
	}
	dir, ok := args[0].(*object.String)
	if !ok {
		return nil, fmt.Errorf("os.chdir: expected string, got %s", args[0].Type())
	}
	if err := ren.GetOS(ctx).Chdir(dir.Value()); err != nil {
		return nil, object.NewError(err)
	}
	return object.Nil, nil
}

// Getwd returns the script's current working directory. It takes no arguments.
func Getwd(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 0 {
		return nil, object.NewArgsError("os.getwd", 0, len(args))
	}
	dir, err := ren.GetOS(ctx).Getwd()
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewString(dir), nil
}

// TempDir returns the default directory for temporary files. It takes no
// arguments.
func TempDir(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 0 {
		return nil, object.NewArgsError("os.temp_dir", 0, len(args))
	}
	return object.NewString(ren.GetOS(ctx).TempDir()), nil
}

// Getenv returns the value of an environment variable, or an empty string if it
// is unset. It takes a single key argument.
func Getenv(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("os.getenv", 1, len(args))
	}
	key, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	return object.NewString(ren.GetOS(ctx).Getenv(key)), nil
}

// Setenv sets an environment variable. It takes two string arguments: the key
// and the value.
func Setenv(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, object.NewArgsError("os.setenv", 2, len(args))
	}
	key, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	value, err := object.AsString(args[1])
	if err != nil {
		return nil, err
	}
	if err := ren.GetOS(ctx).Setenv(key, value); err != nil {
		return nil, object.NewError(err)
	}
	return object.Nil, nil
}

// Unsetenv removes an environment variable. It takes a single key argument.
func Unsetenv(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("os.unsetenv", 1, len(args))
	}
	key, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	if err := ren.GetOS(ctx).Unsetenv(key); err != nil {
		return nil, object.NewError(err)
	}
	return object.Nil, nil
}

// UserCacheDir returns the default root directory for user-specific cached
// data. It takes no arguments.
func UserCacheDir(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 0 {
		return nil, object.NewArgsError("os.user_cache_dir", 0, len(args))
	}
	dir, err := ren.GetOS(ctx).UserCacheDir()
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewString(dir), nil
}

// UserConfigDir returns the default root directory for user-specific
// configuration data. It takes no arguments.
func UserConfigDir(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 0 {
		return nil, object.NewArgsError("os.user_config_dir", 0, len(args))
	}
	dir, err := ren.GetOS(ctx).UserConfigDir()
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewString(dir), nil
}

// UserHomeDir returns the current user's home directory. It takes no arguments.
func UserHomeDir(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 0 {
		return nil, object.NewArgsError("os.user_home_dir", 0, len(args))
	}
	dir, err := ren.GetOS(ctx).UserHomeDir()
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewString(dir), nil
}

// Environ returns the environment as a list of "key=value" strings. It takes no
// arguments.
func Environ(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 0 {
		return nil, object.NewArgsError("os.environ", 0, len(args))
	}
	envVars := ren.GetOS(ctx).Environ()
	items := make([]object.Object, len(envVars))
	for i, envVar := range envVars {
		items[i] = object.NewString(envVar)
	}
	return object.NewList(items), nil
}

// Getpid returns the process ID of the caller. It takes no arguments.
func Getpid(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 0 {
		return nil, object.NewArgsError("os.getpid", 0, len(args))
	}
	return object.NewInt(int64(ren.GetOS(ctx).Getpid())), nil
}

// Getuid returns the numeric user ID of the caller. It takes no arguments.
func Getuid(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 0 {
		return nil, object.NewArgsError("os.getuid", 0, len(args))
	}
	return object.NewInt(int64(ren.GetOS(ctx).Getuid())), nil
}

// Hostname returns the host name reported by the kernel. It takes no arguments.
func Hostname(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 0 {
		return nil, object.NewArgsError("os.hostname", 0, len(args))
	}
	hostname, err := ren.GetOS(ctx).Hostname()
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewString(hostname), nil
}

// CurrentUser returns the current user as a map of its fields. It takes no
// arguments.
func CurrentUser(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 0 {
		return nil, object.NewArgsError("os.current_user", 0, len(args))
	}
	user, err := ren.GetOS(ctx).CurrentUser()
	if err != nil {
		return nil, object.NewError(err)
	}
	return wrapUser(user), nil
}

// LookupUser looks up a user by username and returns it as a map. It takes a
// single username argument.
func LookupUser(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("os.lookup_user", 1, len(args))
	}
	name, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	user, lookupErr := ren.GetOS(ctx).LookupUser(name)
	if lookupErr != nil {
		return nil, object.NewError(lookupErr)
	}
	return wrapUser(user), nil
}

// LookupUid looks up a user by numeric ID and returns it as a map. It takes a
// single uid argument.
func LookupUid(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("os.lookup_uid", 1, len(args))
	}
	uid, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	user, lookupErr := ren.GetOS(ctx).LookupUid(uid)
	if lookupErr != nil {
		return nil, object.NewError(lookupErr)
	}
	return wrapUser(user), nil
}

// LookupGroup looks up a group by name and returns it as a map. It takes a
// single group-name argument.
func LookupGroup(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("os.lookup_group", 1, len(args))
	}
	name, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	group, lookupErr := ren.GetOS(ctx).LookupGroup(name)
	if lookupErr != nil {
		return nil, object.NewError(lookupErr)
	}
	return wrapGroup(group), nil
}

// LookupGid looks up a group by numeric ID and returns it as a map. It takes a
// single gid argument.
func LookupGid(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("os.lookup_gid", 1, len(args))
	}
	gid, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	group, lookupErr := ren.GetOS(ctx).LookupGid(gid)
	if lookupErr != nil {
		return nil, object.NewError(lookupErr)
	}
	return wrapGroup(group), nil
}

// Stdin resolves the "stdin" module attribute to a file object wrapping the
// script's standard input.
func Stdin(ctx context.Context, name string) (object.Object, error) {
	f := ren.GetOS(ctx).Stdin()
	return objects.NewFile(ctx, f, "/dev/stdin"), nil
}

// Stdout resolves the "stdout" module attribute to a file object wrapping the
// script's standard output.
func Stdout(ctx context.Context, name string) (object.Object, error) {
	f := ren.GetOS(ctx).Stdout()
	return objects.NewFile(ctx, f, "/dev/stdout"), nil
}

// wrapUser wraps an os.User as a Risor object
func wrapUser(user ren.User) object.Object {
	items := map[string]object.Object{
		"uid":      object.NewString(user.Uid()),
		"gid":      object.NewString(user.Gid()),
		"username": object.NewString(user.Username()),
		"name":     object.NewString(user.Name()),
		"home_dir": object.NewString(user.HomeDir()),
	}
	return object.NewMap(items)
}

// wrapGroup wraps an os.Group as a Risor object
func wrapGroup(group ren.Group) object.Object {
	items := map[string]object.Object{
		"gid":  object.NewString(group.Gid()),
		"name": object.NewString(group.Name()),
	}
	return object.NewMap(items)
}

// Module returns the "os" module with all of its functions and dynamic
// attributes registered.
func Module() *object.Module {
	return object.NewBuiltinsModule("os", map[string]object.Object{
		"args":            object.NewBuiltin("args", Args),
		"chdir":           object.NewBuiltin("chdir", Chdir),
		"current_user":    object.NewBuiltin("current_user", CurrentUser),
		"environ":         object.NewBuiltin("environ", Environ),
		"exit":            object.NewBuiltin("exit", Exit),
		"getenv":          object.NewBuiltin("getenv", Getenv),
		"getpid":          object.NewBuiltin("getpid", Getpid),
		"getuid":          object.NewBuiltin("getuid", Getuid),
		"getwd":           object.NewBuiltin("getwd", Getwd),
		"hostname":        object.NewBuiltin("hostname", Hostname),
		"lookup_gid":      object.NewBuiltin("lookup_gid", LookupGid),
		"lookup_group":    object.NewBuiltin("lookup_group", LookupGroup),
		"lookup_uid":      object.NewBuiltin("lookup_uid", LookupUid),
		"lookup_user":     object.NewBuiltin("lookup_user", LookupUser),
		"setenv":          object.NewBuiltin("setenv", Setenv),
		"temp_dir":        object.NewBuiltin("temp_dir", TempDir),
		"unsetenv":        object.NewBuiltin("unsetenv", Unsetenv),
		"user_cache_dir":  object.NewBuiltin("user_cache_dir", UserCacheDir),
		"user_config_dir": object.NewBuiltin("user_config_dir", UserConfigDir),
		"user_home_dir":   object.NewBuiltin("user_home_dir", UserHomeDir),
		"stdin":           object.NewDynamicAttr("stdin", Stdin),
		"stdout":          object.NewDynamicAttr("stdout", Stdout),
	})
}
