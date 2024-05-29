package jira

import (
	"crypto/tls"
	"net/http"

	gojira "github.com/andygrunwald/go-jira"
)

type client struct {
	gojira.Client
}

var _ gojiraClient = (*client)(nil)

var createBearerClient = func(token string, host string, skipTLSVerify bool) (gojiraClient, error) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: skipTLSVerify}

	tp := gojira.BearerAuthTransport{
		Token:     token,
		Transport: customTransport,
	}

	gojiraClient, err := gojira.NewClient(tp.Client(), host)

	if err != nil {
		return nil, err
	}

	return &client{*gojiraClient}, nil
}

func (c *client) getIssue(identifier string) (*gojira.Issue, *gojira.Response, error) {
	return c.Issue.Get(identifier, &gojira.GetQueryOptions{Fields: "issuetype,summary"})
}
