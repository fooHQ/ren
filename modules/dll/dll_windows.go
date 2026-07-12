//go:build windows

package dll

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"unsafe"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/deepnoodle-ai/risor/v2/pkg/op"
	"golang.org/x/sys/windows"
)

var _ object.Object = (*Handle)(nil)

// HANDLE is the Risor type name of a loaded-library handle.
const HANDLE = "handle"

// Handle is a Risor object wrapping an open dynamic-link library. Its lifetime
// is bound to the loading context: the library is released when the handle is
// closed or when the context is done.
type Handle struct {
	dll    *windows.DLL
	closed chan struct{}
}

// newHandle wraps an open DLL and starts a goroutine that releases it when the
// handle is closed or the context is done.
func newHandle(ctx context.Context, dll *windows.DLL) *Handle {
	h := &Handle{
		dll:    dll,
		closed: make(chan struct{}),
	}
	h.startCleanup(ctx)
	return h
}

// startCleanup releases the underlying library once the handle is closed or the
// context is cancelled.
func (h *Handle) startCleanup(ctx context.Context) {
	go func() {
		select {
		case <-h.closed:
		case <-ctx.Done():
		}
		_ = h.dll.Release()
	}()
}

// isClosed reports whether the handle has been closed.
func (h *Handle) isClosed() bool {
	select {
	case <-h.closed:
		return true
	default:
		return false
	}
}

// Type returns the Risor type name of the handle.
func (h *Handle) Type() object.Type {
	return HANDLE
}

// Inspect returns a human-readable representation of the handle.
func (h *Handle) Inspect() string {
	return fmt.Sprintf("handle(path=%s)", h.dll.Name)
}

// String returns a string representation of the handle.
func (h *Handle) String() string {
	return h.Inspect()
}

// IsTruthy reports whether the handle is truthy; it is always true.
func (h *Handle) IsTruthy() bool {
	return true
}

// Interface returns the underlying *windows.DLL.
func (h *Handle) Interface() any {
	return h.dll
}

// Equals reports whether other is the same handle instance.
func (h *Handle) Equals(other object.Object) bool {
	return h == other
}

// RunOperation always returns an error; handles support no binary operations.
func (h *Handle) RunOperation(opType op.BinaryOpType, right object.Object) (object.Object, error) {
	return nil, object.TypeErrorf("unsupported operation for handle: %v", opType)
}

// Attrs returns the attribute specifications for the handle's methods.
func (h *Handle) Attrs() []object.AttrSpec {
	return handleMethods.Specs()
}

// GetAttr returns the named method of the handle.
func (h *Handle) GetAttr(name string) (object.Object, bool) {
	return handleMethods.GetAttr(h, name)
}

// SetAttr always returns an error; handle attributes are read-only.
func (h *Handle) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("handle has no attribute %q", name)
}

// handleMethods holds the methods exposed on handle objects (lookup, close).
var handleMethods = object.NewMethodRegistry[*Handle](HANDLE)

func init() {
	handleMethods.Define("lookup").
		Doc("Look up a procedure in the library by name and return a callable proc.").
		Arg("name").
		Returns(PROC).
		Impl(func(h *Handle, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, object.NewArgsError("handle.lookup", 1, len(args))
			}
			if h.isClosed() {
				return nil, fmt.Errorf("handle.lookup: handle is closed")
			}
			name, err := object.AsString(args[0])
			if err != nil {
				return nil, err
			}
			p, procErr := h.dll.FindProc(name)
			if procErr != nil {
				return nil, object.NewError(procErr)
			}
			return newProc(h, p), nil
		})
	handleMethods.Define("close").
		Doc("Free the library handle.").
		Impl(func(h *Handle, ctx context.Context, args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, object.NewArgsError("handle.close", 0, len(args))
			}
			// close may be called more than once; closing the channel twice
			// would panic. Only the VM goroutine closes it, so this
			// check-then-close needs no further synchronization.
			if h.isClosed() {
				return object.Nil, nil
			}
			close(h.closed)
			return object.Nil, nil
		})
}

var (
	_ object.Object   = (*Proc)(nil)
	_ object.Callable = (*Proc)(nil)
)

// PROC is the Risor type name of a procedure looked up from a library.
const PROC = "proc"

// Proc is a callable Risor object representing a procedure exported by a loaded
// library. It retains its owning handle so a call can refuse to run against a
// library that has been closed.
type Proc struct {
	handle *Handle
	proc   *windows.Proc
}

// newProc wraps a resolved procedure together with the handle it belongs to.
func newProc(h *Handle, p *windows.Proc) *Proc {
	return &Proc{handle: h, proc: p}
}

// Type returns the Risor type name of the proc.
func (p *Proc) Type() object.Type {
	return PROC
}

// Inspect returns a human-readable representation of the proc.
func (p *Proc) Inspect() string {
	return fmt.Sprintf("proc(name=%s)", p.proc.Name)
}

// String returns a string representation of the proc.
func (p *Proc) String() string {
	return p.Inspect()
}

// IsTruthy reports whether the proc is truthy; it is always true.
func (p *Proc) IsTruthy() bool {
	return true
}

// Interface returns the underlying *windows.Proc.
func (p *Proc) Interface() any {
	return p.proc
}

// Equals reports whether other is the same proc instance.
func (p *Proc) Equals(other object.Object) bool {
	return p == other
}

// RunOperation always returns an error; procs support no binary operations.
func (p *Proc) RunOperation(opType op.BinaryOpType, right object.Object) (object.Object, error) {
	return nil, object.TypeErrorf("unsupported operation for proc: %v", opType)
}

// Attrs returns nil; a proc exposes no attributes and is instead called directly.
func (p *Proc) Attrs() []object.AttrSpec {
	return nil
}

// GetAttr always reports that no attribute exists; a proc is called directly.
func (p *Proc) GetAttr(name string) (object.Object, bool) {
	return nil, false
}

// SetAttr always returns an error; proc attributes are read-only.
func (p *Proc) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("proc has no attribute %q", name)
}

// toUintptr converts a Risor object to a uintptr suitable for a Windows proc
// call. pin receives any Go object that must be kept alive until after the
// call returns (use runtime.KeepAlive on the returned slice).
func toUintptr(obj object.Object, pin *[]any) (uintptr, error) {
	switch v := obj.(type) {
	case *object.Int:
		return uintptr(v.Value()), nil
	case *object.Bool:
		if v.Value() {
			return 1, nil
		}
		return 0, nil
	case *object.NilType:
		return 0, nil
	case *object.String:
		ptr, err := windows.UTF16PtrFromString(v.Value())
		if err != nil {
			return 0, err
		}
		*pin = append(*pin, ptr)
		return uintptr(unsafe.Pointer(ptr)), nil
	case *object.Bytes:
		b := v.Value()
		if len(b) == 0 {
			return 0, nil
		}
		*pin = append(*pin, b)
		return uintptr(unsafe.Pointer(&b[0])), nil
	default:
		return 0, fmt.Errorf("proc.call: expected int, bool, nil, string, or bytes, got %s", obj.Type())
	}
}

var _ object.Object = (*CallResult)(nil)

// CALL_RESULT is the Risor type name of a procedure call's result.
const CALL_RESULT = "call_result"

// CallResult is a Risor object holding the outcome of a procedure call: the
// return value and the error code (errno) set by the call.
type CallResult struct {
	value uintptr
	errno windows.Errno
}

// newCallResult builds a call_result from a procedure's return value and errno.
func newCallResult(value uintptr, errno windows.Errno) *CallResult {
	return &CallResult{value: value, errno: errno}
}

// Type returns the Risor type name of the call result.
func (r *CallResult) Type() object.Type {
	return CALL_RESULT
}

// Inspect returns a human-readable representation of the call result.
func (r *CallResult) Inspect() string {
	return fmt.Sprintf("call_result(value=%d, errno=%d)", r.value, r.errno)
}

// String returns a string representation of the call result.
func (r *CallResult) String() string {
	return r.Inspect()
}

// IsTruthy reports whether the call result is truthy; it is always true.
func (r *CallResult) IsTruthy() bool {
	return true
}

// Interface returns the call result itself.
func (r *CallResult) Interface() any {
	return r
}

// Equals reports whether other is the same call_result instance.
func (r *CallResult) Equals(other object.Object) bool {
	return r == other
}

// RunOperation always returns an error; call results support no binary operations.
func (r *CallResult) RunOperation(opType op.BinaryOpType, right object.Object) (object.Object, error) {
	return nil, object.TypeErrorf("unsupported operation for call_result: %v", opType)
}

// Attrs returns the attribute specifications for the call result's fields.
func (r *CallResult) Attrs() []object.AttrSpec {
	return callResultMethods.Specs()
}

// GetAttr returns the named field (value or errno) of the call result.
func (r *CallResult) GetAttr(name string) (object.Object, bool) {
	return callResultMethods.GetAttr(r, name)
}

// SetAttr always returns an error; call_result attributes are read-only.
func (r *CallResult) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("call_result has no attribute %q", name)
}

// callResultMethods holds the attributes exposed on call_result objects.
var callResultMethods = object.NewMethodRegistry[*CallResult](CALL_RESULT)

func init() {
	callResultMethods.Define("value").
		Doc("The return value of the procedure call.").
		Returns("int").
		Getter(func(r *CallResult) object.Object {
			return object.NewInt(int64(r.value))
		})
	callResultMethods.Define("errno").
		Doc("The error code set by the procedure. Zero means no error occurred.").
		Returns("int").
		Getter(func(r *CallResult) object.Object {
			return object.NewInt(int64(r.errno))
		})
}

// Call invokes the procedure with the given arguments and returns a
// call_result. Implementing object.Callable makes a proc directly callable
// from Risor, e.g. proc(1, 2, "text").
func (p *Proc) Call(ctx context.Context, args ...object.Object) (object.Object, error) {
	// Refuse to call into a library that has been closed; its code may be
	// unmapped, which would fault the process rather than panic recoverably.
	if p.handle.isClosed() {
		return nil, fmt.Errorf("proc.call: library has been closed")
	}
	uargs := make([]uintptr, len(args))
	var pin []any
	for i, arg := range args {
		u, err := toUintptr(arg, &pin)
		if err != nil {
			return nil, err
		}
		uargs[i] = u
	}
	r1, _, lastErr := p.proc.Call(uargs...)
	runtime.KeepAlive(pin)
	errno, _ := errors.AsType[windows.Errno](lastErr)
	return newCallResult(r1, errno), nil
}

// Load opens the dynamic-link library at the given path and returns a handle
// to it. It takes a single string argument.
func Load(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("dll.load", 1, len(args))
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	lib, loadErr := windows.LoadLibraryEx(path, 0, 0)
	if loadErr != nil {
		return nil, object.NewError(loadErr)
	}
	dll := &windows.DLL{Name: path, Handle: lib}
	return newHandle(ctx, dll), nil
}

// Module returns the "dll" module with its load function registered.
func Module() *object.Module {
	return object.NewBuiltinsModule("dll", map[string]object.Object{
		"load": object.NewBuiltin("load", Load),
	})
}
