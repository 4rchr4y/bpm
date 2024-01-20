package parser

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

type BundleParser struct {
	decoder bpTOMLDecoder
}

func NewBundleParser(decoder bpTOMLDecoder) *BundleParser {
	return &BundleParser{
		decoder: decoder,
	}
}

type ParseInput struct {
	FileName string
	Files    map[string][]byte
}

// TODO: check that file is not empty; isEmpty(content)
func (bp *BundleParser) Parse(input *ParseInput) (*bundle.Bundle, error) {
	b := &bundle.Bundle{
		FileName:  filepath.Clean(input.FileName),
		RegoFiles: make(map[string]*bundle.RawRegoFile),
	}

	for filePath, content := range input.Files {
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

		case isBPMWorkFile(filePath):
			bpmwork, err := bp.parseBPMWorkFile(content)
			if err != nil {
				return nil, err
			}

			b.BpmWorkFile = bpmwork
		}
	}

	return b, nil
}

func (bp *BundleParser) parseRegoFile(fileContent []byte, filePath string) (*ast.Module, error) {
	parsed, err := ast.ParseModule(filePath, string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing file contents: %v", err)
	}

	return parsed, nil
}

func (bp *BundleParser) parseBPMWorkFile(fileContent []byte) (*bundle.BpmWorkFile, error) {
	var bpmwork bundle.BpmWorkFile
	if err := bp.decoder.Decode(string(fileContent), &bpmwork); err != nil {
		return nil, fmt.Errorf("error parsing bpm.work content: %v", err)
	}

	return &bpmwork, nil
}

func (bp *BundleParser) parseBPMLockFile(fileContent []byte) (*bundle.BundleLockFile, error) {
	var bundlelock bundle.BundleLockFile
	if err := bp.decoder.Decode(string(fileContent), &bundlelock); err != nil {
		return nil, fmt.Errorf("error parsing bundle.lock content: %v", err)
	}

	return &bundlelock, nil
}

func (bp *BundleParser) parseBPMFile(fileContent []byte) (*bundle.BundleFile, error) {
	var bundlefile bundle.BundleFile
	if err := bp.decoder.Decode(string(fileContent), &bundlefile); err != nil {
		return nil, fmt.Errorf("error parsing bundle.toml content: %v", err)
	}

	return &bundlefile, nil
}

func isRegoFile(filePath string) bool    { return filepath.Ext(filePath) == constant.RegoFileExt }
func isBPMFile(filePath string) bool     { return filePath == constant.BundleFileName }
func isBPMLockFile(filePath string) bool { return filePath == constant.LockFileName }
func isBPMWorkFile(filePath string) bool { return filePath == constant.WorkFileName }
