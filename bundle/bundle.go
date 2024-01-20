package bundle

type Bundle struct {
	FileName       string
	BundleFile     *BundleFile
	BundleLockFile *BundleLockFile
	BpmWorkFile    *BpmWorkFile
	RegoFiles      map[string]*RawRegoFile
}

func (b *Bundle) UpdateLock() bool {
	if len(b.RegoFiles) < 1 {
		// no rego files, then nothing to update
		return false
	}

	if b.BundleLockFile == nil {
		b.BundleLockFile = &BundleLockFile{
			Version: 1,
			Modules: make([]*ModuleDef, len(b.RegoFiles)),
		}
	}

	var i uint
	for path, file := range b.RegoFiles {
		b.BundleLockFile.Modules[i] = &ModuleDef{
			Name:     file.Package(),
			Source:   path,
			Checksum: file.Sum(),
			Dependencies: func() []string {
				result := make([]string, len(file.Parsed.Imports))
				for i, _import := range file.Parsed.Imports {
					result[i] = _import.Path.String()
				}

				return result
			}(),
		}

		i++
	}

	return true
}

func (b *Bundle) Validation() error {
	panic("not implemented")
}
