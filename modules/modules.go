package modules

import (
	"github.com/risor-io/risor/object"

	modbase64 "github.com/foohq/ren/modules/base64"
	modbuiltins "github.com/foohq/ren/modules/builtins"
	modbytes "github.com/foohq/ren/modules/bytes"
	modcli "github.com/foohq/ren/modules/cli"
	moderrors "github.com/foohq/ren/modules/errors"
	modexec "github.com/foohq/ren/modules/exec"
	modfilepath "github.com/foohq/ren/modules/filepath"
	modfmt "github.com/foohq/ren/modules/fmt"
	modhttp "github.com/foohq/ren/modules/http"
	modjson "github.com/foohq/ren/modules/json"
	modmath "github.com/foohq/ren/modules/math"
	modnet "github.com/foohq/ren/modules/net"
	modos "github.com/foohq/ren/modules/os"
	modrand "github.com/foohq/ren/modules/rand"
	modregexp "github.com/foohq/ren/modules/regexp"
	modshlex "github.com/foohq/ren/modules/shlex"
	modstrconv "github.com/foohq/ren/modules/strconv"
	modstrings "github.com/foohq/ren/modules/strings"
	modtime "github.com/foohq/ren/modules/time"
	modurlpath "github.com/foohq/ren/modules/urlpath"
)

var (
	modules = map[string]*object.Module{
		"base64":   modbase64.Module(),
		"builtins": modbuiltins.Module(),
		"bytes":    modbytes.Module(),
		"cli":      modcli.Module(),
		"errors":   moderrors.Module(),
		"exec":     modexec.Module(),
		"filepath": modfilepath.Module(),
		"fmt":      modfmt.Module(),
		"http":     modhttp.Module(),
		"json":     modjson.Module(),
		"math":     modmath.Module(),
		"net":      modnet.Module(),
		"os":       modos.Module(),
		"rand":     modrand.Module(),
		"regexp":   modregexp.Module(),
		"shlex":    modshlex.Module(),
		"strconv":  modstrconv.Module(),
		"strings":  modstrings.Module(),
		"time":     modtime.Module(),
		"urlpath":  modurlpath.Module(),
	}
	builtins = map[string]map[string]object.Object{
		"base64":   modbase64.Builtins(),
		"builtins": modbuiltins.Builtins(),
		"bytes":    modbytes.Builtins(),
		"cli":      modcli.Builtins(),
		"errors":   moderrors.Builtins(),
		"exec":     modexec.Builtins(),
		"filepath": modfilepath.Builtins(),
		"fmt":      modfmt.Builtins(),
		"http":     modhttp.Builtins(),
		"json":     modjson.Builtins(),
		"math":     modmath.Builtins(),
		"net":      modnet.Builtins(),
		"os":       modos.Builtins(),
		"rand":     modrand.Builtins(),
		"regexp":   modregexp.Builtins(),
		"shlex":    modshlex.Builtins(),
		"strconv":  modstrconv.Builtins(),
		"strings":  modstrings.Builtins(),
		"time":     modtime.Builtins(),
		"urlpath":  modurlpath.Builtins(),
	}
)

func Modules() []string {
	result := make([]string, 0, len(modules))
	for name := range modules {
		result = append(result, name)
	}
	return result
}

func Globals() map[string]any {
	result := make(map[string]any, len(modules)+len(builtins))
	for name, module := range modules {
		if module == nil {
			continue
		}
		result[name] = module
	}
	for _, builtin := range builtins {
		if builtin == nil {
			continue
		}
		for name, fn := range builtin {
			result[name] = fn
		}
	}
	return result
}

func GlobalNames() []string {
	result := make([]string, len(modules)+len(builtins))
	for name, module := range modules {
		if module == nil {
			continue
		}
		result = append(result, name)
	}
	for _, builtin := range builtins {
		if builtin == nil {
			continue
		}
		for name := range builtin {
			result = append(result, name)
		}
	}
	return result
}

// BuildTag returns build tag for a module name. The function does not check existence of the module.
func BuildTag(name string) string {
	return "module_" + name
}
