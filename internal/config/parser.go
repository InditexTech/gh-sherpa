package config

import (
	"embed"
	"io"
	"text/template"
)

//go:embed templates/*.tmpl
var embeddedTemplates embed.FS

type configFileTemplateData struct {
	JiraData     JiraTemplateConfiguration
	GithubData   GithubTemplateConfiguration
	BranchesData BranchesTemplateConfiguration
}

type JiraTemplateConfiguration struct {
	Jira
}

type GithubTemplateConfiguration struct {
	Github
}

type BranchesTemplateConfiguration struct {
	Branches
}

func writeTemplatedConfigFile(wr io.Writer, templateData configFileTemplateData) error {
	t, err := template.ParseFS(embeddedTemplates, "templates/*.tmpl")
	if err != nil {
		return err
	}

	if err := t.ExecuteTemplate(wr, "configuration", templateData); err != nil {
		return err
	}

	return nil
}
