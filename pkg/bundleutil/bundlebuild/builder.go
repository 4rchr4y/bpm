package bundlebuild

import "github.com/4rchr4y/bpm/pkg/bundle"

type BundleOptFn func(*bundle.Bundle)

func WithIgnoreList(ignoreFile *bundle.IgnoreFile) BundleOptFn {
	return func(b *bundle.Bundle) {
		b.IgnoreFile = ignoreFile
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
