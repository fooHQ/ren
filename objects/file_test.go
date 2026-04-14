// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

package objects_test

import (
	"context"
	"io"
	"testing"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/foohq/ren/objects"
	"github.com/foohq/ren/testutils"
)

// Define a type that implements ren.File and io.Seeker for testing
type mockSeekableFile struct {
	testutils.MockFile
}

func (m *mockSeekableFile) Seek(offset int64, whence int) (int64, error) {
	args := m.Called(offset, whence)
	return int64(args.Int(0)), args.Error(1)
}

func TestFile(t *testing.T) {
	ctx := context.Background()
	m := &testutils.MockFile{}
	f := objects.NewFile(ctx, m, "/test.txt")

	require.Equal(t, object.Type(objects.FILE), f.Type())
	require.Equal(t, m, f.Interface())
	require.Equal(t, m, f.Value())
	require.True(t, f.IsTruthy())
	require.Equal(t, "file(path=/test.txt)", f.Inspect())
}

func TestFileMethods(t *testing.T) {
	ctx := context.Background()
	m := &testutils.MockFile{}
	f := objects.NewFile(ctx, m, "/test.txt")

	// Test name()
	res, ok := f.GetAttr("name")
	require.True(t, ok)
	val, err := res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, object.NewString("/test.txt"), val)

	// Test info()
	fiMock := &testutils.MockFileInfo{}
	m.On("Stat").Return(fiMock, nil)
	res, ok = f.GetAttr("info")
	require.True(t, ok)
	val, err = res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, objects.NewFileInfo(fiMock), val)

	// Test read()
	m.On("Read", mock.Anything).Return(5, nil).Run(func(args mock.Arguments) {
		p := args.Get(0).([]byte)
		copy(p, "hello")
	})
	res, ok = f.GetAttr("read")
	require.True(t, ok)
	fn := res.(*object.Builtin)
	buf := object.NewBytes(make([]byte, 5))
	val, err = fn.Call(ctx, buf)
	require.NoError(t, err)
	require.Equal(t, object.NewBytes([]byte("hello")), val)

	// Test write()
	m.On("Write", []byte("world")).Return(5, nil)
	res, ok = f.GetAttr("write")
	require.True(t, ok)
	val, err = res.(*object.Builtin).Call(ctx, object.NewBytes([]byte("world")))
	require.NoError(t, err)
	require.Equal(t, object.NewInt(5), val)

	// Test close()
	m.On("Close").Return(nil)
	res, ok = f.GetAttr("close")
	require.True(t, ok)
	val, err = res.(*object.Builtin).Call(ctx)
	require.NoError(t, err)
	require.Equal(t, object.Nil, val)
	// Test seek() on non-seeker
	f2 := objects.NewFile(ctx, m, "/test.txt")
	res, ok = f2.GetAttr("seek")
	require.True(t, ok)
	_, err = res.(*object.Builtin).Call(ctx, object.NewInt(0), object.NewInt(0))
	require.Error(t, err)
	require.Contains(t, err.Error(), "does not support seeking")
}

func TestFileSeek(t *testing.T) {
	ctx := context.Background()
	m := &mockSeekableFile{}
	f := objects.NewFile(ctx, m, "/test.txt")

	m.On("Seek", int64(10), io.SeekStart).Return(10, nil)
	res, ok := f.GetAttr("seek")
	require.True(t, ok)
	val, err := res.(*object.Builtin).Call(ctx, object.NewInt(10), object.NewInt(io.SeekStart))
	require.NoError(t, err)
	require.Equal(t, object.NewInt(10), val)
}

func TestFileEquals(t *testing.T) {
	ctx := context.Background()
	m1 := &testutils.MockFile{}
	m2 := &testutils.MockFile{}
	f1 := objects.NewFile(ctx, m1, "/test.txt")
	f2 := objects.NewFile(ctx, m1, "/test.txt")
	f3 := objects.NewFile(ctx, m2, "/other.txt")

	require.True(t, f1.Equals(f1))
	require.False(t, f1.Equals(f2)) // Currently f1 == other, which means pointer equality
	require.False(t, f1.Equals(f3))
}

func TestFileMarshalJSON(t *testing.T) {
	ctx := context.Background()
	m := &testutils.MockFile{}
	f := objects.NewFile(ctx, m, "/test.txt")
	_, err := f.MarshalJSON()
	require.Error(t, err)
	require.Equal(t, "type error: unable to marshal file", err.Error())
}
