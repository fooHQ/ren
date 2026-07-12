package ren

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/deepnoodle-ai/risor/v2/pkg/compiler"
	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/deepnoodle-ai/risor/v2/pkg/vm"

	"github.com/foohq/urlpath"
)

// Import URL schemes.
const (
	// builtinScheme imports a built-in module registered with the runtime,
	// e.g. import("builtin://os").
	builtinScheme = "builtin"
	// packageScheme imports a module from the package, e.g.
	// import("package://path/to/mod"). It is the explicit form of a bare
	// package path.
	packageScheme = "package"
)

// moduleExt is the file extension of compiled modules inside a package.
const moduleExt = ".json"

// Importer resolves and loads modules for the import builtin. It is placed on
// the context by Run and retrieved by the import builtin at call time.
//
// Imports are resolved by URL scheme:
//
//   - import("path/to/mod")           loads a compiled module from the package,
//     relative to the package root.
//   - import("package://path/to/mod") is the explicit form of the above.
//   - import("builtin://name")        returns a built-in module registered with
//     the runtime.
//
// A package module is executed at most once; its exports are cached and reused
// for subsequent imports. The packager compiles each module as a self-contained
// function returning a map of its top-level names, so running the module yields
// that map directly and its exported functions are usable in the importing VM.
type Importer struct {
	pkg      *zip.Reader
	builtins map[string]*object.Module
	env      map[string]any

	mu      sync.Mutex
	cache   map[string]object.Object
	loading map[string]struct{}
}

func newImporter(pkg *zip.Reader, builtins map[string]*object.Module, env map[string]any) *Importer {
	return &Importer{
		pkg:      pkg,
		builtins: builtins,
		env:      env,
		cache:    make(map[string]object.Object),
		loading:  make(map[string]struct{}),
	}
}

type importerContextKey struct{}

// WithImporter returns a new context carrying the given Importer.
func WithImporter(ctx context.Context, imp *Importer) context.Context {
	return context.WithValue(ctx, importerContextKey{}, imp)
}

// GetImporter returns the Importer from the context, or nil if none is set.
func GetImporter(ctx context.Context) *Importer {
	imp, _ := ctx.Value(importerContextKey{}).(*Importer)
	return imp
}

// Import resolves name to a module. See the Importer documentation for the
// supported import forms.
func (imp *Importer) Import(ctx context.Context, name string) (object.Object, error) {
	scheme, err := urlpath.Scheme(name)
	if err != nil {
		return nil, fmt.Errorf("cannot import %q: %w", name, err)
	}

	switch scheme {
	case "", packageScheme:
		return imp.importPackage(ctx, name)
	case builtinScheme:
		return imp.importBuiltin(name)
	default:
		return nil, fmt.Errorf("cannot import %q: unsupported scheme %q", name, scheme)
	}
}

func (imp *Importer) importBuiltin(name string) (object.Object, error) {
	modName := strings.TrimPrefix(name, builtinScheme+"://")
	mod, ok := imp.builtins[modName]
	if !ok {
		return nil, fmt.Errorf("cannot import %q: no such built-in module", name)
	}
	return mod, nil
}

func (imp *Importer) importPackage(ctx context.Context, name string) (object.Object, error) {
	pkgName := strings.TrimPrefix(name, packageScheme+"://")
	pth, err := packagePath(pkgName)
	if err != nil {
		return nil, fmt.Errorf("cannot import %q: %w", name, err)
	}

	imp.mu.Lock()
	if mod, ok := imp.cache[pth]; ok {
		imp.mu.Unlock()
		return mod, nil
	}
	if _, ok := imp.loading[pth]; ok {
		imp.mu.Unlock()
		return nil, fmt.Errorf("cannot import %q: import cycle detected", name)
	}
	imp.loading[pth] = struct{}{}
	imp.mu.Unlock()

	defer func() {
		imp.mu.Lock()
		delete(imp.loading, pth)
		imp.mu.Unlock()
	}()

	mod, err := imp.loadPackage(ctx, name, pth)
	if err != nil {
		return nil, err
	}

	imp.mu.Lock()
	imp.cache[pth] = mod
	imp.mu.Unlock()

	return mod, nil
}

// loadPackage reads and runs a compiled module from the package, returning it
// as a module object. The packager compiles a module as a function that returns
// a map of its top-level names, so running it yields that map; the exported
// functions are self-contained closures usable directly by the importing
// script. The map is wrapped in a module so that attribute access is not
// shadowed by the built-in methods of a map (get, keys, values, ...).
func (imp *Importer) loadPackage(ctx context.Context, name, pth string) (object.Object, error) {
	b, err := imp.readFile(pth)
	if err != nil {
		return nil, fmt.Errorf("cannot import %q: %w", name, err)
	}

	code, err := compiler.UnmarshalCode(b)
	if err != nil {
		return nil, fmt.Errorf("cannot import %q: %w", name, err)
	}

	result, err := vm.Run(ctx, code.ToBytecode(), vm.WithGlobals(imp.env))
	if err != nil {
		return nil, fmt.Errorf("cannot import %q: %w", name, err)
	}

	exports, ok := result.(*object.Map)
	if !ok {
		return nil, fmt.Errorf("cannot import %q: module did not produce exports", name)
	}

	return object.NewBuiltinsModule(name, exports.Value()), nil
}

func (imp *Importer) readFile(pth string) ([]byte, error) {
	f, err := imp.pkg.Open(pth)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	return io.ReadAll(f)
}

// packagePath converts an import name into a slash path within the package,
// rooted at the package root. It rejects paths that escape the root.
func packagePath(name string) (string, error) {
	pth, err := urlpath.Path(name)
	if err != nil {
		return "", err
	}
	pth = strings.TrimPrefix(pth, "/")
	if pth == "" || pth == "." || strings.HasPrefix(pth, "./") || pth == ".." || strings.HasPrefix(pth, "../") {
		return "", fmt.Errorf("invalid import path")
	}

	if !strings.HasSuffix(pth, moduleExt) {
		pth += moduleExt
	}
	return pth, nil
}
