// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

package objects

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/deepnoodle-ai/risor/v2/pkg/op"

	"github.com/foohq/ren"
)

var _ object.Object = (*DirEntry)(nil)

// DIRENTRY is the Risor type name of a directory entry object.
const DIRENTRY = "dir_entry"

// DirEntry is a Risor object wrapping a single ren.DirEntry, as returned when
// listing a directory.
type DirEntry struct {
	value ren.DirEntry
}

// NewDirEntry wraps a directory entry as a Risor object.
func NewDirEntry(value ren.DirEntry) *DirEntry {
	return &DirEntry{
		value: value,
	}
}

// Attrs returns the attribute specifications for the directory entry's methods.
func (d *DirEntry) Attrs() []object.AttrSpec {
	return dirEntryMethods.Specs()
}

// Inspect returns a human-readable representation of the directory entry.
func (d *DirEntry) Inspect() string {
	return fmt.Sprintf("dir_entry(name=%s, type=%s)", d.value.Name(),
		fileModeTypeString(d.value.Type()))
}

// Type returns the Risor type name of the directory entry.
func (d *DirEntry) Type() object.Type {
	return DIRENTRY
}

// Interface returns the underlying ren.DirEntry.
func (d *DirEntry) Interface() any {
	return d.value
}

// String returns a string representation of the directory entry.
func (d *DirEntry) String() string {
	return fmt.Sprintf("dir_entry(%v)", d.value)
}

// Value returns the underlying ren.DirEntry.
func (d *DirEntry) Value() ren.DirEntry {
	return d.value
}

// Equals reports whether other is the same directory entry instance.
func (d *DirEntry) Equals(other object.Object) bool {
	return d == other
}

// GetAttr returns the named method of the directory entry.
func (d *DirEntry) GetAttr(name string) (object.Object, bool) {
	return dirEntryMethods.GetAttr(d, name)
}

// SetAttr always returns an error; dir_entry attributes are read-only.
func (d *DirEntry) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("dir_entry has no attribute %q", name)
}

// IsTruthy reports whether the directory entry is truthy; it is always true.
func (d *DirEntry) IsTruthy() bool {
	return true
}

// RunOperation always returns an error; dir_entry supports no binary operations.
func (d *DirEntry) RunOperation(opType op.BinaryOpType, right object.Object) (object.Object, error) {
	return nil, object.TypeErrorf("unsupported operation for dir_entry: %v ", opType)
}

// MarshalJSON encodes the directory entry, including its name, type, whether it
// is a directory, and its file info.
func (d *DirEntry) MarshalJSON() ([]byte, error) {
	info, _ := d.value.Info()
	return json.Marshal(struct {
		Name  string    `json:"name"`
		Type  string    `json:"type"`
		IsDir bool      `json:"is_dir"`
		Info  *FileInfo `json:"info"`
	}{
		Name:  d.value.Name(),
		Type:  fileModeTypeString(d.value.Type()),
		IsDir: d.value.IsDir(),
		Info:  NewFileInfo(info),
	})
}

// dirEntryMethods holds the methods exposed on dir_entry objects (name, info).
var dirEntryMethods = object.NewMethodRegistry[*DirEntry](DIRENTRY)

func init() {
	dirEntryMethods.Define("name").
		Doc(""). // TODO
		Returns("string").
		Impl(func(d *DirEntry, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("dir_entry.name", 0, len(args))
			}
			return object.NewString(d.value.Name()), nil
		})
	dirEntryMethods.Define("info").
		Doc(""). // TODO
		Returns(FILEINFO).
		Impl(func(d *DirEntry, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("dir_entry.info", 0, len(args))
			}
			info, err := d.value.Info()
			if err != nil {
				return nil, err
			}
			return NewFileInfo(info), nil
		})
}
