package encode

import (
	"bytes"

	"github.com/4rchr4y/bpm/bundle/bundlefile"
	"github.com/4rchr4y/bpm/bundle/lockfile"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type BundleEncoder struct{}

func NewBundleEncoder() *BundleEncoder {
	return &BundleEncoder{}
}

func (e *BundleEncoder) DecodeBundleFile(content []byte) (*bundlefile.File, error) {
	bundlefile := new(bundlefile.File)

	if err := hclsimple.Decode("bundle.hcl", content, nil, bundlefile); err != nil {
		return nil, err
	}

	return bundlefile, nil
}

func (e *BundleEncoder) EncodeBundleFile(bundlefile *bundlefile.File) []byte {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(bundlefile, f.Body())

	return bytes.TrimSpace(f.Bytes())
}

func (e *BundleEncoder) EncodeLockFile(lockfile *lockfile.File) []byte {
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
