package bundle

const (
	BundleFileName = "bundle.toml"
	LockFileName   = "bundle.lock"
	WorkFileName   = "bpm.work"
	IgnoreFileName = ".bpmignore"
)

type ManifestFile interface {
	Name() string

	manifestFile()
}
