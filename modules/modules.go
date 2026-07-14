// Package modules is the registry of built-in modules that Ren scripts can
// import via the "builtin://" scheme (for example, import("builtin://os")).
package modules

import (
	"maps"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"

	moddll "github.com/foohq/ren/modules/dll"
	modfilepath "github.com/foohq/ren/modules/filepath"
	modfs "github.com/foohq/ren/modules/fs"
	modos "github.com/foohq/ren/modules/os"
)

var modules = map[string]*object.Module{
	//"cli":      modcli.Module(),
	"dll": moddll.Module(),
	//"exec":     modexec.Module(),
	"filepath": modfilepath.Module(),
	"fs":       modfs.Module(),
	//"http":     modhttp.Module(),
	//"net":      modnet.Module(),
	"os": modos.Module(),
}

// Modules returns a copy of the registry, mapping each built-in module's name
// to its object.Module.
func Modules() map[string]*object.Module {
	result := make(map[string]*object.Module, len(modules))
	maps.Copy(result, modules)
	return result
}

// ModuleDocs bundles a module's name, its module-level documentation, and the
// documentation for each name it exposes.
type ModuleDocs struct {
	Name  string
	Doc   string
	Funcs []object.FuncSpec
}

// Docs returns documentation for every built-in module, in a stable order
// suitable for rendering a reference.
func Docs() []ModuleDocs {
	return []ModuleDocs{
		{Name: "os", Doc: modos.ModuleDoc(), Funcs: modos.Docs()},
		{Name: "fs", Doc: modfs.ModuleDoc(), Funcs: modfs.Docs()},
		{Name: "filepath", Doc: modfilepath.ModuleDoc(), Funcs: modfilepath.Docs()},
		{Name: "dll", Doc: moddll.ModuleDoc(), Funcs: moddll.Docs()},
	}
}
