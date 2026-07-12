// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

// Package ren executes compiled Ren script packages.
//
// A package is a zip archive containing an entrypoint plus any modules and
// data files it imports. Run and its helpers load such a package, wire up the
// execution environment — built-in functions, importable modules, a
// scheme-based virtual filesystem, and an OS abstraction — and run the
// entrypoint on the embedded Risor virtual machine. Behaviour is configured
// through Option values passed to the Run functions.
package ren

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"maps"
	"os"

	"github.com/deepnoodle-ai/risor/v2"
	"github.com/deepnoodle-ai/risor/v2/pkg/compiler"
	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/deepnoodle-ai/risor/v2/pkg/parser"
)

// Option is a function that configures the execution of a script.
type Option func(*options)

// WithModule adds a custom module to the script execution environment.
func WithModule(module *object.Module) Option {
	return func(o *options) {
		if module == nil {
			return
		}
		o.modules = append(o.modules, module)
	}
}

// WithBuiltin adds a custom builtin function to the script execution environment.
func WithBuiltin(builtin *object.Builtin) Option {
	return func(options *options) {
		if builtin == nil {
			return
		}
		options.builtins = append(options.builtins, builtin)
	}
}

// WithFilesystem registers a filesystem for a specific URI scheme.
func WithFilesystem(scheme string, fs FS) Option {
	return func(options *options) {
		if fs == nil {
			return
		}
		if options.filesystems == nil {
			options.filesystems = make(map[string]FS)
		}
		options.filesystems[scheme] = fs
	}
}

// WithStdin sets the standard input file for the script.
func WithStdin(f File) Option {
	return func(o *options) {
		o.stdin = f
	}
}

// WithStdout sets the standard output file for the script.
func WithStdout(f File) Option {
	return func(o *options) {
		o.stdout = f
	}
}

// WithArgs sets the command line arguments for the script.
func WithArgs(args []string) Option {
	return func(o *options) {
		o.args = args
	}
}

// WithExitHandler sets the handler for os.exit calls in the script.
func WithExitHandler(handler ExitHandler) Option {
	return func(o *options) {
		o.exitHandler = handler
	}
}

// RunBytes executes a Ren script provided as a byte slice.
func RunBytes(ctx context.Context, b []byte, opts ...Option) error {
	reader := bytes.NewReader(b)
	return Run(ctx, reader, reader.Size(), opts...)
}

// RunFile executes a Ren script from a file.
func RunFile(ctx context.Context, filename string, opts ...Option) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	inf, err := f.Stat()
	if err != nil {
		return err
	}

	return Run(ctx, f, inf.Size(), opts...)
}

// Run executes a Ren script from an io.ReaderAt.
func Run(ctx context.Context, reader io.ReaderAt, size int64, opt ...Option) error {
	var opts options
	for _, o := range opt {
		o(&opts)
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

	builtins := opts.Builtins()

	env := make(map[string]any, len(builtins))
	maps.Copy(env, builtins)

	if !isOS(ctx) {
		ctx = WithOS(ctx, &osMiddleware{
			fs:          opts.Filesystems(),
			stdin:       opts.Stdin(),
			stdout:      opts.Stdout(),
			args:        opts.Args(),
			exitHandler: opts.ExitHandler(),
		})
	}

	ctx = WithImporter(ctx, newImporter(zr, opts.Modules(), env))

	_, err = risor.Run(
		ctx,
		code.ToBytecode(),
		risor.WithEnv(env),
		risor.WithFilename(code.Filename()),
	)
	if err != nil {
		return &Error{err}
	}

	return nil
}

func readEntrypoint(zr *zip.Reader) ([]byte, error) {
	f, err := zr.Open("entrypoint.json")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type options struct {
	stdin       File
	stdout      File
	args        []string
	exitHandler ExitHandler
	filesystems map[string]FS
	builtins    []*object.Builtin
	modules     []*object.Module
}

func (o *options) Builtins() map[string]any {
	result := make(map[string]any, len(o.builtins))
	for _, builtin := range o.builtins {
		result[builtin.Name()] = builtin
	}
	return result
}

func (o *options) Modules() map[string]*object.Module {
	result := make(map[string]*object.Module, len(o.modules))
	for _, module := range o.modules {
		result[module.Name().Value()] = module
	}
	return result
}

func (o *options) Filesystems() map[string]FS {
	result := make(map[string]FS, len(o.filesystems))
	result["file"] = &localFS{}
	maps.Copy(result, o.filesystems)
	return result
}

func (o *options) Stdin() File {
	if o.stdin != nil {
		return o.stdin
	}
	return os.Stdin
}

func (o *options) Stdout() File {
	if o.stdout != nil {
		return o.stdout
	}
	return os.Stdout
}

func (o *options) Args() []string {
	if o.args != nil {
		return o.args
	}
	return []string{}
}

func (o *options) ExitHandler() ExitHandler {
	if o.exitHandler != nil {
		return o.exitHandler
	}
	return func(code int) {}
}

// Error represents an error that occurred during script execution.
type Error struct {
	err error
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.err
}

// Error returns the error message.
func (e *Error) Error() string {
	if parserErr, ok := errors.AsType[parser.ParserError](e.err); ok {
		return parserErr.FriendlyErrorMessage()
	}
	return e.err.Error()
}
