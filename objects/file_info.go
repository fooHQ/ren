// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

package objects

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/deepnoodle-ai/risor/v2/pkg/op"

	"github.com/foohq/ren"
)

var _ object.Object = (*FileInfo)(nil)

const FILEINFO = "file_info"

type FileInfo struct {
	value ren.FileInfo
}

func NewFileInfo(value ren.FileInfo) *FileInfo {
	return &FileInfo{
		value: value,
	}
}

func (f *FileInfo) Attrs() []object.AttrSpec {
	return fileInfoMethods.Specs()
}

func (f *FileInfo) Inspect() string {
	return f.String()
}

func (f *FileInfo) Type() object.Type {
	return FILEINFO
}

func (f *FileInfo) Interface() interface{} {
	return f.value
}

func (f *FileInfo) String() string {
	v := f.value
	return fmt.Sprintf("file_info(name=%s, mode=%s, size=%d, mod_time=%v)",
		v.Name(), v.Mode().String(), v.Size(), v.ModTime().Format(time.RFC3339))
}

func (f *FileInfo) Value() ren.FileInfo {
	return f.value
}

func (f *FileInfo) Equals(other object.Object) bool {
	return f == other
}

func (f *FileInfo) GetAttr(name string) (object.Object, bool) {
	return fileInfoMethods.GetAttr(f, name)
}

func (f *FileInfo) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("file_info has no attribute %q", name)
}

func (f *FileInfo) IsTruthy() bool {
	return true
}

func (f *FileInfo) RunOperation(opType op.BinaryOpType, right object.Object) (object.Object, error) {
	return nil, object.TypeErrorf("unsupported operation for file_info: %v", opType)
}

func (f *FileInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IsDir   bool      `json:"is_dir"`
		Mode    *FileMode `json:"mode"`
		ModTime time.Time `json:"mod_time"`
		Name    string    `json:"name"`
		Size    int64     `json:"size"`
	}{
		IsDir:   f.value.IsDir(),
		Mode:    NewFileMode(f.value.Mode()),
		ModTime: f.value.ModTime(),
		Name:    f.value.Name(),
		Size:    f.value.Size(),
	})
}

var fileInfoMethods = object.NewMethodRegistry[*FileInfo](FILEINFO)

func init() {
	fileInfoMethods.Define("name").
		Doc(""). // TODO
		Impl(func(f *FileInfo, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file_info.name", 0, len(args))
			}
			return object.NewString(f.value.Name()), nil
		})
	fileInfoMethods.Define("size").
		Doc(""). // TODO
		Impl(func(f *FileInfo, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file_info.size", 0, len(args))
			}
			return object.NewInt(f.value.Size()), nil
		})
	fileInfoMethods.Define("mod_time").
		Doc(""). // TODO
		Impl(func(f *FileInfo, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file_info.mod_time", 0, len(args))
			}
			return object.NewTime(f.value.ModTime()), nil
		})
	fileInfoMethods.Define("mode").
		Doc(""). // TODO
		Impl(func(f *FileInfo, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file_info.mode", 0, len(args))
			}
			return NewFileMode(f.value.Mode()), nil
		})
}
