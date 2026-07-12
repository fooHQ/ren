package builtins_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/foohq/ren/builtins"
)

// TestDocsMatchBuiltins guards that Docs covers exactly the set of registered
// global builtins: every builtin is documented, and no documentation refers to
// a builtin that is not registered.
func TestDocsMatchBuiltins(t *testing.T) {
	documented := make([]string, 0, len(builtins.Docs()))
	for _, spec := range builtins.Docs() {
		documented = append(documented, spec.Name)
	}

	registered := make([]string, 0, len(builtins.Builtins()))
	for name := range builtins.Builtins() {
		registered = append(registered, name)
	}

	require.ElementsMatch(t, registered, documented)
}
