package fileifier

import (
	"fmt"
	"path/filepath"

	"github.com/4rchr4y/bpm/bundle"
	"github.com/4rchr4y/bpm/constant"

	"github.com/open-policy-agent/opa/ast"
)

type bpTOMLDecoder interface {
	Decode(data string, v interface{}) error
}

type Fileifier struct {
	decoder bpTOMLDecoder
}

func NewFileifier(decoder bpTOMLDecoder) *Fileifier {
	return &Fileifier{
		decoder: decoder,
	}
}

// TODO: check that file is not empty; isEmpty(content)
func (bp *Fileifier) Fileify(files map[string][]byte) (*bundle.Bundle, error) {
	b := &bundle.Bundle{
		RegoFiles: make(map[string]*bundle.RawRegoFile),
	}

	for filePath, content := range files {
		switch {
		case isRegoFile(filePath):
			parsed, err := bp.parseRegoFile(content, filePath)
			if err != nil {
				return nil, err
			}

			b.RegoFiles[filePath] = &bundle.RawRegoFile{
				Path:   filePath,
				Parsed: parsed,
			}

		case isBPMFile(filePath):
			bundlefile, err := bp.parseBPMFile(content)
			if err != nil {
				return nil, err
			}

			b.BundleFile = bundlefile

		case isBPMLockFile(filePath):
			bundlelock, err := bp.parseBPMLockFile(content)
			if err != nil {
				return nil, err
			}

			b.BundleLockFile = bundlelock
		}
	}

	return b, nil
}

func (bp *Fileifier) parseRegoFile(fileContent []byte, filePath string) (*ast.Module, error) {
	parsed, err := ast.ParseModule(filePath, string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing file contents: %v", err)
	}

	return parsed, nil
}

func (bp *Fileifier) parseBPMLockFile(fileContent []byte) (*bundle.BundleLockFile, error) {
	var bundlelock bundle.BundleLockFile
	if err := bp.decoder.Decode(string(fileContent), &bundlelock); err != nil {
		return nil, fmt.Errorf("error parsing bundle.lock content: %v", err)
	}

	return &bundlelock, nil
}

func (bp *Fileifier) parseBPMFile(fileContent []byte) (*bundle.BundleFile, error) {
	var bundlefile bundle.BundleFile
	if err := bp.decoder.Decode(string(fileContent), &bundlefile); err != nil {
		return nil, fmt.Errorf("error parsing bundle.toml content: %v", err)
	}

	return &bundlefile, nil
}

func isRegoFile(filePath string) bool    { return filepath.Ext(filePath) == constant.RegoFileExt }
func isBPMFile(filePath string) bool     { return filePath == constant.BundleFileName }
func isBPMLockFile(filePath string) bool { return filePath == constant.LockFileName }
