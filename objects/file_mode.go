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

// FILEMODE is the Risor type name of a file mode object.
const FILEMODE = "file_mode"

var _ object.Object = (*FileMode)(nil)

// FileMode is a Risor object wrapping a ren.FileMode. It carries a file's type
// and permission bits and can be compared against another file_mode or an int.
type FileMode struct {
	value ren.FileMode
}

// NewFileMode wraps a file mode as a Risor object.
func NewFileMode(value ren.FileMode) *FileMode {
	return &FileMode{
		value: value,
	}
}

// Attrs returns the attribute specifications for the file mode's methods.
func (m *FileMode) Attrs() []object.AttrSpec {
	return fileModeMethods.Specs()
}

// Inspect returns a human-readable representation of the file mode.
func (m *FileMode) Inspect() string {
	return fmt.Sprintf("file_mode(%s)", m.value)
}

// Type returns the Risor type name of the file mode.
func (m *FileMode) Type() object.Type {
	return FILEMODE
}

// Interface returns the underlying ren.FileMode.
func (m *FileMode) Interface() any {
	return m.value
}

// String returns the standard textual representation of the file mode.
func (m *FileMode) String() string {
	return m.value.String()
}

// Value returns the underlying ren.FileMode.
func (m *FileMode) Value() ren.FileMode {
	return m.value
}

// Compare orders the file mode against another file_mode or int, returning -1,
// 0, or 1. It errors for any other type.
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

// Equals reports whether the file mode equals another file_mode or int of the
// same value.
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

// GetAttr returns the named method of the file mode.
func (m *FileMode) GetAttr(name string) (object.Object, bool) {
	return fileModeMethods.GetAttr(m, name)
}

// SetAttr always returns an error; file_mode attributes are read-only.
func (m *FileMode) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("file_mode has no attribute %q", name)
}

// IsTruthy reports whether the file mode is non-zero.
func (m *FileMode) IsTruthy() bool {
	return m.value != ren.FileMode(0)
}

// RunOperation always returns an error; file_mode supports no binary operations.
func (m *FileMode) RunOperation(opType op.BinaryOpType, right object.Object) (object.Object, error) {
	return nil, object.TypeErrorf("unsupported operation for file_mode: %v", opType)
}

// MarshalJSON encodes the file mode, including whether it is a directory or
// regular file, its permission string, and its type.
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

// fileModeTypeString returns a short name for a file mode's type, such as
// "dir", "regular", "symlink", or "unknown".
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

// fileModeMethods holds the methods exposed on file_mode objects (is_dir,
// is_regular, perm, type).
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
