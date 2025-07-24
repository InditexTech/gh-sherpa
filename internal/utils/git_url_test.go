package utils

import (
	"testing"
)

func TestExtractRepoFromURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		// HTTPS URLs
		{
			name:     "HTTPS URL with github.com",
			url:      "https://github.com/owner/repo.git",
			expected: "owner/repo",
		},
		{
			name:     "HTTPS URL without .git suffix",
			url:      "https://github.com/owner/repo",
			expected: "owner/repo",
		},
		// SSH URLs
		{
			name:     "SSH URL with git@github.com",
			url:      "git@github.com:owner/repo.git",
			expected: "owner/repo",
		},
		{
			name:     "SSH URL without .git suffix",
			url:      "git@github.com:owner/repo",
			expected: "owner/repo",
		},

		// Organization and user repositories
		{
			name:     "Organization repository",
			url:      "https://github.com/InditexTech/gh-sherpa.git",
			expected: "InditexTech/gh-sherpa",
		},
		{
			name:     "User repository",
			url:      "https://github.com/m4rii0/gh-sherpa.git",
			expected: "m4rii0/gh-sherpa",
		},

		// Special characters in repository names
		{
			name:     "Repository with hyphens",
			url:      "https://github.com/owner/my-awesome-repo.git",
			expected: "owner/my-awesome-repo",
		},
		{
			name:     "Repository with underscores",
			url:      "https://github.com/owner/my_awesome_repo.git",
			expected: "owner/my_awesome_repo",
		},
		{
			name:     "Repository with numbers",
			url:      "https://github.com/owner/repo123.git",
			expected: "owner/repo123",
		},
		{
			name:     "Repository with dots",
			url:      "https://github.com/owner/repo.name.git",
			expected: "owner/repo.name",
		},

		// Edge cases with different URL formats
		{
			name:     "URL with www subdomain",
			url:      "https://www.github.com/owner/repo.git",
			expected: "owner/repo",
		},
		// SSH URL variations
		{
			name:     "SSH URL with ssh:// protocol",
			url:      "ssh://git@github.com/owner/repo.git",
			expected: "owner/repo",
		},
		{
			name:     "SSH URL with different user",
			url:      "user@github.com:owner/repo.git",
			expected: "owner/repo",
		},
		// Non-GitHub URLs (should return original URL)
		{
			name:     "GitLab URL",
			url:      "https://gitlab.com/owner/repo.git",
			expected: "https://gitlab.com/owner/repo.git",
		},
		{
			name:     "Bitbucket URL",
			url:      "https://bitbucket.org/owner/repo.git",
			expected: "https://bitbucket.org/owner/repo.git",
		},
		{
			name:     "Self-hosted Git server",
			url:      "https://git.company.com/owner/repo.git",
			expected: "https://git.company.com/owner/repo.git",
		},
		// Invalid or malformed URLs
		{
			name:     "Empty URL",
			url:      "",
			expected: "",
		},
		{
			name:     "Malformed URL",
			url:      "not-a-url",
			expected: "not-a-url",
		},
		// URLs with different formats
		{
			name:     "Owner with hyphens",
			url:      "https://github.com/my-organization/repo.git",
			expected: "my-organization/repo",
		},
		{
			name:     "Owner with dots",
			url:      "https://github.com/my.organization/repo.git",
			expected: "my.organization/repo",
		},
		// Edge cases for incomplete URLs
		{
			name:     "URL with only owner (no repo)",
			url:      "https://github.com/owner",
			expected: "https://github.com/owner",
		},
		{
			name:     "URL with only owner and trailing slash",
			url:      "https://github.com/owner/",
			expected: "https://github.com/owner/",
		},
		{
			name:     "URL with empty path segments",
			url:      "https://github.com/owner//repo",
			expected: "owner/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractRepoFromURL(tt.url)
			if result != tt.expected {
				t.Errorf("ExtractRepoFromURL(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}
