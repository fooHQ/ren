// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

package objects

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/deepnoodle-ai/risor/v2/pkg/op"

	"github.com/foohq/ren"
)

const FILEMODE = "file_mode"

var _ object.Object = (*FileMode)(nil)

type FileMode struct {
	value ren.FileMode
}

func NewFileMode(value ren.FileMode) *FileMode {
	return &FileMode{
		value: value,
	}
}

func (m *FileMode) Attrs() []object.AttrSpec {
	return fileModeMethods.Specs()
}

func (m *FileMode) Inspect() string {
	return fmt.Sprintf("file_mode(%s)", m.value)
}

func (m *FileMode) Type() object.Type {
	return FILEMODE
}

func (m *FileMode) Interface() any {
	return m.value
}

func (m *FileMode) String() string {
	return m.value.String()
}

func (m *FileMode) Value() ren.FileMode {
	return m.value
}

func (m *FileMode) Compare(other object.Object) (int, error) {
	switch other := other.(type) {
	case *FileMode:
		if m.value < other.value {
			return -1, nil
		} else if m.value > other.value {
			return 1, nil
		}
		return 0, nil
	case *object.Int:
		if int64(m.value) < other.Value() {
			return -1, nil
		} else if int64(m.value) > other.Value() {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, object.TypeErrorf("unable to compare file_mode and %s", other.Type())
	}
}

func (m *FileMode) Equals(other object.Object) bool {
	switch other := other.(type) {
	case *FileMode:
		if m.value == other.value {
			return true
		}
	case *object.Int:
		if int64(m.value) == other.Value() {
			return true
		}
	}
	return false
}

func (m *FileMode) GetAttr(name string) (object.Object, bool) {
	return fileModeMethods.GetAttr(m, name)
}

func (m *FileMode) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("file_mode has no attribute %q", name)
}

func (m *FileMode) IsTruthy() bool {
	return m.value != ren.FileMode(0)
}

func (m *FileMode) RunOperation(opType op.BinaryOpType, right object.Object) (object.Object, error) {
	return nil, object.TypeErrorf("unsupported operation for file_mode: %v", opType)
}

func (m *FileMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IsDir     bool   `json:"is_dir"`
		IsRegular bool   `json:"is_regular"`
		Perm      string `json:"perm"`
		Type      string `json:"type"`
	}{
		IsDir:     m.value.IsDir(),
		IsRegular: m.value.IsRegular(),
		Perm:      m.value.String(),
		Type:      fileModeTypeString(m.value),
	})
}

func fileModeTypeString(m ren.FileMode) string {
	switch {
	case m.IsDir():
		return "dir"
	case m.IsRegular():
		return "regular"
	case m&fs.ModeSymlink != 0:
		return "symlink"
	case m&fs.ModeNamedPipe != 0:
		return "named_pipe"
	case m&fs.ModeSocket != 0:
		return "socket"
	case m&fs.ModeDevice != 0:
		return "device"
	case m&fs.ModeCharDevice != 0:
		return "char_device"
	case m&fs.ModeIrregular != 0:
		return "irregular"
	default:
		return "unknown"
	}
}

var fileModeMethods = object.NewMethodRegistry[*FileMode](FILEMODE)

func init() {
	fileModeMethods.Define("is_dir").
		Doc(""). // TODO
		Returns("bool").
		Impl(func(f *FileMode, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file_mode.is_dir", 0, len(args))
			}
			return object.NewBool(f.value.IsDir()), nil
		})
	fileModeMethods.Define("is_regular").
		Doc(""). // TODO
		Returns("bool").
		Impl(func(f *FileMode, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file_mode.is_regular", 0, len(args))
			}
			return object.NewBool(f.value.IsRegular()), nil
		})
	fileModeMethods.Define("perm").
		Doc(""). // TODO
		Returns("string").
		Impl(func(f *FileMode, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file_mode.perm", 0, len(args))
			}
			return object.NewString(f.value.String()), nil
		})
	fileModeMethods.Define("type").
		Doc(""). // TODO
		Returns("string").
		Impl(func(f *FileMode, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("file_mode.type", 0, len(args))
			}
			return object.NewString(fileModeTypeString(f.value)), nil
		})
}
