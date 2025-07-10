package importer

import (
	"context"
	"fmt"
	"io"
	"io/fs"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/object"
)

// FSImporterOptions configure an Importer that can read from a filesystem
// implementing the `fs.FS` interface.
type FSImporterOptions struct {
	// Global names that should be available when the module is compiled.
	GlobalNames []string

	// The filesystem to search for Risor modules.
	SourceFS fs.FS
}

// FSImporter is an Importer that can read Risor code modules from a filesystem
// implementing the `fs.FS` interface.
type FSImporter struct {
	globalNames []string
	sourceFS    fs.FS
}

// NewFSImporter returns an Importer that can read Risor code modules from a
// filesystem implementing the `fs.FS` interface.
func NewFSImporter(opts FSImporterOptions) *FSImporter {
	return &FSImporter{
		globalNames: opts.GlobalNames,
		sourceFS:    opts.SourceFS,
	}
}

// Import a module by name.
func (i *FSImporter) Import(_ context.Context, name string) (*object.Module, error) {
	source, found := i.readFileWithExtension(name)
	if !found {
		return nil, fmt.Errorf("import error: module %q not found", name)
	}

	code, err := compiler.UnmarshalCode(source)
	if err != nil {
		return nil, err
	}

	return object.NewModule(name, code), nil
}

func (i *FSImporter) readFileWithExtension(name string) ([]byte, bool) {
	f, err := i.sourceFS.Open(name + ".json")
	if err != nil {
		return nil, false
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, false
	}

	return b, true
}
