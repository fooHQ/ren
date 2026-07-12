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

// FILEINFO is the Risor type name of a file info object.
const FILEINFO = "file_info"

// FileInfo is a Risor object wrapping a ren.FileInfo, describing a file's name,
// size, mode, and modification time.
type FileInfo struct {
	value ren.FileInfo
}

// NewFileInfo wraps file information as a Risor object.
func NewFileInfo(value ren.FileInfo) *FileInfo {
	return &FileInfo{
		value: value,
	}
}

// Attrs returns the attribute specifications for the file info's methods.
func (f *FileInfo) Attrs() []object.AttrSpec {
	return fileInfoMethods.Specs()
}

// Inspect returns a human-readable representation of the file info.
func (f *FileInfo) Inspect() string {
	return f.String()
}

// Type returns the Risor type name of the file info.
func (f *FileInfo) Type() object.Type {
	return FILEINFO
}

// Interface returns the underlying ren.FileInfo.
func (f *FileInfo) Interface() any {
	return f.value
}

// String returns a string representation of the file info.
func (f *FileInfo) String() string {
	v := f.value
	return fmt.Sprintf("file_info(name=%s, mode=%s, size=%d, mod_time=%v)",
		v.Name(), v.Mode().String(), v.Size(), v.ModTime().Format(time.RFC3339))
}

// Value returns the underlying ren.FileInfo.
func (f *FileInfo) Value() ren.FileInfo {
	return f.value
}

// Equals reports whether other is the same file_info instance.
func (f *FileInfo) Equals(other object.Object) bool {
	return f == other
}

// GetAttr returns the named method of the file info.
func (f *FileInfo) GetAttr(name string) (object.Object, bool) {
	return fileInfoMethods.GetAttr(f, name)
}

// SetAttr always returns an error; file_info attributes are read-only.
func (f *FileInfo) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("file_info has no attribute %q", name)
}

// IsTruthy reports whether the file info is truthy; it is always true.
func (f *FileInfo) IsTruthy() bool {
	return true
}

// RunOperation always returns an error; file_info supports no binary operations.
func (f *FileInfo) RunOperation(opType op.BinaryOpType, right object.Object) (object.Object, error) {
	return nil, object.TypeErrorf("unsupported operation for file_info: %v", opType)
}

// MarshalJSON encodes the file info, including its name, size, mode,
// modification time, and whether it is a directory.
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

// fileInfoMethods holds the methods exposed on file_info objects (name, size,
// mod_time, mode).
var fileInfoMethods = object.NewMethodRegistry[*FileInfo](FILEINFO)

func init() {
	fileInfoMethods.Define("name").
		Doc(""). // TODO
		Returns("string").
		Impl(func(f *FileInfo, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file_info.name", 0, len(args))
			}
			return object.NewString(f.value.Name()), nil
		})
	fileInfoMethods.Define("size").
		Doc(""). // TODO
		Returns("int").
		Impl(func(f *FileInfo, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file_info.size", 0, len(args))
			}
			return object.NewInt(f.value.Size()), nil
		})
	fileInfoMethods.Define("mod_time").
		Doc(""). // TODO
		Returns("time").
		Impl(func(f *FileInfo, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file_info.mod_time", 0, len(args))
			}
			return object.NewTime(f.value.ModTime()), nil
		})
	fileInfoMethods.Define("mode").
		Doc(""). // TODO
		Returns(FILEMODE).
		Impl(func(f *FileInfo, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file_info.mode", 0, len(args))
			}
			return NewFileMode(f.value.Mode()), nil
		})
}
