package config

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	"github.com/stretchr/testify/require"
)

func TestWriteTemplatedConfigFile(t *testing.T) {

	t.Run("Written config file should be the same as the expected default-config", func(t *testing.T) {
		originalGetConfigFile := GetConfigFile
		defer func() {
			GetConfigFile = originalGetConfigFile
		}()

		GetConfigFile = func() (ConfigFile, error) {
			tempDirName, err := os.MkdirTemp("", "sherpa-config-test")
			return ConfigFile{
				Path: filepath.Join(tempDirName, configPath),
				Name: configName,
				Type: configType,
			}, err
		}

		err := Initialize(false)
		require.NoError(t, err)

		defaultConfigFile, err := os.ReadFile("default-config.yml")
		require.NoError(t, err)

		cfg := GetConfig()
		var buff bytes.Buffer
		templateData := configFileTemplateData{
			JiraData: JiraTemplateConfiguration{
				Jira: cfg.Jira,
			},
			GithubData: GithubTemplateConfiguration{
				Github: cfg.Github,
			},
			BranchesData: BranchesTemplateConfiguration{
				Branches: cfg.Branches,
			},
		}

		err = writeTemplatedConfigFile(&buff, templateData)
		require.NoError(t, err)

		require.Equal(t, string(defaultConfigFile), buff.String())

	})
}

func TestJiraTemplateConfiguration(t *testing.T) {
	tmpl, err := template.ParseFS(embeddedTemplates, "templates/*.tmpl")
	require.NoError(t, err)

	t.Run("Should generate empty configuration", func(t *testing.T) {
		jiraData := JiraTemplateConfiguration{
			Jira: Jira{},
		}

		var buff bytes.Buffer
		err := tmpl.ExecuteTemplate(&buff, "jiraConfiguration", jiraData)
		require.NoError(t, err)

		require.Equal(t, `# Jira configuration -----------------------------------#
jira:
  # Jira authentication configuration
  auth:
    # Jira authentication url to generate PAT
    # WARNING: Replace it with your actual Jira authentication url
    host: 
    # Jira already generated PAT
    # WARNING: Replace it with your actual Jira PAT
    token: 
    # Jira insecure TLS configuration
    skip_tls_verify: false
  # Jira issue types configuration
  issue_types: {}
`, buff.String())
	})

	t.Run("Should generate configuration with values", func(t *testing.T) {
		jiraData := JiraTemplateConfiguration{
			Jira: Jira{
				Auth: JiraAuth{
					Host:        "https://jira.example.com",
					Token:       "jira-pat",
					InsecureTLS: true,
				},
				IssueTypes: map[issue_types.IssueType][]string{
					issue_types.Bug:         {"1"},
					issue_types.Feature:     {"2", "3", "4"},
					issue_types.Improvement: {},
				},
			},
		}

		var buff bytes.Buffer
		err := tmpl.ExecuteTemplate(&buff, "jiraConfiguration", jiraData)
		require.NoError(t, err)

		require.Equal(t, `# Jira configuration -----------------------------------#
jira:
  # Jira authentication configuration
  auth:
    # Jira authentication url to generate PAT
    # WARNING: Replace it with your actual Jira authentication url
    host: https://jira.example.com
    # Jira already generated PAT
    # WARNING: Replace it with your actual Jira PAT
    token: jira-pat
    # Jira insecure TLS configuration
    skip_tls_verify: true
  # Jira issue types configuration
  issue_types:
    bug: ["1"]
    feature: ["2", "3", "4"]
    improvement: []
`, buff.String())

	})
}

func TestGithubTemplateConfiguration(t *testing.T) {
	tmpl, err := template.ParseFS(embeddedTemplates, "templates/*.tmpl")
	require.NoError(t, err)

	t.Run("Should generate empty configuration", func(t *testing.T) {
		githubData := GithubTemplateConfiguration{
			Github: Github{},
		}

		var buff bytes.Buffer
		err := tmpl.ExecuteTemplate(&buff, "githubConfiguration", githubData)
		require.NoError(t, err)

		require.Equal(t, `# Github configuration --------------------------------#
github:
  # Github issue labels configuration
  issue_labels: {}
`, buff.String())
	})

	t.Run("Should generate configuration with values", func(t *testing.T) {
		githubData := GithubTemplateConfiguration{
			Github: Github{
				IssueLabels: map[issue_types.IssueType][]string{
					issue_types.Bugfix:        {"kind/bug", "kind/bugfix"},
					issue_types.Feature:       {"kind/feature"},
					issue_types.Refactoring:   {"kind/refactoring"},
					issue_types.Documentation: {},
				},
			},
		}

		var buff bytes.Buffer
		err := tmpl.ExecuteTemplate(&buff, "githubConfiguration", githubData)
		require.NoError(t, err)

		require.Equal(t, `# Github configuration --------------------------------#
github:
  # Github issue labels configuration
  issue_labels:
    bugfix: ["kind/bug", "kind/bugfix"]
    documentation: []
    feature: ["kind/feature"]
    refactoring: ["kind/refactoring"]
`, buff.String())
	})
}

func TestBranchesTemplateConfiguration(t *testing.T) {
	tmpl, err := template.ParseFS(embeddedTemplates, "templates/*.tmpl")
	require.NoError(t, err)

	t.Run("Should generate empty configuration", func(t *testing.T) {
		branchesData := BranchesTemplateConfiguration{
			Branches: Branches{},
		}

		var buff bytes.Buffer
		err := tmpl.ExecuteTemplate(&buff, "branchesConfiguration", branchesData)
		require.NoError(t, err)

		require.Equal(t, `# Branches configuration ------------------------------#
branches:
  # Branch prefix configuration
  prefixes: {}
`, buff.String())
	})

	t.Run("Should generate configuration with values", func(t *testing.T) {
		branchesData := BranchesTemplateConfiguration{
			Branches: Branches{
				Prefixes: map[issue_types.IssueType]string{
					issue_types.Bug:         "fix",
					issue_types.Feature:     "feat",
					issue_types.Improvement: "",
				},
			},
		}

		var buff bytes.Buffer
		err := tmpl.ExecuteTemplate(&buff, "branchesConfiguration", branchesData)
		require.NoError(t, err)

		require.Equal(t, `# Branches configuration ------------------------------#
branches:
  # Branch prefix configuration
  prefixes:
    bug: fix
    feature: feat
    improvement: 
`, buff.String())
	})

}
