package bundle

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/4rchr4y/bpm/constant"
	"github.com/4rchr4y/bpm/pkg/bundle/bundlefile"
	"github.com/4rchr4y/bpm/pkg/bundle/lockfile"
	"github.com/4rchr4y/bpm/pkg/bundle/regofile"
)

// TODO: move to ignorefile folder
type IgnoreFile struct {
	List map[string]struct{}
}

func NewIgnoreFile() *IgnoreFile {
	return &IgnoreFile{
		List: make(map[string]struct{}),
	}
}

func (f *IgnoreFile) Store(fileName string) {
	if fileName != "" {
		f.List[fileName] = struct{}{}
		return
	}
}

func (f *IgnoreFile) Some(path string) bool {
	if path == "" || len(f.List) == 0 {
		return false
	}

	dir := filepath.Dir(path)
	if dir == "." {
		return false
	}

	topLevelDir := strings.Split(dir, string(filepath.Separator))[0]
	_, found := f.List[topLevelDir]
	return found
}

type BundleRaw struct {
	BundleFile *bundlefile.File
	LockFile   *lockfile.File
	RegoFiles  map[string]*regofile.File
	OtherFiles map[string][]byte
}

type BundleOptFn func(*Bundle)

func WithIgnoreList(ignoreFile *IgnoreFile) BundleOptFn {
	return func(b *Bundle) {
		b.IgnoreFile = ignoreFile
	}
}

func WithVersion(v *VersionExpr) BundleOptFn {
	return func(b *Bundle) {
		b.Version = v
	}
}

func (br *BundleRaw) ToBundle(v *VersionExpr, ignoreFile *IgnoreFile) (*Bundle, error) {
	if br.BundleFile == nil {
		return nil, fmt.Errorf("can't find '%s' file", constant.BundleFileName)
	}

	b := &Bundle{
		Version:    v,
		BundleFile: br.BundleFile,
		LockFile:   br.LockFile,
		RegoFiles:  br.RegoFiles,
		IgnoreFile: ignoreFile,
		OtherFiles: br.OtherFiles,
	}

	if b.LockFile == nil {
		b.LockFile = lockfile.Init()
	}

	modules := make([]*lockfile.ModDecl, len(b.RegoFiles))
	var i int
	for filePath, f := range b.RegoFiles {
		modules[i] = &lockfile.ModDecl{
			Package: f.Package(),
			Source:  filePath,
			Sum:     f.Sum(),
		}

		i++
	}

	sort.Slice(modules, func(i, j int) bool { return modules[i].Package < modules[j].Package })
	b.LockFile.Modules = &lockfile.ModulesBlock{List: modules}
	b.LockFile.Sum = b.Sum()

	return b, nil
}

type Bundle struct {
	Version    *VersionExpr
	BundleFile *bundlefile.File
	LockFile   *lockfile.File
	IgnoreFile *IgnoreFile
	RegoFiles  map[string]*regofile.File
	OtherFiles map[string][]byte
}

func (b *Bundle) Name() string       { return b.BundleFile.Package.Name }
func (b *Bundle) Repository() string { return b.BundleFile.Package.Repository }

// Sum computes the overall checksum of the bundle. This method ensures the integrity
// of the bundle by taking into account various components that make up the bundle.
// It generates a SHA-256 hash that represents the combined checksum of the bundle's
// manifest file, all Rego policy files, and any other files included in the bundle.
//
// The process for generating the checksum is as follows:
//
//  1. Bundle manifest file checksum: The checksum of the bundle's main manifest file
//     is added to the hash first. This step ensures that the manifest file, which
//     acts as the primary descriptor of the bundle, has not been tampered with or
//     altered. The integrity of the manifest file is crucial for the bundle's consistency.
//
//  2. Rego policy files checksums: The checksums of all Rego policy files that constitute
//     the bundle are collected and added to the hash. This step guarantees that none of
//     the policy code files have been damaged or modified. Policy files are essential
//     components of the bundle, defining the policies that are enforced.
//
//  3. Other files checksums: Finally, the checksums of all other files included in the
//     bundle are collected and added to the hash. This ensures the integrity of the entire
//     bundle by verifying that no additional files have been altered. This includes any
//     supporting files that may not directly contain policy code but are nonetheless
//     part of the bundle's contents.
//
// The combined SHA-256 hash of these components is then encoded to a hexadecimal string
// and returned as the bundle's overall checksum. This checksum can be used to verify
// the bundle's integrity at a later time, ensuring that the bundle has not been
// altered since the checksum was generated.
func (b *Bundle) Sum() string {
	hasher := sha256.New()
	hasher.Write([]byte(b.BundleFile.Sum())) // Add checksum of the bundle file

	for _, k := range sortedMap(b.RegoFiles) {
		hasher.Write([]byte(b.RegoFiles[k].Sum())) // Add checksums of all Rego files
	}

	for _, k := range sortedMap(b.OtherFiles) {
		hasher.Write(b.OtherFiles[k]) // Add checksums of all other files
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func sortedMap[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
