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

	"github.com/foohq/ren/builtins"
	"github.com/foohq/ren/importer"
	renos "github.com/foohq/ren/internal/os"
	"github.com/foohq/ren/modules"
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
	conf, err := buildConfig(opt...)
	if err != nil {
		return err
	}

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

	imp := importer.NewImporter(importer.ImporterOptions{
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

func buildConfig(opt ...Option) (*risor.Config, error) {
	var opts Options
	for _, o := range opt {
		o(&opts)
	}
	return opts.toConfig(), nil
}

func readEntrypoint(zr *zip.Reader) ([]byte, error) {
	f, err := zr.Open("entrypoint.json")
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
	opts []renos.Option
}

func (o *Options) toConfig() *risor.Config {
	var opts = []risor.Option{
		risor.WithOS(renos.New(o.opts...)),
		risor.WithoutDefaultGlobals(),
		risor.WithGlobals(builtins.Globals()),
		risor.WithGlobals(modules.Globals()),
	}
	return risor.NewConfig(opts...)
}

type Option func(*Options)

func WithEnvVar(name, value string) Option {
	return func(o *Options) {
		o.opts = append(o.opts, renos.WithEnvVar(name, value))
	}
}

func WithArgs(args []string) Option {
	return func(o *Options) {
		o.opts = append(o.opts, renos.WithArgs(args))
	}
}

func WithStdin(file risoros.File) Option {
	return func(o *Options) {
		o.opts = append(o.opts, renos.WithStdin(file))
	}
}
func WithStdout(file risoros.File) Option {
	return func(o *Options) {
		o.opts = append(o.opts, renos.WithStdout(file))
	}
}

func WithWorkDir(dir string) Option {
	return func(o *Options) {
		o.opts = append(o.opts, renos.WithWorkDir(dir))
	}
}

func WithFilesystems(fss map[string]risoros.FS) Option {
	return func(o *Options) {
		o.opts = append(o.opts, renos.WithFilesystems(fss))
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
