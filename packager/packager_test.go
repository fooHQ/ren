package packager_test

import (
	"archive/zip"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/foohq/ren/packager"
)

func TestBuild(t *testing.T) {
	out := packager.NewFilename("helo")
	defer os.Remove(out)
	err := packager.Build("testdata/helo", out)
	require.NoError(t, err)

	zr, err := zip.OpenReader(out)
	require.NoError(t, err)
	defer zr.Close()

	var names []string
	for _, f := range zr.File {
		names = append(names, f.Name)
	}
	require.EqualValues(t, names, []string{"entrypoint.json", "main.json", "next2.json"})
}

func TestBuildMissingEntrypoint(t *testing.T) {
	out := packager.NewFilename("helo")
	err := packager.Build("testdata/noentrypoint", out)
	require.ErrorIs(t, err, packager.ErrMissingEntrypoint)
}

func TestBuildInvalidMain(t *testing.T) {
	out := packager.NewFilename("helo")
	err := packager.Build("testdata/inventrypoint", out)
	require.ErrorIs(t, err, packager.ErrMissingEntrypoint)
}
