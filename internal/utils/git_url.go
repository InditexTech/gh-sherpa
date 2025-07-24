package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func ExtractRepoFromURL(gitURL string) string {
	if strings.Contains(gitURL, "@github.com") {
		return extractFromSSHURL(gitURL)
	}

	if strings.HasPrefix(gitURL, "http://") || strings.HasPrefix(gitURL, "https://") {
		return extractFromHTTPSURL(gitURL)
	}

	return gitURL
}

func extractFromSSHURL(sshURL string) string {
	if !strings.Contains(sshURL, "@github.com") {
		return sshURL
	}

	if strings.HasPrefix(sshURL, "ssh://") {
		parts := strings.Split(sshURL, "github.com/")
		if len(parts) != 2 {
			return sshURL
		}

		pathParts := strings.Split(parts[1], "/")
		if len(pathParts) >= 2 {
			owner := pathParts[0]
			repo := strings.TrimSuffix(pathParts[1], ".git")
			return fmt.Sprintf("%s/%s", owner, repo)
		}

		return sshURL
	}

	parts := strings.Split(sshURL, "@github.com:")
	if len(parts) != 2 {
		return sshURL
	}

	repoParts := strings.Split(parts[1], "/")
	if len(repoParts) >= 2 {
		owner := repoParts[0]
		repo := strings.TrimSuffix(repoParts[1], ".git")
		return fmt.Sprintf("%s/%s", owner, repo)
	} else if len(repoParts) == 1 && strings.Contains(repoParts[0], ".git") {
		ownerRepo := strings.TrimSuffix(repoParts[0], ".git")
		return ownerRepo
	}

	return sshURL
}

func extractFromHTTPSURL(httpsURL string) string {
	parsedURL, err := url.Parse(httpsURL)
	if err != nil {
		return httpsURL
	}

	if !isGitHubHost(parsedURL.Host) {
		return httpsURL
	}

	path := strings.Trim(parsedURL.Path, "/")
	if path == "" {
		return httpsURL
	}

	path = strings.TrimSuffix(path, ".git")

	// Split path into parts and filter out empty segments
	parts := strings.Split(path, "/")
	var validParts []string
	for _, part := range parts {
		if part != "" {
			validParts = append(validParts, part)
		}
	}

	// We need exactly 2 parts: owner and repo
	if len(validParts) == 2 {
		owner := validParts[0]
		repo := validParts[1]
		return fmt.Sprintf("%s/%s", owner, repo)
	}

	// If we don't have exactly 2 valid parts, return original URL
	return httpsURL
}

func isGitHubHost(host string) bool {
	return host == "github.com" || host == "www.github.com"
}
