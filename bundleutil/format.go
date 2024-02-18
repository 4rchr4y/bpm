package bundleutil

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/4rchr4y/bpm/bundle/lockfile"
)

// regex to find and replace extra newlines after {
var bracketNormalizerRegex *regexp.Regexp = regexp.MustCompile(`(?m){\s*\n+`)

// regex to remove quotes around keywords
var lockfileKeywordNormalizerRegex *regexp.Regexp

func init() {
	lockfileKeywordNormalizerPattern := fmt.Sprintf(`"(%s)"`, joinKeywordList("|", lockfile.Keywords[:]))
	lockfileKeywordNormalizerRegex = regexp.MustCompile(lockfileKeywordNormalizerPattern)
}

// FormatBundleFile designed for post-processing and bundlefile formatting
func FormatBundleFile(content []byte) []byte {
	content = bytes.TrimSpace(content)
	content = bracketNormalizerRegex.ReplaceAll(content, []byte("{\n"))
	content = bytes.Replace(content, []byte("{\n}"), []byte("{}"), -1)

	return content
}

// FormatBundleFile designed for post-processing and lockfile formatting
func FormatLockFile(content []byte) []byte {
	content = bytes.TrimSpace(content)
	content = bracketNormalizerRegex.ReplaceAll(content, []byte("{\n"))
	content = lockfileKeywordNormalizerRegex.ReplaceAll(content, []byte("$1"))
	content = bytes.Replace(content, []byte("{\n}"), []byte("{}"), -1)

	return content
}

// FormatBundleFile designed to create a single place for forming a string
// representing the repository and the bundle version
func FormatSourceVersion(source, version string) string {
	return source + "@" + version
}

func joinKeywordList(separator string, keywords []string) string {
	var builder strings.Builder

	for i, keyword := range keywords {
		if i > 0 {
			builder.WriteString(separator)
		}
		builder.WriteString(keyword)
	}

	return builder.String()
}
