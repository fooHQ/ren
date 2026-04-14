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

const DIRENTRY = "dir_entry"

type DirEntry struct {
	value ren.DirEntry
}

func NewDirEntry(value ren.DirEntry) *DirEntry {
	return &DirEntry{
		value: value,
	}
}

func (d *DirEntry) Attrs() []object.AttrSpec {
	return dirEntryMethods.Specs()
}

func (d *DirEntry) Inspect() string {
	return fmt.Sprintf("dir_entry(name=%s, type=%s)", d.value.Name(),
		fileModeTypeString(d.value.Type()))
}

func (d *DirEntry) Type() object.Type {
	return DIRENTRY
}

func (d *DirEntry) Interface() any {
	return d.value
}

func (d *DirEntry) String() string {
	return fmt.Sprintf("dir_entry(%v)", d.value)
}

func (d *DirEntry) Value() ren.DirEntry {
	return d.value
}

func (d *DirEntry) Equals(other object.Object) bool {
	return d == other
}

func (d *DirEntry) GetAttr(name string) (object.Object, bool) {
	return dirEntryMethods.GetAttr(d, name)
}

func (d *DirEntry) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("dir_entry has no attribute %q", name)
}

func (d *DirEntry) IsTruthy() bool {
	return true
}

func (d *DirEntry) RunOperation(opType op.BinaryOpType, right object.Object) (object.Object, error) {
	return nil, object.TypeErrorf("unsupported operation for dir_entry: %v ", opType)
}

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
