package bundleutil

import (
	"github.com/4rchr4y/bpm/pkg/bundle"
)

func FormatVersion(source string, version string) string {
	return source + "@" + version
}

func FormatVersionFromBundle(b *bundle.Bundle) string {
	return FormatVersion(b.Repository(), b.Version.String())
}
