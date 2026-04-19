package modules

import (
	"maps"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"

	modmath "github.com/deepnoodle-ai/risor/v2/pkg/modules/math"
	modrand "github.com/deepnoodle-ai/risor/v2/pkg/modules/rand"
	modregexp "github.com/deepnoodle-ai/risor/v2/pkg/modules/regexp"

	modfilepath "github.com/foohq/ren/modules/filepath"
	modfs "github.com/foohq/ren/modules/fs"
	modos "github.com/foohq/ren/modules/os"
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
