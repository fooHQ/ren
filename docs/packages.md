# Package format

A Ren package is a `.zip` file containing compiled Risor scripts and any other
files copied from the source directory. Each Risor script (`.risor` or `.rsr`)
is parsed, compiled to bytecode, and written as a `.json` file; non-script files
are copied verbatim.

## Entrypoint

Every package must contain an `entrypoint.json`, produced by compiling an
`entrypoint.risor` (or `entrypoint.rsr`) in the root of the source directory.
This is the script the runtime executes to bootstrap the package. It is expected
to contain the bootstrapping code for whatever the package does.

```
scripts/cat/                 cat.zip
├── entrypoint.risor    ->   ├── entrypoint.json
└── lib/                     └── lib/
    └── read.risor               └── read.json
```

## Modules within a package

A script imports another module from the same package by its path, relative to
the package root:

```risor
read := import("lib/read")          // loads lib/read.json
read := import("package://lib/read") // explicit form
```

Each module is compiled as a self-contained function that returns a map of its
top-level names. The runtime runs a module at most once and caches its exports,
so repeated imports share a single instance. Import cycles are detected and
reported as an error.

See the [runtime reference](runtime.md#imports) for the `builtin://` scheme used
to reach modules registered with the runtime, and the
[examples](../examples) for complete packages.
