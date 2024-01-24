package bfencoder

import (
	"bytes"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) DecodeBundleFile(content []byte) (*bundle.BundleFile, error) {
	bundlefile := new(bundle.BundleFile)

	if err := hclsimple.Decode("bundle.hcl", content, nil, bundlefile); err != nil {
		return nil, err
	}

	return bundlefile, nil
}

func (e *Encoder) EncodeBundleFile(bundlefile *bundle.BundleFile) []byte {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(bundlefile, f.Body())

	return bytes.TrimSpace(f.Bytes())
}

func (e *Encoder) EncodeLockFile(lockfile *bundle.BundleLockFile) []byte {
	f := hclwrite.NewEmptyFile()

	return f.Bytes()
}

func transformStringList(data []string) []cty.Value {
	result := make([]cty.Value, len(data))
	for i := range data {
		result[i] = cty.StringVal(data[i])
	}

	return result
}
