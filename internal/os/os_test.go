package os_test

import (
	"testing"

	risoros "github.com/risor-io/risor/os"
	"github.com/stretchr/testify/require"

	"github.com/foohq/ren/internal/os"
	"github.com/foohq/ren/testutils"
)

// TODO: add the rest of the tests!

func TestOS_Args(t *testing.T) {
	args := []string{
		"first",
		"second",
		"third",
	}
	o := os.New(
		os.WithArgs(args),
	)
	actualArgs := o.Args()
	require.Equal(t, args, actualArgs)
}

func TestOS_Create(t *testing.T) {
	testCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(testCh),
		"test": testutils.NewFS(testCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		input  string
		result testutils.CreateResult
	}{
		{
			input: "test://private/",
			result: testutils.CreateResult{
				Name: "/",
			},
		},
		{
			input: "test://private/data/form.txt",
			result: testutils.CreateResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "test://private/data/../form.txt",
			result: testutils.CreateResult{
				Name: "/form.txt",
			},
		},
		{
			input: "test://private/data/../../form.txt",
			result: testutils.CreateResult{
				Name: "/form.txt",
			},
		},
		{
			input: "/",
			result: testutils.CreateResult{
				Name: "/",
			},
		},
		{
			input: "/data/form.txt",
			result: testutils.CreateResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "/data/../form.txt",
			result: testutils.CreateResult{
				Name: "/form.txt",
			},
		},
		{
			input: "/data/../../form.txt",
			result: testutils.CreateResult{
				Name: "/form.txt",
			},
		},
		{
			input: "file:///data/form.txt",
			result: testutils.CreateResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "../data/form.txt",
			result: testutils.CreateResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "./data/form.txt",
			result: testutils.CreateResult{
				Name: "/ren/data/form.txt",
			},
		},
	}

	for i, test := range tests {
		_, err := o.Create(test.input)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-testCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_Mkdir(t *testing.T) {
	resultCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(resultCh),
		"test": testutils.NewFS(resultCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		input  string
		result testutils.MkdirResult
	}{
		{
			input: "test://private/",
			result: testutils.MkdirResult{
				Name: "/",
				Perm: 0777,
			},
		},
		{
			input: "test://private/data/form",
			result: testutils.MkdirResult{
				Name: "/data/form",
				Perm: 0777,
			},
		},
		{
			input: "test://private/data/../form",
			result: testutils.MkdirResult{
				Name: "/form",
				Perm: 0777,
			},
		},
		{
			input: "test://private/data/../../form",
			result: testutils.MkdirResult{
				Name: "/form",
				Perm: 0777,
			},
		},
		{
			input: "/",
			result: testutils.MkdirResult{
				Name: "/",
				Perm: 0777,
			},
		},
		{
			input: "/data/form",
			result: testutils.MkdirResult{
				Name: "/data/form",
				Perm: 0777,
			},
		},
		{
			input: "/data/../form",
			result: testutils.MkdirResult{
				Name: "/form",
				Perm: 0777,
			},
		},
		{
			input: "/data/../../form",
			result: testutils.MkdirResult{
				Name: "/form",
				Perm: 0777,
			},
		},
		{
			input: "file:///data/form",
			result: testutils.MkdirResult{
				Name: "/data/form",
				Perm: 0777,
			},
		},
		{
			input: "../data/form",
			result: testutils.MkdirResult{
				Name: "/data/form",
				Perm: 0777,
			},
		},
		{
			input: "./data/form.txt",
			result: testutils.MkdirResult{
				Name: "/ren/data/form.txt",
				Perm: 0777,
			},
		},
	}

	for i, test := range tests {
		err := o.Mkdir(test.input, 0777)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-resultCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_MkdirAll(t *testing.T) {
	resultCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(resultCh),
		"test": testutils.NewFS(resultCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		input  string
		result testutils.MkdirAllResult
	}{
		{
			input: "test://private/",
			result: testutils.MkdirAllResult{
				Path: "/",
				Perm: 0777,
			},
		},
		{
			input: "test://private/data/form",
			result: testutils.MkdirAllResult{
				Path: "/data/form",
				Perm: 0777,
			},
		},
		{
			input: "test://private/data/../form",
			result: testutils.MkdirAllResult{
				Path: "/form",
				Perm: 0777,
			},
		},
		{
			input: "test://private/data/../../form",
			result: testutils.MkdirAllResult{
				Path: "/form",
				Perm: 0777,
			},
		},
		{
			input: "/",
			result: testutils.MkdirAllResult{
				Path: "/",
				Perm: 0777,
			},
		},
		{
			input: "/data/form",
			result: testutils.MkdirAllResult{
				Path: "/data/form",
				Perm: 0777,
			},
		},
		{
			input: "/data/../form",
			result: testutils.MkdirAllResult{
				Path: "/form",
				Perm: 0777,
			},
		},
		{
			input: "/data/../../form",
			result: testutils.MkdirAllResult{
				Path: "/form",
				Perm: 0777,
			},
		},
		{
			input: "file:///data/form",
			result: testutils.MkdirAllResult{
				Path: "/data/form",
				Perm: 0777,
			},
		},
		{
			input: "../data/form",
			result: testutils.MkdirAllResult{
				Path: "/data/form",
				Perm: 0777,
			},
		},
		{
			input: "./data/form.txt",
			result: testutils.MkdirAllResult{
				Path: "/ren/data/form.txt",
				Perm: 0777,
			},
		},
	}

	for i, test := range tests {
		err := o.MkdirAll(test.input, 0777)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-resultCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_Open(t *testing.T) {
	resultCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(resultCh),
		"test": testutils.NewFS(resultCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		input  string
		result testutils.OpenResult
	}{
		{
			input: "test://private/",
			result: testutils.OpenResult{
				Name: "/",
			},
		},
		{
			input: "test://private/data/form.txt",
			result: testutils.OpenResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "test://private/data/../form.txt",
			result: testutils.OpenResult{
				Name: "/form.txt",
			},
		},
		{
			input: "test://private/data/../../form.txt",
			result: testutils.OpenResult{
				Name: "/form.txt",
			},
		},
		{
			input: "/",
			result: testutils.OpenResult{
				Name: "/",
			},
		},
		{
			input: "/data/form.txt",
			result: testutils.OpenResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "/data/../form.txt",
			result: testutils.OpenResult{
				Name: "/form.txt",
			},
		},
		{
			input: "/data/../../form.txt",
			result: testutils.OpenResult{
				Name: "/form.txt",
			},
		},
		{
			input: "file:///data/form.txt",
			result: testutils.OpenResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "../data/form.txt",
			result: testutils.OpenResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "./data/form.txt",
			result: testutils.OpenResult{
				Name: "/ren/data/form.txt",
			},
		},
	}

	for i, test := range tests {
		_, err := o.Open(test.input)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-resultCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_OpenFile(t *testing.T) {
	resultCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(resultCh),
		"test": testutils.NewFS(resultCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		input  string
		result testutils.OpenFileResult
	}{
		{
			input: "test://private/",
			result: testutils.OpenFileResult{
				Name: "/",
				Flag: 1313,
				Perm: 0777,
			},
		},
		{
			input: "test://private/data/form.txt",
			result: testutils.OpenFileResult{
				Name: "/data/form.txt",
				Flag: 1313,
				Perm: 0777,
			},
		},
		{
			input: "test://private/data/../form.txt",
			result: testutils.OpenFileResult{
				Name: "/form.txt",
				Flag: 1313,
				Perm: 0777,
			},
		},
		{
			input: "test://private/data/../../form.txt",
			result: testutils.OpenFileResult{
				Name: "/form.txt",
				Flag: 1313,
				Perm: 0777,
			},
		},
		{
			input: "/",
			result: testutils.OpenFileResult{
				Name: "/",
				Flag: 1313,
				Perm: 0777,
			},
		},
		{
			input: "/data/form.txt",
			result: testutils.OpenFileResult{
				Name: "/data/form.txt",
				Flag: 1313,
				Perm: 0777,
			},
		},
		{
			input: "/data/../form.txt",
			result: testutils.OpenFileResult{
				Name: "/form.txt",
				Flag: 1313,
				Perm: 0777,
			},
		},
		{
			input: "/data/../../form.txt",
			result: testutils.OpenFileResult{
				Name: "/form.txt",
				Flag: 1313,
				Perm: 0777,
			},
		},
		{
			input: "file:///data/form.txt",
			result: testutils.OpenFileResult{
				Name: "/data/form.txt",
				Flag: 1313,
				Perm: 0777,
			},
		},
		{
			input: "../data/form.txt",
			result: testutils.OpenFileResult{
				Name: "/data/form.txt",
				Flag: 1313,
				Perm: 0777,
			},
		},
		{
			input: "./data/form.txt",
			result: testutils.OpenFileResult{
				Name: "/ren/data/form.txt",
				Flag: 1313,
				Perm: 0777,
			},
		},
	}

	for i, test := range tests {
		_, err := o.OpenFile(test.input, 1313, 0777)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-resultCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_ReadFile(t *testing.T) {
	resultCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(resultCh),
		"test": testutils.NewFS(resultCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		input  string
		result testutils.ReadFileResult
	}{
		{
			input: "test://private/",
			result: testutils.ReadFileResult{
				Name: "/",
			},
		},
		{
			input: "test://private/data/form.txt",
			result: testutils.ReadFileResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "test://private/data/../form.txt",
			result: testutils.ReadFileResult{
				Name: "/form.txt",
			},
		},
		{
			input: "test://private/data/../../form.txt",
			result: testutils.ReadFileResult{
				Name: "/form.txt",
			},
		},
		{
			input: "/",
			result: testutils.ReadFileResult{
				Name: "/",
			},
		},
		{
			input: "/data/form.txt",
			result: testutils.ReadFileResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "/data/../form.txt",
			result: testutils.ReadFileResult{
				Name: "/form.txt",
			},
		},
		{
			input: "/data/../../form.txt",
			result: testutils.ReadFileResult{
				Name: "/form.txt",
			},
		},
		{
			input: "file:///data/form.txt",
			result: testutils.ReadFileResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "../data/form.txt",
			result: testutils.ReadFileResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "./data/form.txt",
			result: testutils.ReadFileResult{
				Name: "/ren/data/form.txt",
			},
		},
	}

	for i, test := range tests {
		_, err := o.ReadFile(test.input)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-resultCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_Remove(t *testing.T) {
	resultCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(resultCh),
		"test": testutils.NewFS(resultCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		input  string
		result testutils.RemoveResult
	}{
		{
			input: "test://private/",
			result: testutils.RemoveResult{
				Name: "/",
			},
		},
		{
			input: "test://private/data/form.txt",
			result: testutils.RemoveResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "test://private/data/../form.txt",
			result: testutils.RemoveResult{
				Name: "/form.txt",
			},
		},
		{
			input: "test://private/data/../../form.txt",
			result: testutils.RemoveResult{
				Name: "/form.txt",
			},
		},
		{
			input: "/",
			result: testutils.RemoveResult{
				Name: "/",
			},
		},
		{
			input: "/data/form.txt",
			result: testutils.RemoveResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "/data/../form.txt",
			result: testutils.RemoveResult{
				Name: "/form.txt",
			},
		},
		{
			input: "/data/../../form.txt",
			result: testutils.RemoveResult{
				Name: "/form.txt",
			},
		},
		{
			input: "file:///data/form.txt",
			result: testutils.RemoveResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "../data/form.txt",
			result: testutils.RemoveResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "./data/form.txt",
			result: testutils.RemoveResult{
				Name: "/ren/data/form.txt",
			},
		},
	}

	for i, test := range tests {
		err := o.Remove(test.input)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-resultCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_RemoveAll(t *testing.T) {
	resultCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(resultCh),
		"test": testutils.NewFS(resultCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		input  string
		result testutils.RemoveAllResult
	}{
		{
			input: "test://private/",
			result: testutils.RemoveAllResult{
				Path: "/",
			},
		},
		{
			input: "test://private/data/form.txt",
			result: testutils.RemoveAllResult{
				Path: "/data/form.txt",
			},
		},
		{
			input: "test://private/data/../form.txt",
			result: testutils.RemoveAllResult{
				Path: "/form.txt",
			},
		},
		{
			input: "test://private/data/../../form.txt",
			result: testutils.RemoveAllResult{
				Path: "/form.txt",
			},
		},
		{
			input: "/",
			result: testutils.RemoveAllResult{
				Path: "/",
			},
		},
		{
			input: "/data/form.txt",
			result: testutils.RemoveAllResult{
				Path: "/data/form.txt",
			},
		},
		{
			input: "/data/../form.txt",
			result: testutils.RemoveAllResult{
				Path: "/form.txt",
			},
		},
		{
			input: "/data/../../form.txt",
			result: testutils.RemoveAllResult{
				Path: "/form.txt",
			},
		},
		{
			input: "file:///data/form.txt",
			result: testutils.RemoveAllResult{
				Path: "/data/form.txt",
			},
		},
		{
			input: "../data/form.txt",
			result: testutils.RemoveAllResult{
				Path: "/data/form.txt",
			},
		},
		{
			input: "./data/form.txt",
			result: testutils.RemoveAllResult{
				Path: "/ren/data/form.txt",
			},
		},
	}

	for i, test := range tests {
		err := o.RemoveAll(test.input)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-resultCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_Rename(t *testing.T) {
	resultCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(resultCh),
		"test": testutils.NewFS(resultCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		src    string
		dst    string
		result testutils.RenameResult
	}{
		{
			src: "test://private/foo.txt",
			dst: "test://private/bar.txt",
			result: testutils.RenameResult{
				OldPath: "/foo.txt",
				NewPath: "/bar.txt",
			},
		},
		{
			src: "/private/foo.txt",
			dst: "/private/bar.txt",
			result: testutils.RenameResult{
				OldPath: "/private/foo.txt",
				NewPath: "/private/bar.txt",
			},
		},
		{
			src: "/private/foo.txt",
			dst: "../bar.txt",
			result: testutils.RenameResult{
				OldPath: "/private/foo.txt",
				NewPath: "/bar.txt",
			},
		},
		{
			src: "./private/foo.txt",
			dst: "bar.txt",
			result: testutils.RenameResult{
				OldPath: "/ren/private/foo.txt",
				NewPath: "/ren/bar.txt",
			},
		},
		{
			src: "../foo.txt",
			dst: "./bar.txt",
			result: testutils.RenameResult{
				OldPath: "/foo.txt",
				NewPath: "/ren/bar.txt",
			},
		},
	}

	for i, test := range tests {
		err := o.Rename(test.src, test.dst)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-resultCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_Rename_ErrCrossingFSBoundaries(t *testing.T) {
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(nil),
		"test": testutils.NewFS(nil),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		src string
		dst string
	}{
		{
			src: "test://private/foo.txt",
			dst: "./private/bar.txt",
		},
		{
			src: "./private/bar.txt",
			dst: "test://private/foo.txt",
		},
	}

	for i, test := range tests {
		err := o.Rename(test.src, test.dst)
		require.ErrorIs(t, err, os.ErrCrossingFSBoundaries, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_Stat(t *testing.T) {
	resultCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(resultCh),
		"test": testutils.NewFS(resultCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		input  string
		result testutils.StatResult
	}{
		{
			input: "test://private/",
			result: testutils.StatResult{
				Name: "/",
			},
		},
		{
			input: "test://private/data/form.txt",
			result: testutils.StatResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "test://private/data/../form.txt",
			result: testutils.StatResult{
				Name: "/form.txt",
			},
		},
		{
			input: "test://private/data/../../form.txt",
			result: testutils.StatResult{
				Name: "/form.txt",
			},
		},
		{
			input: "/",
			result: testutils.StatResult{
				Name: "/",
			},
		},
		{
			input: "/data/form.txt",
			result: testutils.StatResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "/data/../form.txt",
			result: testutils.StatResult{
				Name: "/form.txt",
			},
		},
		{
			input: "/data/../../form.txt",
			result: testutils.StatResult{
				Name: "/form.txt",
			},
		},
		{
			input: "file:///data/form.txt",
			result: testutils.StatResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "../data/form.txt",
			result: testutils.StatResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "./data/form.txt",
			result: testutils.StatResult{
				Name: "/ren/data/form.txt",
			},
		},
	}

	for i, test := range tests {
		_, err := o.Stat(test.input)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-resultCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_Symlink(t *testing.T) {
	resultCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(resultCh),
		"test": testutils.NewFS(resultCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		src    string
		dst    string
		result testutils.SymlinkResult
	}{
		{
			src: "test://private/foo.txt",
			dst: "test://private/bar.txt",
			result: testutils.SymlinkResult{
				OldName: "/foo.txt",
				NewName: "/bar.txt",
			},
		},
		{
			src: "/private/foo.txt",
			dst: "/private/bar.txt",
			result: testutils.SymlinkResult{
				OldName: "/private/foo.txt",
				NewName: "/private/bar.txt",
			},
		},
		{
			src: "/private/foo.txt",
			dst: "../bar.txt",
			result: testutils.SymlinkResult{
				OldName: "/private/foo.txt",
				NewName: "/bar.txt",
			},
		},
		{
			src: "./private/foo.txt",
			dst: "bar.txt",
			result: testutils.SymlinkResult{
				OldName: "/ren/private/foo.txt",
				NewName: "/ren/bar.txt",
			},
		},
		{
			src: "../foo.txt",
			dst: "./bar.txt",
			result: testutils.SymlinkResult{
				OldName: "/foo.txt",
				NewName: "/ren/bar.txt",
			},
		},
	}

	for i, test := range tests {
		err := o.Symlink(test.src, test.dst)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-resultCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_Symlink_ErrCrossingFSBoundaries(t *testing.T) {
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(nil),
		"test": testutils.NewFS(nil),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		src string
		dst string
	}{
		{
			src: "test://private/foo.txt",
			dst: "./private/bar.txt",
		},
		{
			src: "./private/bar.txt",
			dst: "test://private/foo.txt",
		},
	}

	for i, test := range tests {
		err := o.Symlink(test.src, test.dst)
		require.ErrorIs(t, err, os.ErrCrossingFSBoundaries, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_WriteFile(t *testing.T) {
	resultCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(resultCh),
		"test": testutils.NewFS(resultCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		input  string
		result testutils.WriteFileResult
	}{
		{
			input: "test://private/",
			result: testutils.WriteFileResult{
				Name: "/",
				Data: []byte("test"),
				Perm: 0777,
			},
		},
		{
			input: "test://private/data/form.txt",
			result: testutils.WriteFileResult{
				Name: "/data/form.txt",
				Data: []byte("test"),
				Perm: 0777,
			},
		},
		{
			input: "test://private/data/../form.txt",
			result: testutils.WriteFileResult{
				Name: "/form.txt",
				Data: []byte("test"),
				Perm: 0777,
			},
		},
		{
			input: "test://private/data/../../form.txt",
			result: testutils.WriteFileResult{
				Name: "/form.txt",
				Data: []byte("test"),
				Perm: 0777,
			},
		},
		{
			input: "/",
			result: testutils.WriteFileResult{
				Name: "/",
				Data: []byte("test"),
				Perm: 0777,
			},
		},
		{
			input: "/data/form.txt",
			result: testutils.WriteFileResult{
				Name: "/data/form.txt",
				Data: []byte("test"),
				Perm: 0777,
			},
		},
		{
			input: "/data/../form.txt",
			result: testutils.WriteFileResult{
				Name: "/form.txt",
				Data: []byte("test"),
				Perm: 0777,
			},
		},
		{
			input: "/data/../../form.txt",
			result: testutils.WriteFileResult{
				Name: "/form.txt",
				Data: []byte("test"),
				Perm: 0777,
			},
		},
		{
			input: "file:///data/form.txt",
			result: testutils.WriteFileResult{
				Name: "/data/form.txt",
				Data: []byte("test"),
				Perm: 0777,
			},
		},
		{
			input: "../data/form.txt",
			result: testutils.WriteFileResult{
				Name: "/data/form.txt",
				Data: []byte("test"),
				Perm: 0777,
			},
		},
		{
			input: "./data/form.txt",
			result: testutils.WriteFileResult{
				Name: "/ren/data/form.txt",
				Data: []byte("test"),
				Perm: 0777,
			},
		},
	}

	for i, test := range tests {
		err := o.WriteFile(test.input, []byte("test"), 0777)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-resultCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_ReadDir(t *testing.T) {
	resultCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(resultCh),
		"test": testutils.NewFS(resultCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		input  string
		result testutils.ReadDirResult
	}{
		{
			input: "test://private/",
			result: testutils.ReadDirResult{
				Name: "/",
			},
		},
		{
			input: "test://private/data/form",
			result: testutils.ReadDirResult{
				Name: "/data/form",
			},
		},
		{
			input: "test://private/data/../form",
			result: testutils.ReadDirResult{
				Name: "/form",
			},
		},
		{
			input: "test://private/data/../../form",
			result: testutils.ReadDirResult{
				Name: "/form",
			},
		},
		{
			input: "/",
			result: testutils.ReadDirResult{
				Name: "/",
			},
		},
		{
			input: "/data/form.txt",
			result: testutils.ReadDirResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "/data/../form.txt",
			result: testutils.ReadDirResult{
				Name: "/form.txt",
			},
		},
		{
			input: "/data/../../form.txt",
			result: testutils.ReadDirResult{
				Name: "/form.txt",
			},
		},
		{
			input: "file:///data/form.txt",
			result: testutils.ReadDirResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "../data/form.txt",
			result: testutils.ReadDirResult{
				Name: "/data/form.txt",
			},
		},
		{
			input: "./data/form.txt",
			result: testutils.ReadDirResult{
				Name: "/ren/data/form.txt",
			},
		},
	}

	for i, test := range tests {
		_, err := o.ReadDir(test.input)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-resultCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}

func TestOS_WalkDir(t *testing.T) {
	resultCh := make(chan any, 1)
	fss := map[string]risoros.FS{
		"file": testutils.NewFS(resultCh),
		"test": testutils.NewFS(resultCh),
	}
	o := os.New(
		os.WithWorkDir("/ren"),
		os.WithFilesystems(fss),
	)

	tests := []struct {
		input  string
		result testutils.WalkDirResult
	}{
		{
			input: "test://private/",
			result: testutils.WalkDirResult{
				Root: "/",
				Fn:   nil,
			},
		},
		{
			input: "test://private/data/form",
			result: testutils.WalkDirResult{
				Root: "/data/form",
				Fn:   nil,
			},
		},
		{
			input: "test://private/data/../form",
			result: testutils.WalkDirResult{
				Root: "/form",
				Fn:   nil,
			},
		},
		{
			input: "test://private/data/../../form",
			result: testutils.WalkDirResult{
				Root: "/form",
				Fn:   nil,
			},
		},
		{
			input: "/",
			result: testutils.WalkDirResult{
				Root: "/",
			},
		},
		{
			input: "/data/form.txt",
			result: testutils.WalkDirResult{
				Root: "/data/form.txt",
			},
		},
		{
			input: "/data/../form.txt",
			result: testutils.WalkDirResult{
				Root: "/form.txt",
			},
		},
		{
			input: "/data/../../form.txt",
			result: testutils.WalkDirResult{
				Root: "/form.txt",
			},
		},
		{
			input: "file:///data/form.txt",
			result: testutils.WalkDirResult{
				Root: "/data/form.txt",
			},
		},
		{
			input: "../data/form.txt",
			result: testutils.WalkDirResult{
				Root: "/data/form.txt",
			},
		},
		{
			input: "./data/form.txt",
			result: testutils.WalkDirResult{
				Root: "/ren/data/form.txt",
			},
		},
	}

	for i, test := range tests {
		err := o.WalkDir(test.input, nil)
		require.NoError(t, err, "test %d/%d", i+1, len(tests))

		result := <-resultCh
		require.Equal(t, test.result, result, "test %d/%d", i+1, len(tests))
	}
}
