package regexp

import (
	modregexp "github.com/deepnoodle-ai/risor/v2/pkg/modules/regexp"
	"github.com/deepnoodle-ai/risor/v2/pkg/object"
)

func Module() *object.Module {
	return modregexp.Module()
}
