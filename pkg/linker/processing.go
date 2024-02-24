package linker

import (
	"strings"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/bundle/regofile"
	"github.com/open-policy-agent/opa/ast"
)

type ModuleProcessFn func(*ast.Module)

func ProcessModule(m *ast.Module, processes ...ModuleProcessFn) *ast.Module {
	for _, fn := range processes {
		fn(m)
	}

	return m
}

// WithImportProcessing will traverse all imports to identify all bundle imports.
// This is done to allow users to specify imports without
// the `.date` prefix in each import.
//
// Import without this feature:
// import data.example.folder.file
//
// With this feature:
// import example.folder.file
//
// When traversing imports, bundle imports are identified
// and the `.data` prefix is automatically added to such imports.
func WithImportProcessing(requireList map[string]struct{}) ModuleProcessFn {
	return func(m *ast.Module) {
		for _, importItem := range m.Imports {
			importPathStr := importItem.Path.String()
			sourcePkg := strings.Index(importPathStr, ".")
			if _, exists := requireList[importPathStr[:sourcePkg]]; !exists {
				continue
			}

			value := ast.VarTerm(regofile.ImportPathPrefix + importPathStr)
			value.Location = importItem.Path.Location

			importItem.Path = value
		}
	}
}

func getRequireList(b *bundle.Bundle) map[string]struct{} {
	requireCache := make(map[string]struct{}, len(b.BundleFile.Require.List))
	for _, r := range b.BundleFile.Require.List {
		requireCache[r.Name] = struct{}{}
	}

	return requireCache
}
