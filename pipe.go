package ren

import (
	"errors"
	"io"
	"io/fs"
	"time"

	risoros "github.com/risor-io/risor/os"
)

var _ risoros.File = &Pipe{}

// Pipe implements Risor's os.File interface.
// The type is backed by Go's io.Pipe, therefore, it allows concurrent read/write.
type Pipe struct {
	r *io.PipeReader
	w *io.PipeWriter
}

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

func (f *Pipe) Stat() (risoros.FileInfo, error) {
	return risoros.NewFileInfo(risoros.GenericFileInfoOpts{
		Name:    "grr",
		Size:    0,
		Mode:    0,
		ModTime: time.Time{},
		IsDir:   false,
	}), nil
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
