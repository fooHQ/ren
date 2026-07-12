package packager

import (
	"archive/zip"
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/deepnoodle-ai/risor/v2/pkg/ast"
	"github.com/deepnoodle-ai/risor/v2/pkg/compiler"
	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/deepnoodle-ai/risor/v2/pkg/parser"
)

var (
	ErrMissingEntrypoint = errors.New("missing entrypoint")
)

var exts = []string{".risor", ".rsr"}

const fileExt = "zip"

func NewFilename(name string) string {
	if strings.HasSuffix(name, fileExt) {
		return name
	}
	return name + "." + fileExt
}

func Build(src, dst string, opt ...Option) error {
	var opts options
	for _, o := range opt {
		o(&opts)
	}

	err := isEntrypoint(src)
	if err != nil {
		return err
	}

	tmpDir, err := walkSourceDir(src, &opts)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	tmpZip, err := createTempZip(tmpDir)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(tmpZip)
	}()

	err = os.Rename(tmpZip, dst)
	if err != nil {
		return err
	}

	return nil
}

func walkSourceDir(src string, opts *options) (string, error) {
	tmpDir, err := os.MkdirTemp(".", "ren*")
	if err != nil {
		return "", err
	}

	srcPrefix := ""
	err = filepath.Walk(src, func(srcPth string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		srcPth = filepath.ToSlash(srcPth)
		if srcPrefix == "" {
			srcPrefix = filepath.ToSlash(filepath.Clean(srcPth))
		}

		if info.IsDir() {
			// Skip a directory...
			return nil
		}

		dstPth := filepath.Join(tmpDir, strings.TrimPrefix(srcPth, srcPrefix))

		err = os.MkdirAll(filepath.Dir(dstPth), 0755)
		if err != nil {
			return err
		}

		if isRisorScript(srcPth) {
			rel := strings.TrimPrefix(srcPth, srcPrefix)
			err = compileScript(context.Background(), srcPth, dstPth, opts.GlobalNames(), !isEntrypointFile(rel))
		} else {
			err = copyFile(srcPth, dstPth)
		}
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		_ = os.RemoveAll(tmpDir)
	}

	return tmpDir, err
}

func createTempZip(src string) (string, error) {
	f, err := os.CreateTemp(".", "ren*."+fileExt)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	zw := zip.NewWriter(f)
	defer func() {
		_ = zw.Close()
	}()

	err = zw.AddFS(os.DirFS(src))
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func isEntrypoint(dir string) error {
	for _, ext := range exts {
		info, err := os.Stat(filepath.Join(dir, "entrypoint"+ext))
		if err == nil && info.Mode().IsRegular() {
			return nil
		}
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return ErrMissingEntrypoint
}

func isRisorScript(filename string) bool {
	for _, ext := range exts {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

func compileScript(ctx context.Context, src, dst string, globalNames []string, wrap bool) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	source := string(b)
	prog, err := parseSource(ctx, src, source)
	if err != nil {
		return err
	}

	code, err := compileProgram(src, source, prog, globalNames)
	if err != nil {
		return err
	}

	// A module is compiled as an immediately-invoked function returning a map of
	// its top-level names. Wrapping turns those names into locals, so exported
	// functions become self-contained closures independent of any VM's globals.
	// It is done on the AST, reusing the original statements, so that compiled
	// source locations still point at the module's real lines and columns. The
	// first compile above is what reveals the module's top-level names; the
	// wrapped program must then be compiled to produce the shipped bytecode.
	if wrap {
		names, err := topLevelNames(code, globalNames)
		if err != nil {
			return err
		}

		prog = wrapModule(prog, names)
		code, err = compileProgram(src, source, prog, globalNames)
		if err != nil {
			return err
		}
	}

	out, err := code.MarshalJSON()
	if err != nil {
		return err
	}

	return os.WriteFile(replaceScriptExt(dst), out, 0644)
}

func parseSource(ctx context.Context, filename, source string) (*ast.Program, error) {
	return parser.Parse(ctx, source, &parser.Config{
		Filename: filename,
		MaxDepth: 0,
	})
}

func compileProgram(filename, source string, prog *ast.Program, globalNames []string) (*compiler.Code, error) {
	comp, err := compiler.New(&compiler.Config{
		GlobalNames: globalNames,
		Filename:    filename,
		Source:      source,
	})
	if err != nil {
		return nil, err
	}

	return comp.CompileAST(prog)
}

// topLevelNames compiles a module unwrapped and returns the names it declares
// at the top level, i.e. the globals it introduces beyond the provided names.
func topLevelNames(code *compiler.Code, globalNames []string) ([]string, error) {
	provided := make(map[string]struct{}, len(globalNames))
	for _, name := range globalNames {
		provided[name] = struct{}{}
	}

	var names []string
	for _, name := range code.GlobalNames() {
		if _, ok := provided[name]; ok {
			continue
		}
		names = append(names, name)
	}
	return names, nil
}

// wrapModule rewrites a module program into an immediately-invoked function
// literal that returns a map of the given names. The original statements are
// reused verbatim so their source positions are preserved.
func wrapModule(prog *ast.Program, names []string) *ast.Program {
	items := make([]ast.MapItem, 0, len(names))
	for _, name := range names {
		items = append(items, ast.MapItem{
			Key:   &ast.String{Literal: name, Value: name},
			Value: &ast.Ident{Name: name},
		})
	}

	stmts := make([]ast.Node, 0, len(prog.Stmts)+1)
	stmts = append(stmts, prog.Stmts...)
	stmts = append(stmts, &ast.Return{Value: &ast.Map{Items: items}})

	call := &ast.Call{Fun: &ast.Func{Body: &ast.Block{Stmts: stmts}}}
	return &ast.Program{Stmts: []ast.Node{call}}
}

func isEntrypointFile(rel string) bool {
	rel = strings.TrimPrefix(filepath.ToSlash(rel), "/")
	for _, ext := range exts {
		if rel == "entrypoint"+ext {
			return true
		}
	}
	return false
}

func copyFile(src, dst string) error {
	fileSrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = fileSrc.Close()
	}()

	fileDst, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = fileDst.Close()
	}()

	_, err = io.Copy(fileDst, fileSrc)
	return err
}

func replaceScriptExt(filename string) string {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	return name + ".json"
}

type options struct {
	builtins []*object.Builtin
}

func (o *options) GlobalNames() []string {
	names := make([]string, 0, len(o.builtins))
	for _, builtin := range o.builtins {
		names = append(names, builtin.Name())
	}
	return names
}

type Option func(*options)

func WithBuiltin(builtin *object.Builtin) Option {
	return func(options *options) {
		if builtin == nil {
			return
		}
		options.builtins = append(options.builtins, builtin)
	}
}
