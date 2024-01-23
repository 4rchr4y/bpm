package bundle

func ValidateBundle(b *Bundle) error {
	return nil
}

func UpdateLockFile(b *Bundle) bool {
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
			Require: func() []*Requirement {
				result := make([]*Requirement, len(file.Parsed.Imports))

				for i, _import := range file.Parsed.Imports {
					result[i] = &Requirement{
						Package: _import.Path.String(),
					}
				}

				return result
			}(),
		}

		i++
	}

	return true
}
