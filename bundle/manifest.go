package bundle

type ManifestFile interface {
	FileName() string

	manifestFile()
}
