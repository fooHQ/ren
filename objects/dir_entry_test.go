// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

package objects_test

import (
	"context"
	"io/fs"
	"testing"
	"time"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/stretchr/testify/require"

	"github.com/foohq/ren"
	"github.com/foohq/ren/objects"
	"github.com/foohq/ren/testutils"
)

func TestDirEntry(t *testing.T) {
	m := &testutils.MockDirEntry{}
	m.On("Name").Return("test.txt")
	m.On("Type").Return(ren.FileMode(0644))
	m.On("IsDir").Return(false)

	de := objects.NewDirEntry(m)
	require.Equal(t, object.Type(objects.DIRENTRY), de.Type())
	require.Equal(t, m, de.Interface())
	require.Equal(t, m, de.Value())
	require.True(t, de.IsTruthy())
	require.Equal(t, "dir_entry(name=test.txt, type=regular)", de.Inspect())
}

func TestDirEntryMethods(t *testing.T) {
	ctx := context.Background()
	m := &testutils.MockDirEntry{}
	m.On("Name").Return("test.txt")

	fiMock := &testutils.MockFileInfo{}
	m.On("Info").Return(fiMock, nil)

	de := objects.NewDirEntry(m)

	// Test name()
	res, ok := de.GetAttr("name")
	require.True(t, ok)
	val, err := res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, object.NewString("test.txt"), val)

	// Test info()
	res, ok = de.GetAttr("info")
	require.True(t, ok)
	val, err = res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, objects.NewFileInfo(fiMock), val)
}

func TestDirEntryMarshalJSON(t *testing.T) {
	m := &testutils.MockDirEntry{}
	m.On("Name").Return("test_dir")
	m.On("Type").Return(ren.FileMode(0755) | fs.ModeDir)
	m.On("IsDir").Return(true)

	fiMock := &testutils.MockFileInfo{}
	fiMock.On("Name").Return("test_dir")
	fiMock.On("Size").Return(4096)
	fiMock.On("Mode").Return(ren.FileMode(0755) | fs.ModeDir)
	fiMock.On("ModTime").Return(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	fiMock.On("IsDir").Return(true)
	m.On("Info").Return(fiMock, nil)

	de := objects.NewDirEntry(m)
	bytes, err := de.MarshalJSON()
	require.NoError(t, err)

	expectedJSON := `{
		"name": "test_dir",
		"type": "dir",
		"is_dir": true,
		"info": {
			"is_dir": true,
			"mode": {
				"is_dir": true,
				"is_regular": false,
				"perm": "drwxr-xr-x",
				"type": "dir"
			},
			"mod_time": "2023-01-01T00:00:00Z",
			"name": "test_dir",
			"size": 4096
		}
	}`
	require.JSONEq(t, expectedJSON, string(bytes))
}
