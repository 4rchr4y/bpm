package manager

import (
	"fmt"
	"strings"
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
