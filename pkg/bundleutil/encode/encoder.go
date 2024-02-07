package encode

import (
	"bytes"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) DecodeBundleFile(content []byte) (*bundlefile.File, error) {
	file := new(bundlefile.File)
	if err := hclsimple.Decode(constant.BundleFileName, content, nil, file); err != nil {
		return nil, err
	}

	return file, nil
}

func (e *Encoder) DecodeLockFile(content []byte) (*lockfile.File, error) {
	file := new(lockfile.File)
	if err := hclsimple.Decode(constant.LockFileName, content, nil, file); err != nil {
		return nil, err
	}

	return file, nil
}

func (e *Encoder) EncodeBundleFile(bundlefile *bundlefile.File) []byte {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(bundlefile, f.Body())

	result := bytes.TrimSpace(f.Bytes())
	result = bytes.Replace(result, []byte("{\n\n"), []byte("{\n"), -1)

	return result
}

const lockfileComment = "// This file has been auto-generated by `bpm`.\n// It is not meant to be edited manually."

func (e *Encoder) EncodeLockFile(lockfile *lockfile.File) []byte {
	tempFile := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(lockfile, tempFile.Body())

	f := hclwrite.NewEmptyFile()

	f.Body().AppendUnstructuredTokens([]*hclwrite.Token{
		{Type: hclsyntax.TokenComment, Bytes: []byte(lockfileComment)},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
		{Type: hclsyntax.TokenOBrace, Bytes: tempFile.Bytes()},
	})

	result := bytes.TrimSpace(f.Bytes())
	result = bytes.Replace(result, []byte("\"direct\""), []byte("direct"), -1)
	result = bytes.Replace(result, []byte("\"indirect\""), []byte("indirect"), -1)
	result = bytes.Replace(result, []byte("{\n\n"), []byte("{\n"), -1)

	return result
}