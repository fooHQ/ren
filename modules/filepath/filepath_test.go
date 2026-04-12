package filepath_test

import (
	"context"
	"testing"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/stretchr/testify/require"

	"github.com/foohq/ren"
	modfilepath "github.com/foohq/ren/modules/filepath"
	"github.com/foohq/ren/testutils"
)

func TestAbs(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	abs, err := modfilepath.Abs(ctx, object.NewString("foo"), object.NewString("/"))
	require.NoError(t, err)
	require.IsType(t, &object.String{}, abs)
	require.Equal(t, "/foo", abs.(*object.String).Value())
}

func TestBase(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	base, err := modfilepath.Base(ctx, object.NewString("/foo/bar.txt"))
	require.NoError(t, err)
	require.IsType(t, &object.String{}, base)
	require.Equal(t, "bar.txt", base.(*object.String).Value())
}

func TestClean(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	clean, err := modfilepath.Clean(ctx, object.NewString("/foo/../foo/bar//baz"))
	require.NoError(t, err)
	require.IsType(t, &object.String{}, clean)
	require.Equal(t, "/foo/bar/baz", clean.(*object.String).Value())
}

func TestDir(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	dir, err := modfilepath.Dir(ctx, object.NewString("/foo/bar/baz.txt"))
	require.NoError(t, err)
	require.IsType(t, &object.String{}, dir)
	require.Equal(t, "/foo/bar", dir.(*object.String).Value())
}

func TestExt(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	ext, err := modfilepath.Ext(ctx, object.NewString("bar/baz.txt"))
	require.NoError(t, err)
	require.IsType(t, &object.String{}, ext)
	require.Equal(t, ".txt", ext.(*object.String).Value())
}

func TestIsAbs(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	isAbsTrue, err := modfilepath.IsAbs(ctx, object.NewString("/foo/bar"))
	require.NoError(t, err)
	require.IsType(t, &object.Bool{}, isAbsTrue)
	require.True(t, isAbsTrue.(*object.Bool).Value())

	isAbsFalse, err := modfilepath.IsAbs(ctx, object.NewString("foo/bar"))
	require.NoError(t, err)
	require.IsType(t, &object.Bool{}, isAbsFalse)
	require.False(t, isAbsFalse.(*object.Bool).Value())
}

func TestJoin(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	join, err := modfilepath.Join(ctx, object.NewString("foo"), object.NewString("bar"), object.NewString("baz.txt"))
	require.NoError(t, err)
	require.IsType(t, &object.String{}, join)
	require.Equal(t, "foo/bar/baz.txt", join.(*object.String).Value())
}

func TestMatch(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	result, err := modfilepath.Match(ctx, object.NewString("*.txt"), object.NewString("file.txt"))
	require.NoError(t, err)
	require.IsType(t, &object.Bool{}, result)
	require.True(t, result.(*object.Bool).Value())

	result, err = modfilepath.Match(ctx, object.NewString("*.txt"), object.NewString("file.jpg"))
	require.NoError(t, err)
	require.IsType(t, &object.Bool{}, result)
	require.False(t, result.(*object.Bool).Value())
}

func TestSplit(t *testing.T) {
	m := &testutils.MockOS{}
	ctx := ren.WithOS(context.Background(), m)
	split, err := modfilepath.Split(ctx, object.NewString("/foo/bar/baz.txt"))
	require.NoError(t, err)
	require.IsType(t, &object.List{}, split)
	l := split.(*object.List)
	items := l.Value()
	require.Len(t, items, 2)
	require.Equal(t, "/foo/bar", items[0].(*object.String).Value())
	require.Equal(t, "baz.txt", items[1].(*object.String).Value())
}
