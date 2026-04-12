// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

package ren

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"

	"github.com/foohq/urlpath"
)

var (
	// ErrFSNotFound is returned when a filesystem is not found for a given scheme.
	ErrFSNotFound = errors.New("filesystem not found")
	// ErrCrossingFSBoundaries is returned when an operation crosses filesystem boundaries.
	ErrCrossingFSBoundaries = errors.New("crossing filesystem boundaries")
)

// FS is an interface for a filesystem.
type FS interface {
	Mkdir(name string, perm FileMode) error
	MkdirAll(path string, perm FileMode) error
	MkdirTemp(dir, pattern string) (string, error)
	OpenFile(name string, flag int, perm FileMode) (File, error)
	ReadFile(name string) ([]byte, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldpath, newpath string) error
	Stat(name string) (FileInfo, error)
	Symlink(oldname, newname string) error
	WriteFile(name string, data []byte, perm FileMode) error
	ReadDir(name string) ([]DirEntry, error)
}

var _ FS = fsMiddleware{}

type fsMiddleware map[string]FS

func (f fsMiddleware) Mkdir(name string, perm FileMode) error {
	fs, err := f.lookupFS(name)
	if err != nil {
		return err
	}
	pth, err := urlpath.Path(name)
	if err != nil {
		return err
	}
	err = fs.Mkdir(pth, perm)
	if err != nil {
		return fmt.Errorf("mkdir %s: %w", name, err)
	}
	return nil
}

func (f fsMiddleware) MkdirAll(path string, perm FileMode) error {
	fs, err := f.lookupFS(path)
	if err != nil {
		return err
	}
	pth, err := urlpath.Path(path)
	if err != nil {
		return err
	}
	err = fs.MkdirAll(pth, perm)
	if err != nil {
		return fmt.Errorf("mkdir %s: %w", path, err)
	}
	return nil
}

func (f fsMiddleware) MkdirTemp(dir, pattern string) (string, error) {
	fs, err := f.lookupFS(dir)
	if err != nil {
		return "", err
	}
	pth, err := urlpath.Path(dir)
	if err != nil {
		return "", err
	}
	dir, err = fs.MkdirTemp(pth, pattern)
	if err != nil {
		return "", fmt.Errorf("mkdir %s: %w", dir, err)
	}
	return dir, nil
}

func (f fsMiddleware) OpenFile(name string, flag int, perm FileMode) (File, error) {
	fs, err := f.lookupFS(name)
	if err != nil {
		return nil, err
	}
	pth, err := urlpath.Path(name)
	if err != nil {
		return nil, err
	}
	file, err := fs.OpenFile(pth, flag, perm)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	}
	return file, nil
}

func (f fsMiddleware) ReadFile(name string) ([]byte, error) {
	fs, err := f.lookupFS(name)
	if err != nil {
		return nil, err
	}
	pth, err := urlpath.Path(name)
	if err != nil {
		return nil, err
	}
	b, err := fs.ReadFile(pth)
	if err != nil {
		return nil, fmt.Errorf("open %s: %v", name, err)
	}
	return b, nil
}

func (f fsMiddleware) Remove(name string) error {
	fs, err := f.lookupFS(name)
	if err != nil {
		return err
	}
	pth, err := urlpath.Path(name)
	if err != nil {
		return err
	}
	err = fs.Remove(pth)
	if err != nil {
		return fmt.Errorf("remove %s: %w", name, err)
	}
	return nil
}

func (f fsMiddleware) RemoveAll(path string) error {
	fs, err := f.lookupFS(path)
	if err != nil {
		return err
	}
	pth, err := urlpath.Path(path)
	if err != nil {
		return err
	}
	err = fs.RemoveAll(pth)
	if err != nil {
		return fmt.Errorf("remove %s: %w", path, err)
	}
	return nil
}

func (f fsMiddleware) Rename(oldPath, newPath string) error {
	oldFS, err := f.lookupFS(oldPath)
	if err != nil {
		return err
	}
	newFS, err := f.lookupFS(newPath)
	if err != nil {
		return err
	}
	if oldFS != newFS {
		return ErrCrossingFSBoundaries
	}
	oldPth, err := urlpath.Path(oldPath)
	if err != nil {
		return err
	}
	newPth, err := urlpath.Path(newPath)
	if err != nil {
		return err
	}
	err = oldFS.Rename(oldPth, newPth)
	if err != nil {
		return fmt.Errorf("rename %s %s: %w", oldPath, newPath, err)
	}
	return nil
}

func (f fsMiddleware) Stat(name string) (FileInfo, error) {
	fs, err := f.lookupFS(name)
	if err != nil {
		return nil, err
	}
	pth, err := urlpath.Path(name)
	if err != nil {
		return nil, err
	}
	info, err := fs.Stat(pth)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", name, err)
	}
	return info, nil
}

func (f fsMiddleware) Symlink(oldName, newName string) error {
	oldFS, err := f.lookupFS(oldName)
	if err != nil {
		return err
	}
	newFS, err := f.lookupFS(newName)
	if err != nil {
		return err
	}
	if oldFS != newFS {
		return ErrCrossingFSBoundaries
	}
	oldPth, err := urlpath.Path(oldName)
	if err != nil {
		return err
	}
	newPth, err := urlpath.Path(newName)
	if err != nil {
		return err
	}
	err = oldFS.Symlink(oldPth, newPth)
	if err != nil {
		return fmt.Errorf("symlink %s %s: %w", oldPth, newPth, err)
	}
	return nil
}

func (f fsMiddleware) WriteFile(name string, data []byte, perm FileMode) error {
	fs, err := f.lookupFS(name)
	if err != nil {
		return err
	}
	pth, err := urlpath.Path(name)
	if err != nil {
		return err
	}
	err = fs.WriteFile(pth, data, perm)
	if err != nil {
		return fmt.Errorf("open %s: %w", name, err)
	}
	return nil
}

func (f fsMiddleware) ReadDir(name string) ([]DirEntry, error) {
	fs, err := f.lookupFS(name)
	if err != nil {
		return nil, err
	}

	pth, err := urlpath.Path(name)
	if err != nil {
		return nil, err
	}

	results, err := fs.ReadDir(pth)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	}

	entries := make([]DirEntry, 0, len(results))
	entries = append(entries, results...)

	return entries, nil
}

func (f fsMiddleware) lookupFS(pth string) (FS, error) {
	scheme, err := urlpath.Scheme(pth)
	if err != nil {
		return nil, err
	}

	if scheme == "" {
		scheme = "file"
	}

	fs, ok := f[scheme]
	if !ok {
		return nil, ErrFSNotFound
	}
	return fs, nil
}

var _ FS = &localFS{}

type localFS struct{}

func (f *localFS) Mkdir(name string, perm FileMode) error {
	err := os.Mkdir(name, perm)
	if err != nil {
		return errors.Unwrap(err)
	}
	return nil
}

func (f *localFS) MkdirAll(path string, perm FileMode) error {
	err := os.MkdirAll(path, perm)
	if err != nil {
		return errors.Unwrap(err)
	}
	return nil
}

func (f *localFS) MkdirTemp(dir, pattern string) (string, error) {
	dir, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return "", errors.Unwrap(err)
	}
	return dir, nil
}

func (f *localFS) OpenFile(name string, flag int, perm FileMode) (File, error) {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, errors.Unwrap(err)
	}
	return file, nil
}

func (f *localFS) ReadFile(name string) ([]byte, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, errors.Unwrap(err)
	}
	return b, nil
}

func (f *localFS) Remove(name string) error {
	err := os.Remove(name)
	if err != nil {
		return errors.Unwrap(err)
	}
	return nil
}

func (f *localFS) RemoveAll(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return errors.Unwrap(err)
	}
	return nil
}

func (f *localFS) Rename(oldPath, newPath string) error {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return errors.Unwrap(err)
	}
	return nil
}

func (f *localFS) Stat(name string) (FileInfo, error) {
	info, err := os.Stat(name)
	if err != nil {
		return nil, errors.Unwrap(err)
	}
	return info, nil
}

func (f *localFS) Symlink(oldName, newName string) error {
	err := os.Symlink(oldName, newName)
	if err != nil {
		return errors.Unwrap(err)
	}
	return nil
}

func (f *localFS) WriteFile(name string, data []byte, perm FileMode) error {
	err := os.WriteFile(name, data, perm)
	if err != nil {
		return errors.Unwrap(err)
	}
	return nil
}

func (f *localFS) ReadDir(name string) ([]DirEntry, error) {
	results, err := os.ReadDir(name)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	entries := make([]DirEntry, 0, len(results))
	entries = append(entries, results...)

	return entries, nil
}

type (
	// FileMode represents a file's mode and permission bits.
	FileMode = fs.FileMode
	// FileInfo describes a file and is returned by Stat.
	FileInfo = fs.FileInfo
	// DirEntry is an entry read from a directory.
	DirEntry = fs.DirEntry
)

// File represents an open file.
type File interface {
	fs.File
	io.Writer
}

var (
	_ File = (*Pipe)(nil)
)

// Pipe implements ren's os.File interface and allows concurrent reads and writes.
type Pipe struct {
	r *io.PipeReader
	w *io.PipeWriter
}

// NewPipe returns a new Pipe.
func NewPipe() *Pipe {
	r, w := io.Pipe()
	return &Pipe{
		r: r,
		w: w,
	}
}

func (f *Pipe) Write(p []byte) (int, error) {
	n, err := f.w.Write(p)
	if errors.Is(err, io.ErrClosedPipe) {
		return n, fs.ErrClosed
	}
	return n, err
}

func (f *Pipe) Stat() (FileInfo, error) {
	return &pipeInfo{
		name:    "grr",
		size:    0,
		mode:    0,
		modTime: time.Time{},
		isDir:   false,
	}, nil
}

func (f *Pipe) Read(p []byte) (int, error) {
	n, err := f.r.Read(p)
	if errors.Is(err, io.ErrClosedPipe) {
		return n, fs.ErrClosed
	}
	return n, err
}

func (f *Pipe) Close() error {
	wErr := f.w.Close()
	rErr := f.r.Close()

	var err error
	if wErr != nil {
		err = wErr
	}
	if rErr != nil {
		err = rErr
	}
	return err
}

var _ FileInfo = (*pipeInfo)(nil)

type pipeInfo struct {
	name    string
	size    int64
	mode    FileMode
	modTime time.Time
	isDir   bool
}

func (fi *pipeInfo) Name() string {
	return fi.name
}

func (fi *pipeInfo) Size() int64 {
	return fi.size
}

func (fi *pipeInfo) Mode() FileMode {
	return fi.mode
}

func (fi *pipeInfo) ModTime() time.Time {
	return fi.modTime
}

func (fi *pipeInfo) IsDir() bool {
	return fi.isDir
}

func (fi *pipeInfo) Sys() any {
	return nil
}
