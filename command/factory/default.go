package factory

func New(version string) *Factory {
	f := &Factory{
		Name:    "bpm",
		Version: version,
	}

	return f
}
