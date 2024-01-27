package fileify

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundle"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/4rchr4y/bpm/pkg/bundle/regofile"
	"github.com/open-policy-agent/opa/ast"
)

type bfEncoder interface {
	DecodeBundleFile(content []byte) (*bundlefile.File, error)
	DecodeLockFile(content []byte) (*lockfile.File, error)
}

type Fileifier struct {
	encoder bfEncoder
}

func NewFileifier(decoder bfEncoder) *Fileifier {
	return &Fileifier{
		encoder: decoder,
	}
}

func (bp *Fileifier) Fileify(files map[string][]byte) (*bundle.Bundle, error) {
	b := &bundle.Bundle{
		RegoFiles:  make(map[string]*regofile.File),
		OtherFiles: make(map[string][]byte),
	}

	ignoreFileContent, exist := files[constant.IgnoreFileName]
	if exist && !isEmpty(ignoreFileContent) {
		ignoreList, err := bp.parseIgnoreFile(ignoreFileContent)
		if err != nil {
			return nil, err
		}

		b.IgnoreFiles = ignoreList
	}

	for filePath, content := range files {
		switch {
		case shouldIgnore(b.IgnoreFiles, filePath):
			continue

		case isRegoFile(filePath):
			parsed, err := bp.parseRegoFile(content, filePath)
			if err != nil {
				return nil, err
			}

			b.RegoFiles[filePath] = &regofile.File{
				Path:   filePath,
				Parsed: parsed,
			}

		case isBPMFile(filePath):
			bundlefile, err := bp.encoder.DecodeBundleFile(content)
			if err != nil {
				return nil, fmt.Errorf("error occurred while decoding %s content: %v", constant.BundleFileName, err)
			}

			b.BundleFile = bundlefile

		case isBPMLockFile(filePath):
			lockfile, err := bp.encoder.DecodeLockFile(content)
			if err != nil {
				return nil, fmt.Errorf("error occurred while decoding %s content: %v", constant.BundleFileName, err)
			}

			b.BundleLockFile = lockfile

		default:
			if isEmpty(ignoreFileContent) {
				continue
			}

			b.OtherFiles[filePath] = content
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

func (bp *Fileifier) parseBPMLockFile(fileContent []byte) (*lockfile.File, error) {
	var bundlelock lockfile.File

	return &bundlelock, nil
}

func (bp *Fileifier) parseIgnoreFile(fileContent []byte) (map[string]struct{}, error) {
	result := make(map[string]struct{})

	scanner := bufio.NewScanner(bytes.NewReader(fileContent))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		result[line] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading '%s' input: %v", constant.IgnoreFileName, err)
	}

	return result, nil
}

func shouldIgnore(ignoreList map[string]struct{}, path string) bool {
	if path == "" || len(ignoreList) == 0 {
		return false
	}

	dir := filepath.Dir(path)
	if dir == "." {
		return false
	}

	topLevelDir := strings.Split(dir, string(filepath.Separator))[0]
	_, found := ignoreList[topLevelDir]
	return found
}

func isRegoFile(filePath string) bool    { return filepath.Ext(filePath) == constant.RegoFileExt }
func isBPMFile(filePath string) bool     { return filePath == constant.BundleFileName }
func isBPMLockFile(filePath string) bool { return filePath == constant.LockFileName }
func isEmpty(content []byte) bool        { return len(content) < 1 }
