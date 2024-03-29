package domain

type ConfigProvider interface {
	GetConfigByIssueTracker(issueTracker string) (IssueTrackerConfig, error)
	GetConfigFilePath() string
	SaveConfigByIssueTracker(issueTracker string, config IssueTrackerConfig) error
}

type Config map[string]IssueTrackerConfig

type IssueTrackerConfig struct {
	Auth AuthConfig `json:"auth,omitempty"`
}

type AuthConfig struct {
	Host  string `json:"host,omitempty"`
	Token string `json:"token,omitempty"`
}
