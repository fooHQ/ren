package ren

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	"github.com/risor-io/risor"
	"github.com/risor-io/risor/compiler"
	risoros "github.com/risor-io/risor/os"
	"github.com/risor-io/risor/parser"
	"github.com/risor-io/risor/vm"

	"github.com/foohq/ren/importer"
)

func RunBytes(ctx context.Context, b []byte, opts ...Option) error {
	reader := bytes.NewReader(b)
	return Run(ctx, reader, reader.Size(), opts...)
}

func RunFile(ctx context.Context, filename string, opts ...Option) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	inf, err := f.Stat()
	if err != nil {
		return err
	}

	return Run(ctx, f, inf.Size(), opts...)
}

func Run(ctx context.Context, reader io.ReaderAt, size int64, opt ...Option) error {
	var opts Options
	for _, o := range opt {
		o(&opts)
	}
	conf := opts.toConfig()

	zr, err := zip.NewReader(reader, size)
	if err != nil {
		return err
	}

	b, err := readEntrypoint(zr)
	if err != nil {
		return err
	}

	code, err := compiler.UnmarshalCode(b)
	if err != nil {
		return err
	}

	imp := importer.NewFSImporter(importer.FSImporterOptions{
		GlobalNames: conf.GlobalNames(),
		SourceFS:    zr,
	})

	vmOpts := conf.VMOpts()
	vmOpts = append(vmOpts, vm.WithImporter(imp))
	_, err = vm.Run(ctx, code, vmOpts...)
	if err != nil {
		return &Error{err}
	}

	return nil
}

func readEntrypoint(zr *zip.Reader) ([]byte, error) {
	f, err := zr.Open("main.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type Options struct {
	os      risoros.OS
	globals map[string]any
}

func (o *Options) Validate() error {
	if o.os == nil {
		return errors.New("ren: OS not specified")
	}

	return nil
}

func (o *Options) toConfig() *risor.Config {
	var opts = []risor.Option{
		risor.WithoutDefaultGlobals(),
	}

	if o.os != nil {
		opts = append(opts, risor.WithOS(o.os))
	}

	if o.globals != nil {
		opts = append(opts, risor.WithGlobals(o.globals))
	}

	return risor.NewConfig(opts...)
}

type Option func(*Options)

func WithOS(os risoros.OS) Option {
	return func(o *Options) {
		o.os = os
	}
}

func WithGlobals(globals map[string]any) Option {
	return func(o *Options) {
		if o.globals == nil {
			o.globals = make(map[string]any)
		}
		for k, v := range globals {
			o.globals[k] = v
		}
	}
}

type Error struct {
	err error
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Error() string {
	var parserErr parser.ParserError
	if errors.As(e.err, &parserErr) {
		return parserErr.FriendlyErrorMessage()
	}
	return e.err.Error()
}
