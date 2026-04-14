// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

package objects_test

import (
	"context"
	"testing"
	"time"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/stretchr/testify/require"

	"github.com/foohq/ren"
	"github.com/foohq/ren/objects"
	"github.com/foohq/ren/testutils"
)

func TestFileInfo(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	m := &testutils.MockFileInfo{}
	m.On("Name").Return("test.txt")
	m.On("Size").Return(100)
	m.On("Mode").Return(ren.FileMode(0644))
	m.On("ModTime").Return(now)
	m.On("IsDir").Return(false)

	fi := objects.NewFileInfo(m)
	require.Equal(t, object.Type(objects.FILEINFO), fi.Type())
	require.Equal(t, m, fi.Interface())
	require.Equal(t, m, fi.Value())
	require.True(t, fi.IsTruthy())

	expectedString := "file_info(name=test.txt, mode=-rw-r--r--, size=100, mod_time=" + now.Format(time.RFC3339) + ")"
	require.Equal(t, expectedString, fi.String())
	require.Equal(t, expectedString, fi.Inspect())
}

func TestFileInfoMethods(t *testing.T) {
	ctx := context.Background()
	now := time.Now().Truncate(time.Second)
	m := &testutils.MockFileInfo{}
	m.On("Name").Return("test.txt")
	m.On("Size").Return(100)
	m.On("Mode").Return(ren.FileMode(0644))
	m.On("ModTime").Return(now)

	fi := objects.NewFileInfo(m)

	// Test name()
	res, ok := fi.GetAttr("name")
	require.True(t, ok)
	val, err := res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, object.NewString("test.txt"), val)

	// Test size()
	res, ok = fi.GetAttr("size")
	require.True(t, ok)
	val, err = res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, object.NewInt(100), val)

	// Test mod_time()
	res, ok = fi.GetAttr("mod_time")
	require.True(t, ok)
	val, err = res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, object.NewTime(now), val)

	// Test mode()
	res, ok = fi.GetAttr("mode")
	require.True(t, ok)
	val, err = res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, objects.NewFileMode(ren.FileMode(0644)), val)
}

func TestFileInfoMarshalJSON(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	m := &testutils.MockFileInfo{}
	m.On("Name").Return("test.txt")
	m.On("Size").Return(100)
	m.On("Mode").Return(ren.FileMode(0644))
	m.On("ModTime").Return(now)
	m.On("IsDir").Return(false)

	fi := objects.NewFileInfo(m)
	bytes, err := fi.MarshalJSON()
	require.NoError(t, err)

	expectedJSON := `{
		"is_dir": false,
		"mode": {
			"is_dir": false,
			"is_regular": true,
			"perm": "-rw-r--r--",
			"type": "regular"
		},
		"mod_time": "` + now.Format(time.RFC3339) + `",
		"name": "test.txt",
		"size": 100
	}`
	require.JSONEq(t, expectedJSON, string(bytes))
}
