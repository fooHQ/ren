package shlex_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/foohq/ren/modules/shlex"
)

// TestDocsResolve guards that every name documented in docs.go is actually
// registered by the module, so the documentation cannot reference functions
// that do not exist.
func TestDocsResolve(t *testing.T) {
	m := shlex.Module()
	seen := make(map[string]bool)
	for _, spec := range shlex.Docs() {
		require.NotEmpty(t, spec.Name)
		require.Falsef(t, seen[spec.Name], "duplicate documentation for %q", spec.Name)
		seen[spec.Name] = true

		_, ok := m.GetAttr(spec.Name)
		require.Truef(t, ok, "documented name %q is not registered by the module", spec.Name)
	}
}
