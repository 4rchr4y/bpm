package bundleutil

import "github.com/4rchr4y/bpm/pkg/bundle"

type BundleOptFn func(*bundle.Bundle)

func WithIgnoreList(ignoreList map[string]struct{}) BundleOptFn {
	return func(b *bundle.Bundle) {
		b.IgnoreFiles = ignoreList
	}
}

func WithVersion(v *bundle.VersionExpr) BundleOptFn {
	return func(b *bundle.Bundle) {
		b.Version = v
	}
}

func WithRepository(url string) BundleOptFn {
	return func(b *bundle.Bundle) {
		b.BundleFile.Package.Repository = url
	}
}
