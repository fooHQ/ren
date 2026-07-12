# Using Ren as a library

Ren embeds Risor scripting into a Go application. You build a package from a
directory of scripts, then run it with a runtime you assemble from options.
Full API reference lives on
[pkg.go.dev](https://pkg.go.dev/github.com/foohq/ren); this guide covers the
shape of the two halves.

```
go get github.com/foohq/ren
```

## Packaging

`packager.Build` compiles the Risor scripts in a source directory into a package
file. Pass the builtins the scripts use so the compiler resolves them at build
time.

```go
var opts []packager.Option
for _, b := range builtins.Builtins() {
	opts = append(opts, packager.WithBuiltin(b))
}
err := packager.Build("./scripts/hello", "hello.zip", opts...)
```

The source directory must contain an `entrypoint.risor` (or `entrypoint.rsr`).
See [Package format](packages.md) for what ends up inside the `.zip`.

## Running

`ren.RunFile` (or `ren.Run` / `ren.RunBytes`) executes a package. The runtime is
configured entirely through options:

| Option | Purpose |
|---|---|
| `WithBuiltin(b)` | Register a global builtin. |
| `WithModule(m)` | Register a module importable via `builtin://<name>`. |
| `WithFilesystem(scheme, fs)` | Back a URL scheme (e.g. `file`) with a filesystem the `fs`/`os` modules operate on. |
| `WithStdin(f)` / `WithStdout(f)` | Wire the script's standard streams. |
| `WithArgs(args)` | Set the arguments returned by `os.args`. |
| `WithExitHandler(fn)` | Handle `os.exit`. |

```go
opts := []ren.Option{
	ren.WithStdout(os.Stdout),
	ren.WithArgs(os.Args[1:]),
}
for _, b := range builtins.Builtins() {
	opts = append(opts, ren.WithBuiltin(b))
}
for _, m := range modules.Modules() {
	opts = append(opts, ren.WithModule(m))
}
err := ren.RunFile(context.Background(), "hello.zip", opts...)
```

`builtins.Builtins()` and `modules.Modules()` return Ren's default runtime — the
same set the [`ren` CLI](cli.md) uses. You can register your own, or a curated
subset, instead. For the functions and modules they contain, see the
[runtime reference](runtime.md).

## Filesystems

Modules like `fs` and `os` never touch the host directly; they dispatch through
filesystems you register per URL scheme with `WithFilesystem`. A script reading
`file://data/input.txt` is served by whatever filesystem is registered for the
`file` scheme. Ready-made filesystems (local disk, in-memory, and more) live in
standalone repositories — see
[github.com/fooHQ?q=filesystem](https://github.com/orgs/fooHQ/repositories?q=filesystem).

```go
localFS, _ := local.NewFS()
opts = append(opts, ren.WithFilesystem("file", localFS))
```
