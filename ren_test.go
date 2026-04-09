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
	}

	for _, tt := range tests {
		t.Run(tt.pth, func(t *testing.T) {
			out := packager.NewFilename(filepath.Base(tt.pth))
			defer func() {
				_ = os.Remove(out)
			}()

			// Build the package
			{
				var opts []packager.Option
				for _, o := range tt.builtins {
					opts = append(opts, packager.WithBuiltin(o))
				}
				for _, o := range tt.modules {
					opts = append(opts, packager.WithModule(o))
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
