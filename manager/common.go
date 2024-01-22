package manager

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func validateRepoURL(repoURL string) error {
	trimmedURL := strings.TrimSpace(repoURL)
	if trimmedURL == "" {
		return fmt.Errorf("passed URL is empty")
	}

	if strings.HasPrefix(trimmedURL, "https://") || strings.HasSuffix(trimmedURL, ".git") {
		return fmt.Errorf("URL must not start with 'https://' and end with '.git'")
	}

	return nil
}

func validatePath(path string) error {
	const (
		minPathLength = 1
		maxPathLength = 255
	)

	if len(path) < minPathLength || len(path) > maxPathLength {
		return fmt.Errorf("length is not within the valid")
	}

	if !utf8.ValidString(path) {
		return fmt.Errorf("contains invalid UTF-8 characters")
	}

	invalidPatterns := []string{"..", "://", "\x00"}
	for _, pattern := range invalidPatterns {
		if strings.Contains(path, pattern) {
			return fmt.Errorf("contains invalid pattern '%s'", pattern)
		}
	}

	if strings.Trim(path, " \t\n\r\x00") != path {
		return fmt.Errorf("begins or ends with whitespace or control characters")
	}

	return nil
}
