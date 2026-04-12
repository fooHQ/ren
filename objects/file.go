// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

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

const FILE = "file"

type File struct {
	ctx    context.Context
	value  ren.File
	path   string
	once   sync.Once
	closed chan bool
}

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

func (f *File) Attrs() []object.AttrSpec {
	return fileMethods.Specs()
}

func (f *File) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("file has no attribute %q", name)
}

func (f *File) IsTruthy() bool {
	return true
}

func (f *File) Inspect() string {
	return fmt.Sprintf("file(path=%s)", f.path)
}

func (f *File) Type() object.Type {
	return FILE
}

func (f *File) GetAttr(name string) (object.Object, bool) {
	return fileMethods.GetAttr(f, name)
}

func (f *File) cleanup() {
	go func() {
		select {
		case <-f.closed:
		case <-f.ctx.Done():
			_ = f.value.Close()
		}
	}()
}

func (f *File) Interface() interface{} {
	return f.value
}

func (f *File) Value() ren.File {
	return f.value
}

func (f *File) String() string {
	return f.Inspect()
}

func (f *File) Equals(other object.Object) bool {
	return f == other
}

func (f *File) RunOperation(opType op.BinaryOpType, right object.Object) (object.Object, error) {
	return nil, object.TypeErrorf("unsupported operation for file: %v ", opType)
}

func (f *File) MarshalJSON() ([]byte, error) {
	return nil, object.TypeErrorf("unable to marshal file")
}

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
				if n == len(slice) {
					return obj, nil
				}
				return object.NewBytes(slice[:n]), nil
			default:
				return nil, object.TypeErrorf("file.read() expected bytes (%s given)", obj.Type())
			}
		})
	fileMethods.Define("write").
		Doc(""). // TODO
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
