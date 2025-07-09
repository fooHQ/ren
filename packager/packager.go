package packager

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrIsEmpty     = errors.New("directory is empty")
	ErrInvalidMain = errors.New("main file is not a regular file")
	ErrMissingMain = errors.New("main file is missing")
)

const fileExt = "zip"

func NewFilename(name string) string {
	if strings.HasSuffix(name, fileExt) {
		return name
	}
	return name + "." + fileExt
}

func Build(src, dst string) error {
	err := isEmpty(src)
	if err != nil {
		return err
	}

	err = isMain(src)
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
		if srcPrefix == "" {
			srcPrefix = srcPth
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
		if err != nil {
			return err
		}

		return nil
	})

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

func isEmpty(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if errors.Is(err, io.EOF) {
		return ErrIsEmpty
	}

	return err
}

func isMain(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	files, err := f.Readdirnames(-1)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return err
		}
		return ErrMissingMain
	}

	for _, name := range files {
		if name != "main.risor" && name != "main.rsr" {
			continue
		}

		info, err := os.Stat(filepath.Join(dir, name))
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return ErrInvalidMain
		}

		return nil
	}

	return ErrMissingMain
}
