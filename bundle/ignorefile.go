package bundle

import (
	"path/filepath"
	"strings"

	"github.com/4rchr4y/bpm/constant"
)

type IgnoreFile struct {
	List     map[string]struct{}
	IsSorted bool
}

func NewIgnoreFile(size ...int) *IgnoreFile {
	var initialSize int
	if len(size) > 0 {
		initialSize = size[0]
	}

	return &IgnoreFile{
		List:     make(map[string]struct{}, initialSize),
		IsSorted: false,
	}
}

func (*IgnoreFile) Filename() string { return constant.IgnoreFileName }

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
