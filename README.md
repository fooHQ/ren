# ren

Ren is a runtime environment and script packager for [Risor](https://github.com/risor-io/risor) scripting language.
Ren can be used as a library to provide scripting capabilities to Go applications or as a utility to package/run/test
Risor scripts.

## Installation

```
go get github.com/foohq/ren
```

## Usage

```go
// The example code illustrates how to build a package and execute it.
// The code assumes some externalities, such as that specific directories must exist on the local filesystem.
// The packager creates a file cat.zip in the current working directory, if packaging was successful.
package main

import (
    "context"
    "os"

    "github.com/foohq/ren"
    "github.com/foohq/ren/filesystems/local"
    "github.com/foohq/ren/modules"
    renos "github.com/foohq/ren/os"
    "github.com/foohq/ren/packager"
    risoros "github.com/risor-io/risor/os"
)

func getArgs() []string {
    if len(os.Args) > 1 {
        return os.Args[1:]
    }
    return []string{}	
}

func main() {
    outputName := "cat.zip"
    // Create a package from a script found in a local directory.
    // Builder recognizes Risor scripts by .risor, .rsr extension.
    // Packager parses, compiles Risor scripts and writes compiled bytecode to .json files.
    // Files are then zipped to create a package file.
    err := packager.Build("/home/user/scripts/cat", outputName)
    if err != nil {
        panic(err)
    }

    // Instantiate a local filesystem.
    // Other compatible synthetic filesystems can be found in standalone repositories.
    // Visit https://github.com/fooHQ?q=filesystem.
    localFS, err := local.NewFS()
    if err != nil {
        panic(err)
    }

    // Instantiate OS which provides context to a running script as well as methods to interface with the operating system.
    ros := renos.New(
        renos.WithStdin(os.Stdin),
        renos.WithStdout(os.Stdout),
        renos.WithArgs(getArgs()),
        renos.WithFilesystems(map[string]risoros.FS{
            "file": localFS,
        }),
    )
    err = ren.RunFile(
        context.Background(),
        outputName,
        ros,
        ren.WithGlobals(modules.Globals()),
    )
    if err != nil {
        panic(err)
    }
}
```

## Command line utility

This repository provides a command line utility which can be used to build a package and execute it. The utility can either
be built from the source, or downloaded as a pre-compiled binary. For pre-compiled binaries, please refer to
repository's releases.

Installing from the source requires [devbox](https://www.jetify.com/docs/devbox/installing_devbox/):

```
$ git clone https://github.com/foohq/ren
$ cd ren
$ devbox run build
$ ./build/ren -h
```

## Package format

The package is a .zip file containing compiled Risor scripts encoded as .json files, and other files that were copied
from the source directory. A special file `entrypoint.json` must exist within the package. `entrypoint.json` file is used to bootstrap
a packaged script. `entrypoint.json` is created by compiling `entrypoint.risor` (or `entrypoint.rsr`) which must exist in the source
directory. The source file is expected to contain appropriate bootstrapping code tailored for the packaged script.

Please, see [examples](./examples) for more information.

## License

This module is distributed under the Apache License Version 2.0 found in the [LICENSE](./LICENSE) file.
