package linker

import (
	"context"
	"fmt"
	"strings"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/open-policy-agent/opa/ast"
)

type linkerInspector interface {
	Inspect(b *bundle.Bundle) error
}

type linkerManifester interface {
	SyncLockfile(ctx context.Context, parent *bundle.Bundle) error
}

type linkerFetcher interface {
	PlainFetch(ctx context.Context, source string, version *bundle.VersionSpec) (*bundle.Bundle, error)
}

type Linker struct {
	Fetcher    linkerFetcher
	Manifester linkerManifester
	Inspector  linkerInspector
}

func (l *Linker) Link(ctx context.Context, b *bundle.Bundle) (map[string]*ast.Module, error) {
	if err := l.Manifester.SyncLockfile(ctx, b); err != nil {
		return nil, err
	}

	if err := l.Inspector.Inspect(b); err != nil {
		return nil, err
	}

	result := make(map[string]*ast.Module)

	requireCache := make(map[string]struct{}, len(b.BundleFile.Require.List))
	for _, r := range b.BundleFile.Require.List {
		requireCache[r.Name] = struct{}{}
	}

	// save all the modules of the head bundle
	for filePath, f := range b.RegoFiles {
		result[filePath] = f.Parsed
		fmt.Println("-----")
		for _, _import := range f.Parsed.Imports {
			importPathStr := _import.Path.String()
			sourcePkg := strings.Index(importPathStr, ".")

			if _, exists := requireCache[importPathStr[:sourcePkg]]; !exists {
				continue
			}

			value := ast.StringTerm("data." + importPathStr)
			value.SetLocation(_import.Path.Location)
			_import.Path = value

			// value := ast.StringTerm("data." + importPathStr)

			// _import.Path = &ast.Term{
			// 	Value:    value.Value,
			// 	Location: _import.Location,
			// }
		}

		fmt.Println(f.Parsed.String())
	}

	// iter over all required bundles
	for _, item := range b.LockFile.Require.List {
		v, err := bundle.ParseVersionExpr(item.Version)
		if err != nil {
			return nil, err
		}

		itemBundle, err := l.Fetcher.PlainFetch(ctx, item.Repository, v)
		if err != nil {
			return nil, err
		}

		// save all the modules of the required bundle
		for filePath, f := range itemBundle.RegoFiles {
			result[filePath] = f.Parsed
		}
	}

	return result, nil
}
