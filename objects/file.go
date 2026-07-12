// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

// Package objects provides the Risor object types that Ren's built-in modules
// expose to scripts, wrapping Go filesystem values such as files, directory
// entries, file info, and file modes.
package objects

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/deepnoodle-ai/risor/v2/pkg/op"

	"github.com/foohq/ren"
)

var _ object.Object = (*File)(nil)

// FILE is the Risor type name of a file object.
const FILE = "file"

// File is a Risor object wrapping an open ren.File. It is closed automatically
// when its context is done, unless the script closes it first.
type File struct {
	ctx    context.Context
	value  ren.File
	path   string
	once   sync.Once
	closed chan bool
}

// NewFile wraps an open file at the given path as a Risor object and starts a
// goroutine that closes it when the context is done.
func NewFile(ctx context.Context, value ren.File, path string) *File {
	f := &File{
		ctx:    ctx,
		value:  value,
		path:   path,
		closed: make(chan bool),
	}
	f.cleanup()
	return f
}

// Attrs returns the attribute specifications for the file's methods.
func (f *File) Attrs() []object.AttrSpec {
	return fileMethods.Specs()
}

// SetAttr always returns an error; file attributes are read-only.
func (f *File) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("file has no attribute %q", name)
}

// IsTruthy reports whether the file is truthy; it is always true.
func (f *File) IsTruthy() bool {
	return true
}

// Inspect returns a human-readable representation of the file.
func (f *File) Inspect() string {
	return fmt.Sprintf("file(path=%s)", f.path)
}

// Type returns the Risor type name of the file.
func (f *File) Type() object.Type {
	return FILE
}

// GetAttr returns the named method of the file.
func (f *File) GetAttr(name string) (object.Object, bool) {
	return fileMethods.GetAttr(f, name)
}

// cleanup closes the wrapped file when the context is done, unless it has
// already been closed.
func (f *File) cleanup() {
	go func() {
		select {
		case <-f.closed:
		case <-f.ctx.Done():
			_ = f.value.Close()
		}
	}()
}

// Interface returns the underlying ren.File.
func (f *File) Interface() any {
	return f.value
}

// Value returns the underlying ren.File.
func (f *File) Value() ren.File {
	return f.value
}

// String returns a string representation of the file.
func (f *File) String() string {
	return f.Inspect()
}

// Equals reports whether other is the same file instance.
func (f *File) Equals(other object.Object) bool {
	return f == other
}

// RunOperation always returns an error; files support no binary operations.
func (f *File) RunOperation(opType op.BinaryOpType, right object.Object) (object.Object, error) {
	return nil, object.TypeErrorf("unsupported operation for file: %v ", opType)
}

// MarshalJSON always returns an error; files cannot be marshalled to JSON.
func (f *File) MarshalJSON() ([]byte, error) {
	return nil, object.TypeErrorf("unable to marshal file")
}

// fileMethods holds the methods exposed on file objects (name, info, read,
// write, close, seek).
var fileMethods = object.NewMethodRegistry[*File](FILE)

func init() {
	fileMethods.Define("name").
		Doc("Return the name of the file.").
		Returns("string").
		Impl(func(f *File, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file.name", 0, len(args))
			}
			return object.NewString(f.path), nil
		})
	fileMethods.Define("info").
		Doc(""). // TODO
		Returns(FILEINFO).
		Impl(func(f *File, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file.info", 0, len(args))
			}
			info, err := f.value.Stat()
			if err != nil {
				return nil, err
			}
			return NewFileInfo(info), nil
		})
	fileMethods.Define("read").
		Doc(""). // TODO
		Arg("buffer").
		Returns("int").
		Impl(func(f *File, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, object.NewArgsError("file.read", 1, len(args))
			}
			switch obj := args[0].(type) {
			case *object.Bytes:
				slice := obj.Value()
				n, ioErr := f.value.Read(slice)
				if ioErr != nil && ioErr != io.EOF {
					return nil, object.NewError(ioErr)
				}
				return object.NewInt(int64(n)), nil
			default:
				return nil, object.TypeErrorf("file.read() expected bytes (%s given)", obj.Type())
			}
		})
	fileMethods.Define("write").
		Doc("").
		Arg("data").
		Returns("int").
		Impl(func(f *File, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, object.NewArgsError("file.write", 1, len(args))
			}
			b, err := object.AsBytes(args[0])
			if err != nil {
				return nil, err
			}
			writer, ok := f.value.(io.Writer)
			if !ok {
				return nil, object.TypeErrorf("this file does not support writing")
			}
			n, ioErr := writer.Write(b)
			if ioErr != nil {
				return nil, object.NewError(ioErr)
			}
			return object.NewInt(int64(n)), nil
		})
	fileMethods.Define("close").
		Doc(""). // TODO
		Impl(func(f *File, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file.close", 0, len(args))
			}
			var err error
			f.once.Do(func() {
				err = f.value.Close()
				close(f.closed)
			})
			if err != nil {
				return nil, object.NewError(err)
			}
			return object.Nil, nil
		})
	fileMethods.Define("seek").
		Doc(""). // TODO
		Args("offset", "whence").
		Returns("int").
		Impl(func(f *File, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 2 {
				return nil, object.NewArgsError("file.seek", 2, len(args))
			}
			offset, err := object.AsInt(args[0])
			if err != nil {
				return nil, err
			}
			whence, err := object.AsInt(args[1])
			if err != nil {
				return nil, err
			}
			seeker, ok := f.value.(io.Seeker)
			if !ok {
				return nil, object.TypeErrorf("this file does not support seeking")
			}
			newPosition, ioErr := seeker.Seek(offset, int(whence))
			if ioErr != nil {
				return nil, object.NewError(ioErr)
			}
			return object.NewInt(newPosition), nil
		})
}
