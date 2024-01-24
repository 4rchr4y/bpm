package bfencoder

import (
	"github.com/4rchr4y/bpm/bundle"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) EncodeBundleFile(bundlefile *bundle.BundleFile) []byte {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	if bundlefile.Package != nil {
		packageBlock := rootBody.AppendNewBlock("package", nil)
		packageBody := packageBlock.Body()
		packageBody.SetAttributeValue("name", cty.StringVal(bundlefile.Package.Name))
		packageBody.SetAttributeValue("author", cty.ListVal(transformStringList(bundlefile.Package.Author)))
		packageBody.SetAttributeValue("repository", cty.StringVal(bundlefile.Package.Repository))
		packageBody.SetAttributeValue("description", cty.StringVal(bundlefile.Package.Description))
	}

	for reqKey, reqVal := range bundlefile.Require {
		requireBlock := rootBody.AppendNewBlock("require", []string{reqKey})
		requireBody := requireBlock.Body()
		requireBody.SetAttributeValue("name", cty.StringVal(reqVal.Name))
		requireBody.SetAttributeValue("version", cty.StringVal(reqVal.Version))
	}

	return f.Bytes()
}

func transformStringList(data []string) []cty.Value {
	result := make([]cty.Value, len(data))
	for i := range data {
		result[i] = cty.StringVal(data[i])
	}

	return result
}
