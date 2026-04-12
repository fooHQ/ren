package fs_test

import (
	"context"
	"os"
	"testing"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/stretchr/testify/require"

	"github.com/foohq/ren"
	modfs "github.com/foohq/ren/modules/fs"
	"github.com/foohq/ren/objects"
	"github.com/foohq/ren/testutils"
)

func TestOpenFile(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	path := "test.txt"
	mode := "r"
	perm := 0644
	f := &testutils.MockFile{}

	m.On("OpenFile", path, os.O_RDONLY, ren.FileMode(perm)).Return(f, nil)

	result, err := modfs.OpenFile(ctx, object.NewString(path), object.NewString(mode), object.NewInt(int64(perm)))
	require.NoError(t, err)
	require.IsType(t, &objects.File{}, result)
	require.Equal(t, f, result.(*objects.File).Value())
}

func TestReadFile(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	path := "test.txt"
	data := []byte("hello world")

	m.On("ReadFile", path).Return(data, nil)

	result, err := modfs.ReadFile(ctx, object.NewString(path))
	require.NoError(t, err)
	require.Equal(t, data, result.(*object.Bytes).Value())
}

func TestReadDir(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	path := "testdir"
	entry1 := &testutils.MockDirEntry{}
	entry1.On("Name").Return("file1")

	m.On("ReadDir", path).Return([]ren.DirEntry{entry1}, nil)

	result, err := modfs.ReadDir(ctx, object.NewString(path))
	require.NoError(t, err)
	items := result.(*object.List).Value()
	require.Len(t, items, 1)
	require.Equal(t, "file1", items[0].(*objects.DirEntry).Value().Name())
}

func TestWriteFile(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	path := "test.txt"
	data := "hello world"
	perm := 0644

	m.On("WriteFile", path, []byte(data), ren.FileMode(perm)).Return(nil)

	result, err := modfs.WriteFile(ctx, object.NewString(path), object.NewString(data), object.NewInt(int64(perm)))
	require.NoError(t, err)
	require.Equal(t, object.Nil, result)
}

func TestRemove(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	path := "test.txt"

	m.On("Remove", path).Return(nil)

	result, err := modfs.Remove(ctx, object.NewString(path))
	require.NoError(t, err)
	require.Equal(t, object.Nil, result)
}

func TestRemoveAll(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	path := "testdir"

	m.On("RemoveAll", path).Return(nil)

	result, err := modfs.RemoveAll(ctx, object.NewString(path))
	require.NoError(t, err)
	require.Equal(t, object.Nil, result)
}

func TestRename(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	oldpath := "old.txt"
	newpath := "new.txt"

	m.On("Rename", oldpath, newpath).Return(nil)

	result, err := modfs.Rename(ctx, object.NewString(oldpath), object.NewString(newpath))
	require.NoError(t, err)
	require.Equal(t, object.Nil, result)
}

func TestStat(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	path := "test.txt"
	info := &testutils.MockFileInfo{}
	info.On("Name").Return("test.txt")

	m.On("Stat", path).Return(info, nil)

	result, err := modfs.Stat(ctx, object.NewString(path))
	require.NoError(t, err)
	require.IsType(t, &objects.FileInfo{}, result)
	require.Equal(t, info, result.(*objects.FileInfo).Value())
}

func TestMkdirAll(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	path := "test/dir"
	perm := 0755

	m.On("MkdirAll", path, ren.FileMode(perm)).Return(nil)

	result, err := modfs.MkdirAll(ctx, object.NewString(path), object.NewInt(int64(perm)))
	require.NoError(t, err)
	require.Equal(t, object.Nil, result)
}

func TestMkdir(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	path := "testdir"
	perm := 0755

	m.On("Mkdir", path, ren.FileMode(perm)).Return(nil)

	result, err := modfs.Mkdir(ctx, object.NewString(path), object.NewInt(int64(perm)))
	require.NoError(t, err)
	require.Equal(t, object.Nil, result)
}

func TestMkdirTemp(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	dir := "/tmp"
	pattern := "test-*"
	tempDir := "/tmp/test-123"

	m.On("MkdirTemp", dir, pattern).Return(tempDir, nil)

	result, err := modfs.MkdirTemp(ctx, object.NewString(dir), object.NewString(pattern))
	require.NoError(t, err)
	require.Equal(t, tempDir, result.(*object.String).Value())
}

func TestSymlink(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	oldname := "old.txt"
	newname := "new.txt"

	m.On("Symlink", oldname, newname).Return(nil)

	result, err := modfs.Symlink(ctx, object.NewString(oldname), object.NewString(newname))
	require.NoError(t, err)
	require.Equal(t, object.Nil, result)
}
