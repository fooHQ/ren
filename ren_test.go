package ren_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/stretchr/testify/require"

	"github.com/foohq/ren"
	"github.com/foohq/ren/builtins"
	"github.com/foohq/ren/modules"
	"github.com/foohq/ren/packager"
)

func TestRunFile(t *testing.T) {
	tests := []struct {
		pth      string
		builtins map[string]*object.Builtin
		modules  map[string]*object.Module
	}{
		{
			pth:      "examples/hello",
			builtins: builtins.Builtins(),
			modules:  modules.Modules(),
		},
		{
			pth:      "examples/imports",
			builtins: builtins.Builtins(),
			modules:  modules.Modules(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.pth, func(t *testing.T) {
			out := packager.NewFilename(filepath.Base(tt.pth))
			defer func() {
				_ = os.Remove(out)
			}()

			// Build the package. Modules are not compiled as globals; scripts
			// reach them through import("builtin://...").
			{
				var opts []packager.Option
				for _, o := range tt.builtins {
					opts = append(opts, packager.WithBuiltin(o))
				}
				err := packager.Build(tt.pth, out, opts...)
				require.NoError(t, err)
			}

			// Run the build package
			{
				var opts []ren.Option
				for _, o := range tt.builtins {
					opts = append(opts, ren.WithBuiltin(o))
				}
				for _, o := range tt.modules {
					opts = append(opts, ren.WithModule(o))
				}
				err := ren.RunFile(context.Background(), out, opts...)
				require.NoError(t, err)
			}
		})
	}
}

// TestModuleErrorLineNumber verifies that a runtime error inside an imported
// module reports the module's real source line and column, despite the packager
// wrapping the module in a synthetic function.
func TestModuleErrorLineNumber(t *testing.T) {
	srcDir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(srcDir, "lib"), 0755))
	require.NoError(t, os.WriteFile(
		filepath.Join(srcDir, "entrypoint.risor"),
		[]byte(`const boom = import("lib/boom")`+"\n"),
		0644,
	))
	// The runtime error (calling an int) is on line 3 of the module.
	require.NoError(t, os.WriteFile(
		filepath.Join(srcDir, "lib", "boom.risor"),
		[]byte("const a = 1\nconst b = 2\nconst x = a(5)\n"),
		0644,
	))

	err := packAndRun(t, srcDir)
	require.Error(t, err)
	require.Contains(t, err.Error(), "(3:13)")
}

// TestModuleImportCycle verifies that a circular import is rejected with a
// cycle error rather than looping or overflowing the stack.
func TestModuleImportCycle(t *testing.T) {
	srcDir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(srcDir, "lib"), 0755))
	require.NoError(t, os.WriteFile(
		filepath.Join(srcDir, "entrypoint.risor"),
		[]byte(`const a = import("lib/a")`+"\n"),
		0644,
	))
	// lib/a imports lib/b, which imports lib/a -> cycle.
	require.NoError(t, os.WriteFile(
		filepath.Join(srcDir, "lib", "a.risor"),
		[]byte(`const b = import("lib/b")`+"\n"),
		0644,
	))
	require.NoError(t, os.WriteFile(
		filepath.Join(srcDir, "lib", "b.risor"),
		[]byte(`const a = import("lib/a")`+"\n"),
		0644,
	))

	err := packAndRun(t, srcDir)
	require.Error(t, err)
	require.Contains(t, err.Error(), "import cycle detected")
}

// packAndRun builds the package rooted at srcDir with the standard builtins and
// runs it, returning any execution error.
func packAndRun(t *testing.T, srcDir string) error {
	t.Helper()

	out := filepath.Join(t.TempDir(), packager.NewFilename("pkg"))

	var buildOpts []packager.Option
	for _, o := range builtins.Builtins() {
		buildOpts = append(buildOpts, packager.WithBuiltin(o))
	}
	require.NoError(t, packager.Build(srcDir, out, buildOpts...))

	var runOpts []ren.Option
	for _, o := range builtins.Builtins() {
		runOpts = append(runOpts, ren.WithBuiltin(o))
	}
	return ren.RunFile(context.Background(), out, runOpts...)
}
