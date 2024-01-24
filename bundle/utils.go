package bundle

import "github.com/4rchr4y/bpm/bundle/lockfile"

func ValidateBundle(b *Bundle) error {
	return nil
}

func UpdateLockFile(b *Bundle) bool {
	if len(b.RegoFiles) < 1 {
		// no rego files, then nothing to update
		return false
	}

	if b.BundleLockFile == nil {
		b.BundleLockFile = &lockfile.File{
			Version: 1,
			Modules: make([]*lockfile.ModuleDef, len(b.RegoFiles)),
		}
	}

	var i uint
	for path, file := range b.RegoFiles {
		b.BundleLockFile.Modules[i] = &lockfile.ModuleDef{
			Name:     file.Package(),
			Source:   path,
			Checksum: file.Sum(),
			Require: func() []*lockfile.Requirement {
				result := make([]*lockfile.Requirement, len(file.Parsed.Imports))

				for i, _import := range file.Parsed.Imports {
					result[i] = &lockfile.Requirement{
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
