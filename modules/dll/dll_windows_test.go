//go:build windows

package dll_test

import (
	"context"
	"testing"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/stretchr/testify/require"

	"github.com/foohq/ren/modules/dll"
)

func callMethod(t *testing.T, obj object.Object, name string, args ...object.Object) (object.Object, error) {
	t.Helper()
	attr, ok := obj.GetAttr(name)
	require.True(t, ok, "missing method %q", name)
	builtin, ok := attr.(*object.Builtin)
	require.True(t, ok, "%q is not callable", name)
	return builtin.Call(context.Background(), args...)
}

func callProc(t *testing.T, p object.Object, args ...object.Object) (object.Object, error) {
	t.Helper()
	fn, ok := p.(object.Callable)
	require.True(t, ok, "proc is not directly callable")
	return fn.Call(context.Background(), args...)
}

func getProp(t *testing.T, obj object.Object, name string) object.Object {
	t.Helper()
	attr, ok := obj.GetAttr(name)
	require.True(t, ok, "missing property %q", name)
	return attr
}

func resultValue(t *testing.T, result object.Object) int64 {
	t.Helper()
	v, ok := getProp(t, result, "value").(*object.Int)
	require.True(t, ok)
	return v.Value()
}

func mustLoad(t *testing.T, path string) object.Object {
	t.Helper()
	h, err := dll.Load(context.Background(), object.NewString(path))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = callMethod(t, h, "close")
	})
	return h
}

func mustLookup(t *testing.T, h object.Object, name string) object.Object {
	t.Helper()
	p, err := callMethod(t, h, "lookup", object.NewString(name))
	require.NoError(t, err)
	return p
}

func TestModule(t *testing.T) {
	mod := dll.Module()
	require.Equal(t, "dll", mod.Name().Value())
	load, ok := mod.GetAttr("load")
	require.True(t, ok)
	require.IsType(t, &object.Builtin{}, load)
}

func TestLoadSuccess(t *testing.T) {
	h := mustLoad(t, "kernel32.dll")
	require.IsType(t, &dll.Handle{}, h)
}

func TestLoadErrors(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		args []object.Object
	}{
		{"nonexistent library", []object.Object{object.NewString("no_such_library_xyz.dll")}},
		{"no arguments", nil},
		{"non-string argument", []object.Object{object.NewInt(1)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dll.Load(ctx, tt.args...)
			require.Error(t, err)
		})
	}
}

func TestLookup(t *testing.T) {
	h := mustLoad(t, "kernel32.dll")

	t.Run("existing proc", func(t *testing.T) {
		p, err := callMethod(t, h, "lookup", object.NewString("GetCurrentProcessId"))
		require.NoError(t, err)
		require.IsType(t, &dll.Proc{}, p)
	})

	t.Run("missing proc", func(t *testing.T) {
		_, err := callMethod(t, h, "lookup", object.NewString("NoSuchProcedureXYZ"))
		require.Error(t, err)
	})
}

func TestCall(t *testing.T) {
	h := mustLoad(t, "kernel32.dll")

	t.Run("no arguments", func(t *testing.T) {
		p := mustLookup(t, h, "GetCurrentProcessId")
		res, err := callProc(t, p)
		require.NoError(t, err)
		require.IsType(t, &dll.CallResult{}, res)
		require.NotZero(t, resultValue(t, res)) // a real PID
	})

	t.Run("integer arguments", func(t *testing.T) {
		p := mustLookup(t, h, "MulDiv") // MulDiv(a, b, c) = a*b/c
		res, err := callProc(t, p, object.NewInt(10), object.NewInt(3), object.NewInt(2))
		require.NoError(t, err)
		require.Equal(t, int64(15), resultValue(t, res))
	})

	t.Run("string argument", func(t *testing.T) {
		p := mustLookup(t, h, "lstrlenW") // wide-string length
		res, err := callProc(t, p, object.NewString("hello"))
		require.NoError(t, err)
		require.Equal(t, int64(5), resultValue(t, res))
	})

	t.Run("nil argument", func(t *testing.T) {
		p := mustLookup(t, h, "lstrlenW") // lstrlenW(NULL) == 0
		res, err := callProc(t, p, object.Nil)
		require.NoError(t, err)
		require.Equal(t, int64(0), resultValue(t, res))
	})

	t.Run("bytes argument is filled in place", func(t *testing.T) {
		p := mustLookup(t, h, "QueryPerformanceCounter") // BOOL, writes an 8-byte counter
		buf := object.NewBytes(make([]byte, 8))
		res, err := callProc(t, p, buf)
		require.NoError(t, err)
		require.NotZero(t, resultValue(t, res)) // nonzero BOOL == success
		require.NotEqual(t, make([]byte, 8), buf.Value(), "buffer should have been written")
	})

	t.Run("unsupported argument type", func(t *testing.T) {
		p := mustLookup(t, h, "GetCurrentProcessId")
		_, err := callProc(t, p, object.NewList(nil))
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected int, bool, nil, string, or bytes")
	})
}

func TestCallResultErrno(t *testing.T) {
	h := mustLoad(t, "kernel32.dll")
	// FreeLibrary(NULL) fails and sets last error to ERROR_INVALID_HANDLE.
	p := mustLookup(t, h, "FreeLibrary")
	res, err := callProc(t, p, object.NewInt(0))
	require.NoError(t, err)
	require.Equal(t, int64(0), resultValue(t, res))

	errno, ok := getProp(t, res, "errno").(*object.Int)
	require.True(t, ok)
	require.NotZero(t, errno.Value())
}

func TestCloseIsIdempotent(t *testing.T) {
	h := mustLoad(t, "kernel32.dll")
	_, err := callMethod(t, h, "close")
	require.NoError(t, err)
	// A second close must be safe (guarded by sync.Once).
	_, err = callMethod(t, h, "close")
	require.NoError(t, err)
}

func TestLookupAfterCloseErrors(t *testing.T) {
	h := mustLoad(t, "kernel32.dll")
	_, err := callMethod(t, h, "close")
	require.NoError(t, err)

	// Must return an error, not panic or use the freed library.
	_, err = callMethod(t, h, "lookup", object.NewString("GetCurrentProcessId"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "closed")
}

func TestCallAfterCloseErrors(t *testing.T) {
	h := mustLoad(t, "kernel32.dll")
	p := mustLookup(t, h, "GetCurrentProcessId") // resolve before closing

	_, err := callMethod(t, h, "close")
	require.NoError(t, err)

	// Calling into the closed library must error rather than fault the process.
	_, err = callProc(t, p)
	require.Error(t, err)
	require.Contains(t, err.Error(), "closed")
}

func TestContextCancellationIsSafe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	h, err := dll.Load(ctx, object.NewString("kernel32.dll"))
	require.NoError(t, err)

	// Cancelling frees the handle from the cleanup goroutine; an explicit
	// close afterwards must not double-free or panic.
	cancel()
	_, err = callMethod(t, h, "close")
	require.NoError(t, err)
}

func TestHandleAttrErrors(t *testing.T) {
	h := mustLoad(t, "kernel32.dll")

	_, ok := h.GetAttr("does_not_exist")
	require.False(t, ok)

	err := h.SetAttr("anything", object.Nil)
	require.Error(t, err)
}
