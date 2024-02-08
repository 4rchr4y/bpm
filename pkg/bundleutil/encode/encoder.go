package encode

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/core"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/4rchr4y/bpm/pkg/bundle/regofile"
	"github.com/4rchr4y/bpm/pkg/bundleutil"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/open-policy-agent/opa/ast"
)

type Encoder struct {
	IO core.IO
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) DecodeIgnoreFile(content []byte) (*bundle.IgnoreFile, error) {
	ignoreFile := bundle.NewIgnoreFile()
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		ignoreFile.Store(strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading '%s' input: %v", constant.IgnoreFileName, err)
	}

	return ignoreFile, nil
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

func (e *Encoder) Fileify(files map[string][]byte, options ...bundleutil.BundleOptFn) (*bundle.Bundle, error) {
	b := &bundle.Bundle{
		RegoFiles:  make(map[string]*regofile.File),
		OtherFiles: make(map[string][]byte),
	}

	for i := range options {
		options[i](b)
	}

	modules := &lockfile.ModulesDecl{
		List: make([]*lockfile.ModDecl, 0),
	}

	for filePath, content := range files {
		switch {
		case isRegoFile(filePath):
			parsed, err := ast.ParseModule(filePath, string(content))
			if err != nil {
				return nil, fmt.Errorf("error parsing file contents: %v", err)
			}

			file := &regofile.File{Path: filePath, Parsed: parsed}
			b.RegoFiles[filePath] = file

			modules.List = append(modules.List, &lockfile.ModDecl{
				Package: file.Package(),
				Source:  filePath,
				Sum:     file.Sum(),
			})

		case isBPMFile(filePath):
			bundlefile, err := e.DecodeBundleFile(content)
			if err != nil {
				return nil, fmt.Errorf("error occurred while decoding %s content: %v", constant.BundleFileName, err)
			}

			b.BundleFile = bundlefile

		case isBPMLockFile(filePath):
			lockfile, err := e.DecodeLockFile(content)
			if err != nil {
				return nil, fmt.Errorf("error occurred while decoding %s content: %v", constant.BundleFileName, err)
			}

			b.LockFile = lockfile

		default:
			b.OtherFiles[filePath] = content
		}
	}

	if b.LockFile == nil {
		e.IO.PrintfWarn("file '%s' in bundle '%s' was not found", constant.LockFileName, b.Name())

		if err := initLockFile(b); err != nil {
			return nil, err
		}
	}

	if len(modules.List) > 0 {
		b.LockFile.Modules = modules.Sort()
	}

	return b, nil
}

func isRegoFile(filePath string) bool    { return filepath.Ext(filePath) == constant.RegoFileExt }
func isBPMFile(filePath string) bool     { return filePath == constant.BundleFileName }
func isBPMLockFile(filePath string) bool { return filePath == constant.LockFileName }
func isEmpty(content []byte) bool        { return len(content) < 1 }

func initLockFile(b *bundle.Bundle) error {
	if b.BundleFile == nil {
		return fmt.Errorf("can't find '%s' file", constant.BundleFileName)
	}

	b.LockFile = &lockfile.File{
		// TODO: set 'edition' from global app context
		Edition: "2024",
		Sum:     b.Sum(),
	}

	return nil
}
