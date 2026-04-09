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
			err = compileScript(context.Background(), srcPth, dstPth, opts.GlobalNames())
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

func compileScript(ctx context.Context, src, dst string, globalNames []string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	prog, err := parser.Parse(ctx, string(b), &parser.Config{
		Filename: src,
		MaxDepth: 0,
	})
	if err != nil {
		return err
	}

	comp, err := compiler.New(&compiler.Config{
		GlobalNames: globalNames,
		Filename:    src,
		Source:      string(b),
	})
	if err != nil {
		return err
	}

	code, err := comp.CompileAST(prog)
	if err != nil {
		return err
	}

	b, err = code.MarshalJSON()
	if err != nil {
		return err
	}

	err = os.WriteFile(replaceScriptExt(dst), b, 0644)
	if err != nil {
		return err
	}
	return nil
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
	modules  []*object.Module
}

func (o *options) GlobalNames() []string {
	names := make([]string, 0, len(o.modules))
	for _, builtin := range o.builtins {
		names = append(names, builtin.Name())
	}
	for _, module := range o.modules {
		names = append(names, module.Name().Value())
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

func WithModule(module *object.Module) Option {
	return func(o *options) {
		if module == nil {
			return
		}
		o.modules = append(o.modules, module)
	}
}
