package utils

import (
	"fmt"
	"strings"
)

// ExtractRepoFromURL extracts owner/repo from a git URL
func ExtractRepoFromURL(url string) string {
	if strings.Contains(url, "github.com") {
		parts := strings.Split(url, "/")
		if len(parts) >= 2 {
			owner := parts[len(parts)-2]
			repo := parts[len(parts)-1]

			repo = strings.TrimSuffix(repo, ".git")

			if strings.Contains(owner, ":") {
				owner = strings.Split(owner, ":")[1]
			}

			return fmt.Sprintf("%s/%s", owner, repo)
		}
	}

	return url
}
