# ren

[![Go Reference](https://pkg.go.dev/badge/github.com/foohq/ren.svg)](https://pkg.go.dev/github.com/foohq/ren)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](LICENSE)

Ren is a runtime environment and script packager for the
[Risor](https://github.com/deepnoodle-ai/risor) scripting language. It compiles a
directory of Risor scripts into a single package file and runs it inside a
sandboxed runtime.

Use Ren as a **library** to add scripting to a Go application, or as a
**command-line utility** to build, run, and distribute Risor scripts.

## Installation

As a library:

```
go get github.com/foohq/ren
```

As a command-line tool — download a pre-compiled binary for your platform from
the [GitHub releases](https://github.com/foohq/ren/releases), or build from
source with [devbox](https://www.jetify.com/docs/devbox/installing_devbox/):

```
$ git clone https://github.com/foohq/ren
$ cd ren
$ devbox run build
$ ./build/ren -h
```

## Quickstart

Package a directory of scripts and run it from the command line:

```
$ ren build ./examples/hello    # produces hello.zip
$ ren run hello.zip
```

The same, embedded in a Go program:

```go
package main

import (
	"context"
	"os"

	"github.com/foohq/ren"
	"github.com/foohq/ren/builtins"
	"github.com/foohq/ren/modules"
	"github.com/foohq/ren/packager"
)

func main() {
	// Compile the Risor scripts in ./scripts/hello into a package file.
	// The builtins are supplied so the compiler recognises them at build time.
	var buildOpts []packager.Option
	for _, b := range builtins.Builtins() {
		buildOpts = append(buildOpts, packager.WithBuiltin(b))
	}
	if err := packager.Build("./scripts/hello", "hello.zip", buildOpts...); err != nil {
		panic(err)
	}

	// Run the package, exposing Ren's builtins and modules to the script.
	runOpts := []ren.Option{
		ren.WithStdout(os.Stdout),
		ren.WithArgs(os.Args[1:]),
	}
	for _, b := range builtins.Builtins() {
		runOpts = append(runOpts, ren.WithBuiltin(b))
	}
	for _, m := range modules.Modules() {
		runOpts = append(runOpts, ren.WithModule(m))
	}
	if err := ren.RunFile(context.Background(), "hello.zip", runOpts...); err != nil {
		panic(err)
	}
}
```

## Documentation

- [Command-line utility](docs/cli.md) — building and running packages with the `ren` CLI.
- [Using Ren as a library](docs/library.md) — embedding Ren in a Go application, and the runtime options.
- [Runtime reference](docs/runtime.md) — the builtins and modules available to scripts.
- [Package format](docs/packages.md) — what a package contains and how imports resolve.
- [Examples](examples) — complete, runnable packages.

The Go API reference is on
[pkg.go.dev](https://pkg.go.dev/github.com/foohq/ren).

## License

This module is distributed under the Apache License Version 2.0 found in the
[LICENSE](./LICENSE) file.
