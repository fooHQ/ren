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

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/parser"

	"github.com/foohq/ren/builtins"
	"github.com/foohq/ren/modules"
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

func Build(src, dst string) error {
	err := isEntrypoint(src)
	if err != nil {
		return err
	}

	tmpDir, err := walkSourceDir(src)
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

func walkSourceDir(src string) (string, error) {
	tmpDir, err := os.MkdirTemp(".", "ren*")
	if err != nil {
		return "", err
	}

	srcPrefix := ""
	err = filepath.Walk(src, func(srcPth string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if srcPrefix == "" {
			srcPrefix = srcPth
		}

		if info.IsDir() {
			// Skip a directory...
			return nil
		}

		dstPth := filepath.Join(tmpDir, strings.TrimPrefix(filepath.ToSlash(srcPth), filepath.ToSlash(srcPrefix)))

		err = os.MkdirAll(filepath.Dir(dstPth), 0755)
		if err != nil {
			return err
		}

		if isRisorScript(srcPth) {
			err = compileScript(context.Background(), srcPth, dstPth)
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

func compileScript(ctx context.Context, src, dst string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	prog, err := parser.Parse(ctx, string(b))
	if err != nil {
		return err
	}

	var globalNames []string
	globalNames = append(globalNames, modules.Modules()...)
	globalNames = append(globalNames, builtins.Builtins()...)

	code, err := compiler.Compile(
		prog,
		compiler.WithFilename(src),
		compiler.WithGlobalNames(globalNames),
	)
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
	defer fileSrc.Close()

	fileDst, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fileDst.Close()

	_, err = io.Copy(fileDst, fileSrc)
	return err
}

func replaceScriptExt(filename string) string {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	return name + ".json"
}
