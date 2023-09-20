package config

import (
	"crypto/tls"
	"fmt"
	"net/http"

	gojira "github.com/andygrunwald/go-jira"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	"github.com/InditexTech/gh-sherpa/internal/interactive"
)

// Jira configuration
type Jira struct {
	Auth       JiraAuth
	IssueTypes JiraIssueTypes `mapstructure:"issue_types"`
}

// JiraAuth Jira authentication configuration
type JiraAuth struct {
	Host  string
	Token string
}

// JiraIssueTypes Jira issue types mapping configuration
type JiraIssueTypes map[issue_types.IssueType][]string

type patRequestBody struct {
	Name               string `json:"name,omitempty" structs:"name,omitempty"`
	ExpirationDuration int    `json:"expirationDuration,omitempty" structs:"expirationDuration,omitempty"`
}

type patResponseBody struct {
	Id         int    `json:"id,omitempty" structs:"id,omitempty"`
	Name       string `json:"name,omitempty" structs:"name,omitempty"`
	CreatedAt  string `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	ExpiringAt string `json:"expiringAt,omitempty" structs:"expiringAt,omitempty"`
	RawToken   string `json:"rawToken,omitempty" structs:"rawToken,omitempty"`
}

func configureJira() error {
	configuredHost := vip.GetString("jira.auth.host")
	host, pat, username, password, patName, err := interactive.AskUserForJiraInputs(configuredHost)
	if err != nil {
		return err
	}
	if pat == "" {
		pat, err = generateJiraPAT(host, username, password, patName)
		if err != nil {
			return err
		}
	}

	vip.Set("jira.auth.host", host)
	vip.Set("jira.auth.token", pat)

	return nil
}

func generateJiraPAT(host, username, password, name string) (pat string, err error) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	tp := gojira.BasicAuthTransport{
		Username:  username,
		Password:  password,
		Transport: customTransport,
	}

	client, err := gojira.NewClient(tp.Client(), host)
	if err != nil {
		return
	}

	patReqBody := patRequestBody{
		Name: name,
	}
	patRequest, err := client.NewRequest("POST", "/rest/pat/latest/tokens", patReqBody)
	if err != nil {
		return "", err
	}

	var patBody patResponseBody
	res, err := client.Do(patRequest, &patBody)
	if err != nil {
		if res == nil {
			err = fmt.Errorf("could not get response from host '%s'. Check your jira configuration", host)
			return "", err
		}

		if res.StatusCode == http.StatusUnauthorized {
			return "", fmt.Errorf("your username or password is invalid")
		}

		return "", fmt.Errorf("could not create a PAT: %w", err)
	}

	return patBody.RawToken, nil
}
