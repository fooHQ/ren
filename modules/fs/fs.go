// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

package fs

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"

	"github.com/foohq/ren"
	"github.com/foohq/ren/objects"
)

func OpenFile(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 3 {
		return nil, object.NewArgsError("os.open_file", 3, len(args))
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	mode, err := object.AsString(args[1])
	if err != nil {
		return nil, err
	}
	perm, err := object.AsInt(args[2])
	if err != nil {
		return nil, err
	}
	flags, err := modeToFlags(mode)
	if err != nil {
		return nil, object.NewValueError(err)
	}
	f, err := ren.GetOS(ctx).OpenFile(path, flags, ren.FileMode(perm))
	if err != nil {
		return nil, object.NewError(err)
	}
	return objects.NewFile(ctx, f, path), nil
}

func ReadFile(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("os.read_file", 1, len(args))
	}
	filename, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	bytes, ioErr := ren.GetOS(ctx).ReadFile(filename)
	if ioErr != nil {
		return nil, object.NewError(ioErr)
	}
	return object.NewBytes(bytes), nil
}

func ReadDir(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("os.read_dir", 1, len(args))
	}
	dirName, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	entries, ioErr := ren.GetOS(ctx).ReadDir(dirName)
	if ioErr != nil {
		return nil, object.NewError(ioErr)
	}
	items := make([]object.Object, 0, len(entries))
	for _, entry := range entries {
		items = append(items, objects.NewDirEntry(entry))
	}
	return object.NewList(items), nil
}

func WriteFile(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 3 {
		return nil, object.NewArgsError("os.write_file", 3, len(args))
	}
	filename, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	var data []byte
	switch arg := args[1].(type) {
	case *object.Bytes:
		data = arg.Value()
	case *object.String:
		data = []byte(arg.Value())
	default:
		return nil, fmt.Errorf("os.write_file: expected byte_slice or string, got %s", args[1].Type())
	}
	perm, err := object.AsInt(args[2])
	if err != nil {
		return nil, err
	}
	if err := ren.GetOS(ctx).WriteFile(filename, data, ren.FileMode(perm)); err != nil {
		return nil, object.NewError(err)
	}
	return object.Nil, nil
}

func Remove(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("os.remove", 1, len(args))
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	if err := ren.GetOS(ctx).Remove(path); err != nil {
		return nil, object.NewError(err)
	}
	return object.Nil, nil
}

func RemoveAll(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("os.remove_all", 1, len(args))
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	if err := ren.GetOS(ctx).RemoveAll(path); err != nil {
		return nil, object.NewError(err)
	}
	return object.Nil, nil
}

func Rename(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, object.NewArgsError("os.rename", 2, len(args))
	}
	oldpath, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	newpath, err := object.AsString(args[1])
	if err != nil {
		return nil, err
	}
	if err := ren.GetOS(ctx).Rename(oldpath, newpath); err != nil {
		return nil, object.NewError(err)
	}
	return object.Nil, nil
}

func Stat(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("os.stat", 1, len(args))
	}
	name, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	info, ioErr := ren.GetOS(ctx).Stat(name)
	if ioErr != nil {
		return nil, object.NewError(ioErr)
	}
	return objects.NewFileInfo(info), nil
}

func MkdirAll(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, object.NewArgsError("os.mkdir_all", 2, len(args))
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	perm, err := object.AsInt(args[1])
	if err != nil {
		return nil, err
	}
	if err := ren.GetOS(ctx).MkdirAll(path, ren.FileMode(perm)); err != nil {
		return nil, object.NewError(err)
	}
	return object.Nil, nil
}

func Mkdir(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, object.NewArgsError("os.mkdir", 2, len(args))
	}
	dir, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	perm, err := object.AsInt(args[1])
	if err != nil {
		return nil, err
	}
	if err := ren.GetOS(ctx).Mkdir(dir, ren.FileMode(perm)); err != nil {
		return nil, object.NewError(err)
	}
	return object.Nil, nil
}

func MkdirTemp(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, object.NewArgsError("os.mkdir_temp", 2, len(args))
	}
	dir, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	pattern, err := object.AsString(args[1])
	if err != nil {
		return nil, err
	}
	tempDir, ioErr := ren.GetOS(ctx).MkdirTemp(dir, pattern)
	if ioErr != nil {
		return nil, object.NewError(ioErr)
	}
	return object.NewString(tempDir), nil
}

func Symlink(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, object.NewArgsError("os.symlink", 2, len(args))
	}
	oldname, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	newname, err := object.AsString(args[1])
	if err != nil {
		return nil, err
	}
	if err := ren.GetOS(ctx).Symlink(oldname, newname); err != nil {
		return nil, object.NewError(err)
	}
	return object.Nil, nil
}

func modeToFlags(mode string) (int, error) {
	switch mode {
	case "r":
		return os.O_RDONLY, nil
	case "r+":
		return os.O_RDWR, nil
	case "w":
		return os.O_WRONLY | os.O_CREATE | os.O_TRUNC, nil
	case "w+":
		return os.O_RDWR | os.O_CREATE | os.O_TRUNC, nil
	case "a":
		return os.O_WRONLY | os.O_CREATE | os.O_APPEND, nil
	case "a+":
		return os.O_RDWR | os.O_CREATE | os.O_APPEND, nil
	case "ax", "xa":
		return os.O_WRONLY | os.O_CREATE | os.O_APPEND | os.O_EXCL, nil
	case "ax+", "xa+":
		return os.O_RDWR | os.O_CREATE | os.O_APPEND | os.O_EXCL, nil
	case "wx", "xw":
		return os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_EXCL, nil
	case "wx+", "xw+":
		return os.O_RDWR | os.O_CREATE | os.O_TRUNC | os.O_EXCL, nil
	default:
		return 0, errors.New("unsupported file mode")
	}
}

func Module() *object.Module {
	return object.NewBuiltinsModule("fs", map[string]object.Object{
		"mkdir":      object.NewBuiltin("mkdir", Mkdir),
		"mkdir_all":  object.NewBuiltin("mkdir_all", MkdirAll),
		"mkdir_temp": object.NewBuiltin("mkdir_temp", MkdirTemp),
		"open_file":  object.NewBuiltin("open_file", OpenFile),
		"read_file":  object.NewBuiltin("read_file", ReadFile),
		"write_file": object.NewBuiltin("write_file", WriteFile),
		"remove":     object.NewBuiltin("remove", Remove),
		"remove_all": object.NewBuiltin("remove_all", RemoveAll),
		"rename":     object.NewBuiltin("rename", Rename),
		"stat":       object.NewBuiltin("stat", Stat),
		"symlink":    object.NewBuiltin("symlink", Symlink),
		"read_dir":   object.NewBuiltin("read_dir", ReadDir),
		// TODO: uncomment these once Risor has a better error handling system
		/*"err_not_exist":  object.NewError(fs.ErrNotExist),
		"err_exist":      object.NewError(fs.ErrExist),
		"err_permission": object.NewError(fs.ErrPermission),
		"err_closed":     object.NewError(fs.ErrClosed),
		"err_invalid":    object.NewError(fs.ErrInvalid),*/
	})
}
