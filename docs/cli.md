# `ren` command-line utility

`ren` packages Risor scripts and runs the resulting packages. Build it from
source with [devbox](https://www.jetify.com/docs/devbox/installing_devbox/):

```
$ git clone https://github.com/foohq/ren
$ cd ren
$ devbox run build
$ ./build/ren -h
```

Pre-compiled binaries are attached to the repository's
[releases](https://github.com/foohq/ren/releases).

```
ren [global options] <command> [command options]

COMMANDS:
   build   Package Risor scripts
   run     Run Risor script from a package
```

## `ren build`

```
ren build [-o <output>] <dir>
```

Compiles the Risor scripts in `<dir>` into a package (a `.zip`; see
[Package format](packages.md)). The directory must contain an
`entrypoint.risor` (or `entrypoint.rsr`) that bootstraps the script.

| Flag | Description |
|---|---|
| `-o`, `--output <file>` | Output file. Defaults to `<dir>.zip` in the current directory. |

```
$ ren build ./scripts/cat
$ ren build -o cat.zip ./scripts/cat
```

Scripts are built with Ren's [global builtins](runtime.md#global-builtins)
available, so the compiler recognises `print`, `import`, `pack`, and the rest.

## `ren run`

```
ren run <pkg> [arg ...]
```

Runs the package `<pkg>`, forwarding any trailing arguments to the script (where
they are available through `os.args`). The script executes with Ren's global
builtins and every [built-in module](runtime.md#modules) registered.

```
$ ren run cat.zip file.txt
```

To make packaged scripts reachable from other packages, or to expose host files
through the `fs` module, use the library API — the CLI runs packages with the
default runtime only. See the [library guide](library.md).
