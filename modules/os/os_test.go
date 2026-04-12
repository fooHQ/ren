package os_test

import (
	"context"
	"errors"
	"testing"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/stretchr/testify/require"

	"github.com/foohq/ren"
	modos "github.com/foohq/ren/modules/os"
	"github.com/foohq/ren/objects"
	"github.com/foohq/ren/testutils"
)

func TestArgs(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	m.On("Args").Return([]string{"arg1", "arg2"})

	result, err := modos.Args(ctx)
	require.NoError(t, err)
	require.IsType(t, &object.List{}, result)
	items := result.(*object.List).Value()
	require.Len(t, items, 2)
	require.Equal(t, "arg1", items[0].(*object.String).Value())
	require.Equal(t, "arg2", items[1].(*object.String).Value())
}

func TestExit(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	m.On("Exit", 0).Return()

	result, err := modos.Exit(ctx, object.NewInt(0))
	require.NoError(t, err)
	require.Equal(t, object.Nil, result)

	m.On("Exit", 1).Return()
	_, err = modos.Exit(ctx, object.NewInt(1))
	require.Error(t, err)
	require.Contains(t, err.Error(), "exited with code 1")
}

func TestChdir(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	m.On("Chdir", "/tmp").Return(nil)

	result, err := modos.Chdir(ctx, object.NewString("/tmp"))
	require.NoError(t, err)
	require.Equal(t, object.Nil, result)

	m.On("Chdir", "/nonexistent").Return(errors.New("not found"))
	_, err = modos.Chdir(ctx, object.NewString("/nonexistent"))
	require.Error(t, err)
}

func TestGetwd(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	m.On("Getwd").Return("/home/user", nil)

	result, err := modos.Getwd(ctx)
	require.NoError(t, err)
	require.Equal(t, "/home/user", result.(*object.String).Value())
}

func TestTempDir(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	m.On("TempDir").Return("/tmp")

	result, err := modos.TempDir(ctx)
	require.NoError(t, err)
	require.Equal(t, "/tmp", result.(*object.String).Value())
}

func TestGetenv(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	m.On("Getenv", "FOO").Return("bar")

	result, err := modos.Getenv(ctx, object.NewString("FOO"))
	require.NoError(t, err)
	require.Equal(t, "bar", result.(*object.String).Value())
}

func TestSetenv(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	m.On("Setenv", "FOO", "bar").Return(nil)

	result, err := modos.Setenv(ctx, object.NewString("FOO"), object.NewString("bar"))
	require.NoError(t, err)
	require.Equal(t, object.Nil, result)
}

func TestUnsetenv(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	m.On("Unsetenv", "FOO").Return(nil)

	result, err := modos.Unsetenv(ctx, object.NewString("FOO"))
	require.NoError(t, err)
	require.Equal(t, object.Nil, result)
}

func TestUserDirs(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)

	m.On("UserCacheDir").Return("/cache", nil)
	result, err := modos.UserCacheDir(ctx)
	require.NoError(t, err)
	require.Equal(t, "/cache", result.(*object.String).Value())

	m.On("UserConfigDir").Return("/config", nil)
	result, err = modos.UserConfigDir(ctx)
	require.NoError(t, err)
	require.Equal(t, "/config", result.(*object.String).Value())

	m.On("UserHomeDir").Return("/home", nil)
	result, err = modos.UserHomeDir(ctx)
	require.NoError(t, err)
	require.Equal(t, "/home", result.(*object.String).Value())
}

func TestEnviron(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	m.On("Environ").Return([]string{"A=1", "B=2"})

	result, err := modos.Environ(ctx)
	require.NoError(t, err)
	items := result.(*object.List).Value()
	require.Len(t, items, 2)
	require.Equal(t, "A=1", items[0].(*object.String).Value())
	require.Equal(t, "B=2", items[1].(*object.String).Value())
}

func TestPidsAndIds(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)

	m.On("Getpid").Return(123)
	result, err := modos.Getpid(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(123), result.(*object.Int).Value())

	m.On("Getuid").Return(456)
	result, err = modos.Getuid(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(456), result.(*object.Int).Value())
}

func TestHostname(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	m.On("Hostname").Return("localhost", nil)

	result, err := modos.Hostname(ctx)
	require.NoError(t, err)
	require.Equal(t, "localhost", result.(*object.String).Value())
}

func TestUserAndGroupLookup(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)

	user := testutils.NewMockUser("1", "2", "user1", "User One", "/home/user1")
	m.On("CurrentUser").Return(user, nil)
	result, err := modos.CurrentUser(ctx)
	require.NoError(t, err)
	require.IsType(t, &object.Map{}, result)
	require.Equal(t, "user1", result.(*object.Map).Value()["username"].(*object.String).Value())

	m.On("LookupUser", "user1").Return(user, nil)
	result, err = modos.LookupUser(ctx, object.NewString("user1"))
	require.NoError(t, err)
	require.Equal(t, "1", result.(*object.Map).Value()["uid"].(*object.String).Value())

	m.On("LookupUid", "1").Return(user, nil)
	result, err = modos.LookupUid(ctx, object.NewString("1"))
	require.NoError(t, err)
	require.Equal(t, "1", result.(*object.Map).Value()["uid"].(*object.String).Value())

	group := testutils.NewMockGroup("2", "group1")
	m.On("LookupGroup", "group1").Return(group, nil)
	result, err = modos.LookupGroup(ctx, object.NewString("group1"))
	require.NoError(t, err)
	require.Equal(t, "2", result.(*object.Map).Value()["gid"].(*object.String).Value())

	m.On("LookupGid", "2").Return(group, nil)
	result, err = modos.LookupGid(ctx, object.NewString("2"))
	require.NoError(t, err)
	require.Equal(t, "2", result.(*object.Map).Value()["gid"].(*object.String).Value())
}

func TestStdin(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)

	stdinPipe := ren.NewPipe()
	m.On("Stdin").Return(stdinPipe)
	result, err := modos.Stdin(ctx, "stdin")
	require.NoError(t, err)
	require.IsType(t, &objects.File{}, result)
	require.Equal(t, stdinPipe, result.(*objects.File).Value())
}

func TestStdout(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)

	stdoutPipe := ren.NewPipe()
	m.On("Stdout").Return(stdoutPipe)
	result, err := modos.Stdout(ctx, "stdout")
	require.NoError(t, err)
	require.IsType(t, &objects.File{}, result)
	require.Equal(t, stdoutPipe, result.(*objects.File).Value())
}
