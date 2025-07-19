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

	tmpDir, err := copyToTempDir(src)
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

func copyToTempDir(src string) (string, error) {
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
			b, err := os.ReadFile(srcPth)
			if err != nil {
				return err
			}

			prog, err := parser.Parse(context.Background(), string(b))
			if err != nil {
				return err
			}

			code, err := compiler.Compile(
				prog,
				compiler.WithFilename(dstPth),
				compiler.WithGlobalNames(modules.GlobalNames()),
			)
			if err != nil {
				return err
			}

			b, err = code.MarshalJSON()
			if err != nil {
				return err
			}

			err = os.WriteFile(replaceScriptExt(dstPth), b, 0644)
			if err != nil {
				return err
			}
		} else {
			fileSrc, err := os.Open(srcPth)
			if err != nil {
				return err
			}
			defer fileSrc.Close()

			fileDst, err := os.Create(dstPth)
			if err != nil {
				return err
			}
			defer fileDst.Close()

			_, err = io.Copy(fileDst, fileSrc)
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

func replaceScriptExt(filename string) string {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	return name + ".json"
}
