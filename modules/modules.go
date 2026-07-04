package modules

import (
	"maps"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"

	moddll "github.com/foohq/ren/modules/dll"
	modfilepath "github.com/foohq/ren/modules/filepath"
	modfs "github.com/foohq/ren/modules/fs"
	modos "github.com/foohq/ren/modules/os"
	modshlex "github.com/foohq/ren/modules/shlex"
)

var modules = map[string]*object.Module{
	//"cli":      modcli.Module(),
	"dll": moddll.Module(),
	//"exec":     modexec.Module(),
	"filepath": modfilepath.Module(),
	"fs":       modfs.Module(),
	//"http":     modhttp.Module(),
	//"net":      modnet.Module(),
	"os":    modos.Module(),
	"shlex": modshlex.Module(),
}

func Modules() map[string]*object.Module {
	result := make(map[string]*object.Module, len(modules))
	maps.Copy(result, modules)
	return result
}
