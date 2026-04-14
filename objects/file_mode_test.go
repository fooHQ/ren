// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

package objects_test

import (
	"context"
	"io/fs"
	"testing"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/stretchr/testify/require"

	"github.com/foohq/ren"
	"github.com/foohq/ren/objects"
)

func TestFileMode(t *testing.T) {
	m := objects.NewFileMode(ren.FileMode(0644))
	require.Equal(t, object.Type(objects.FILEMODE), m.Type())
	require.Equal(t, "-rw-r--r--", m.String())
	require.Equal(t, "file_mode(-rw-r--r--)", m.Inspect())
	require.Equal(t, ren.FileMode(0644), m.Interface())
	require.Equal(t, ren.FileMode(0644), m.Value())
	require.True(t, m.IsTruthy())

	m0 := objects.NewFileMode(ren.FileMode(0))
	require.False(t, m0.IsTruthy())
}

func TestFileModeEquals(t *testing.T) {
	m1 := objects.NewFileMode(ren.FileMode(0644))
	m2 := objects.NewFileMode(ren.FileMode(0644))
	m3 := objects.NewFileMode(ren.FileMode(0755))

	require.True(t, m1.Equals(m2))
	require.False(t, m1.Equals(m3))
	require.True(t, m1.Equals(object.NewInt(0644)))
	require.False(t, m1.Equals(object.NewInt(0755)))
	require.False(t, m1.Equals(object.NewString("-rw-r--r--")))
}

func TestFileModeCompare(t *testing.T) {
	m1 := objects.NewFileMode(ren.FileMode(0644))
	m2 := objects.NewFileMode(ren.FileMode(0644))
	m3 := objects.NewFileMode(ren.FileMode(0755))

	cmp, err := m1.Compare(m2)
	require.NoError(t, err)
	require.Equal(t, 0, cmp)

	cmp, err = m1.Compare(m3)
	require.NoError(t, err)
	require.Equal(t, -1, cmp)

	cmp, err = m3.Compare(m1)
	require.NoError(t, err)
	require.Equal(t, 1, cmp)

	cmp, err = m1.Compare(object.NewInt(0644))
	require.NoError(t, err)
	require.Equal(t, 0, cmp)

	cmp, err = m1.Compare(object.NewInt(0755))
	require.NoError(t, err)
	require.Equal(t, -1, cmp)

	_, err = m1.Compare(object.NewString("-rw-r--r--"))
	require.Error(t, err)
}

func TestFileModeMethods(t *testing.T) {
	ctx := context.Background()

	// Test is_dir
	mDir := objects.NewFileMode(ren.FileMode(0755) | fs.ModeDir)
	res, ok := mDir.GetAttr("is_dir")
	require.True(t, ok, "is_dir attribute should exist")
	fn := res.(*object.Builtin)
	val, err := fn.Call(ctx)
	require.NoError(t, err)
	require.Equal(t, object.True, val)

	mFile := objects.NewFileMode(ren.FileMode(0644))
	res, ok = mFile.GetAttr("is_dir")
	require.True(t, ok)
	val, err = res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, object.False, val)

	// Test is_regular
	res, ok = mFile.GetAttr("is_regular")
	require.True(t, ok)
	val, err = res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, object.True, val)

	// Test perm
	res, ok = mFile.GetAttr("perm")
	require.True(t, ok)
	val, err = res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, object.NewString("-rw-r--r--"), val)

	// Test type
	res, ok = mFile.GetAttr("type")
	require.True(t, ok)
	val, err = res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, object.NewString("regular"), val)

	res, ok = mDir.GetAttr("type")
	require.True(t, ok)
	val, err = res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, object.NewString("dir"), val)
}

func TestFileModeMarshalJSON(t *testing.T) {
	m := objects.NewFileMode(ren.FileMode(0755) | fs.ModeDir)
	bytes, err := m.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, `{"is_dir":true,"is_regular":false,"perm":"drwxr-xr-x","type":"dir"}`, string(bytes))
}
