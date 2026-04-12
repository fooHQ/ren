package modules

import (
	"maps"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"

	modfilepath "github.com/foohq/ren/modules/filepath"
	modfs "github.com/foohq/ren/modules/fs"
	modmath "github.com/foohq/ren/modules/math"
	modos "github.com/foohq/ren/modules/os"
	modrand "github.com/foohq/ren/modules/rand"
	modregexp "github.com/foohq/ren/modules/regexp"
	modshlex "github.com/foohq/ren/modules/shlex"
)

var modules = map[string]*object.Module{
	//"cli":      modcli.Module(),
	//"exec":     modexec.Module(),
	"filepath": modfilepath.Module(),
	"fs":       modfs.Module(),
	//"http":     modhttp.Module(),
	"math": modmath.Module(),
	//"net":      modnet.Module(),
	"os":     modos.Module(),
	"rand":   modrand.Module(),
	"regexp": modregexp.Module(),
	"shlex":  modshlex.Module(),
}

func Modules() map[string]*object.Module {
	result := make(map[string]*object.Module, len(modules))
	maps.Copy(result, modules)
	return result
}
