package importer

import (
	"context"
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestImporter_Import(t *testing.T) {
	// Create test filesystem with the bar.risor fixture
	testFS := fstest.MapFS{
		"foo/bar.json": &fstest.MapFile{
			Data: []byte(`{"code":[{"id":"__main__","name":"__main__","symbol_table_id":"root","instructions":[24,0,71,0,33,0,80],"constants":[{"type":"function","value":{"id":"1","name":"bar","parameters":[],"defaults":[]}}],"source":"func bar() { return 987 }"},{"id":"__main__.0","name":"bar","parent_id":"__main__","symbol_table_id":"root.0","function_id":"1","instructions":[24,0,4],"constants":[{"type":"int","value":987}],"source":"return 987"}],"symbol_table":{"id":"root","symbols":[{"name":"bar","index":0,"is_constant":true}],"symbols_by_name":{"bar":{"name":"bar","index":0,"is_constant":true}},"children":[{"id":"root.0","symbols":[{"name":"bar","index":0,"is_constant":true}],"symbols_by_name":{"bar":{"name":"bar","index":0,"is_constant":true}},"children":[{"id":"root.0.0","symbols":[],"symbols_by_name":{},"is_block":true}]}]}}`),
		},
	}

	t.Run("successfully imports existing module", func(t *testing.T) {
		importer := NewImporter(ImporterOptions{
			SourceFS: testFS,
		})

		module, err := importer.Import(context.Background(), "foo/bar")
		require.NoError(t, err)
		require.NotNil(t, module)
		require.Equal(t, "foo/bar", module.Name().Value())
		require.NotNil(t, module.Code())

		code := module.Code()
		require.Equal(t, []string{"bar"}, code.GlobalNames())
	})

	t.Run("returns error for nonexistent module", func(t *testing.T) {
		importer := NewImporter(ImporterOptions{
			SourceFS: testFS,
		})

		module, err := importer.Import(context.Background(), "nonexistent")
		require.Error(t, err)
		require.Nil(t, module)
		require.Contains(t, err.Error(), "module \"nonexistent\" not found")
	})
}

func TestImporter_WithRealFixtures(t *testing.T) {
	// Create Importer using the fixtures directory in the current directory
	importer := NewImporter(ImporterOptions{
		SourceFS: os.DirFS("fixtures"),
	})

	// Import the bar module from fixtures
	module, err := importer.Import(context.Background(), "foo/bar")
	require.NoError(t, err)
	require.NotNil(t, module)
	require.Equal(t, "foo/bar", module.Name().Value())
	require.NotNil(t, module.Code())

	code := module.Code()
	require.Equal(t, []string{"bar"}, code.GlobalNames())
}
