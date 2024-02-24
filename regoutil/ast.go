package regoutil

import (
	"github.com/4rchr4y/bpm/bundle/bundlefile"
	"github.com/open-policy-agent/opa/ast"
)

func PrepareDocumentParser(schema *bundlefile.Schema) error {
	// adding the name of all required bundles to the list
	// of allowed keywords in the rego environment
	if schema.Require != nil {
		for _, r := range schema.Require.List {
			packageVar := ast.VarTerm(r.Name)
			packageRef := ast.Ref{packageVar}

			ast.RootDocumentNames.Add(packageVar)
			ast.ReservedVars.Add(packageVar.Value.(ast.Var))
			ast.RootDocumentRefs.Add(ast.NewTerm(packageRef))
		}
	}

	return nil
}
