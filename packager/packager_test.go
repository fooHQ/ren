package packager_test

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/stretchr/testify/require"

	"github.com/foohq/ren/packager"
)

func TestBuild(t *testing.T) {
	tests := []struct {
		pth       string
		builtins  []*object.Builtin
		wantFiles []string
		wantErr   error
	}{
		{
			pth:     "testdata/bad_entrypoint",
			wantErr: packager.ErrMissingEntrypoint,
		},
		{
			pth:     "testdata/missing_entrypoint",
			wantErr: packager.ErrMissingEntrypoint,
		},
		{
			pth: "testdata/hello",
			builtins: []*object.Builtin{
				object.NewBuiltin("print", func(ctx context.Context, args ...object.Object) (object.Object, error) {
					return object.Nil, nil
				}),
			},
			wantFiles: []string{
				"entrypoint.json",
				"main.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.pth, func(t *testing.T) {
			out := packager.NewFilename(filepath.Base(tt.pth))
			defer func() {
				_ = os.Remove(out)
			}()

			var opts []packager.Option
			for _, o := range tt.builtins {
				opts = append(opts, packager.WithBuiltin(o))
			}

			err := packager.Build(tt.pth, out, opts...)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)

			zr, err := zip.OpenReader(out)
			require.NoError(t, err)
			defer func() {
				_ = zr.Close()
			}()

			var names []string
			for _, f := range zr.File {
				names = append(names, f.Name)
			}
			require.EqualValues(t, tt.wantFiles, names)
		})
	}
}
