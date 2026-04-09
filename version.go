package ren

import (
	_ "embed"
	"strings"
)

//go:embed version.txt
var version string

// Version returns the current version of Ren.
func Version() string {
	return strings.TrimSpace(version)
}
