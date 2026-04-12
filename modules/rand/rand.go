package rand

import (
	modrand "github.com/deepnoodle-ai/risor/v2/pkg/modules/rand"
	"github.com/deepnoodle-ai/risor/v2/pkg/object"
)

func Module() *object.Module {
	return modrand.Module()
}
