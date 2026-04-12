package filepath

import (
	"context"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"

	"github.com/foohq/urlpath"
)

func Abs(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, object.NewArgsError("filepath.abs", 2, len(args))
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	wd, err := object.AsString(args[1])
	if err != nil {
		return nil, err
	}
	abs, err := urlpath.Abs(path, wd)
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewString(abs), nil
}

func Base(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("filepath.base", 1, len(args))
	}
	path, rerr := object.AsString(args[0])
	if rerr != nil {
		return nil, rerr
	}
	base, err := urlpath.Base(path)
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewString(base), nil
}

func Clean(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("filepath.clean", 1, len(args))
	}
	path, rerr := object.AsString(args[0])
	if rerr != nil {
		return nil, rerr
	}
	cleanPath, err := urlpath.Clean(path)
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewString(cleanPath), nil
}

func Dir(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("filepath.dir", 1, len(args))
	}
	path, rerr := object.AsString(args[0])
	if rerr != nil {
		return nil, rerr
	}
	dirPath, err := urlpath.Dir(path)
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewString(dirPath), nil
}

func Ext(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("filepath.ext", 1, len(args))
	}
	path, rerr := object.AsString(args[0])
	if rerr != nil {
		return nil, rerr
	}
	extension, err := urlpath.Ext(path)
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewString(extension), nil
}

func IsAbs(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("filepath.is_abs", 1, len(args))
	}
	path, rerr := object.AsString(args[0])
	if rerr != nil {
		return nil, rerr
	}
	isAbs, err := urlpath.IsAbs(path)
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewBool(isAbs), nil
}

func Join(ctx context.Context, args ...object.Object) (object.Object, error) {
	paths := make([]string, len(args))
	for i, arg := range args {
		path, rerr := object.AsString(arg)
		if rerr != nil {
			return nil, rerr
		}
		paths[i] = path
	}
	res, err := urlpath.Join(paths...)
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewString(res), nil
}

func Match(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, object.NewArgsError("filepath.match", 2, len(args))
	}
	pattern, rerr := object.AsString(args[0])
	if rerr != nil {
		return nil, rerr
	}
	name, rerr := object.AsString(args[1])
	if rerr != nil {
		return nil, rerr
	}
	matched, err := urlpath.Match(pattern, name)
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewBool(matched), nil
}

func Split(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("filepath.split", 1, len(args))
	}
	path, rerr := object.AsString(args[0])
	if rerr != nil {
		return nil, rerr
	}
	dir, file, err := urlpath.Split(path)
	if err != nil {
		return nil, object.NewError(err)
	}
	return object.NewList([]object.Object{
		object.NewString(dir),
		object.NewString(file),
	}), nil
}

func Module() *object.Module {
	return object.NewBuiltinsModule("filepath", map[string]object.Object{
		"abs":    object.NewBuiltin("abs", Abs),
		"base":   object.NewBuiltin("base", Base),
		"clean":  object.NewBuiltin("clean", Clean),
		"dir":    object.NewBuiltin("dir", Dir),
		"ext":    object.NewBuiltin("ext", Ext),
		"is_abs": object.NewBuiltin("is_abs", IsAbs),
		"join":   object.NewBuiltin("join", Join),
		"match":  object.NewBuiltin("match", Match),
		"split":  object.NewBuiltin("split", Split),
	})
}
