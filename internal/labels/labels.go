package labels

import (
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
)

type IssueLabelsMap map[issue_types.IssueType][]string

type Configuration struct {
	IssueLabels IssueLabelsMap
}

type LabelsProvider struct {
	cfg Configuration
}

var _ domain.LabelProvider = (*LabelsProvider)(nil)

// New returns a new labels provider
func New(cfg Configuration) (*LabelsProvider, error) {
	return &LabelsProvider{
		cfg: cfg,
	}, nil
}

// GetIssueTypeLabel returns the label for the issue type
func (p LabelsProvider) GetIssueTypeLabel(issue domain.Issue) (string, error) {
	//TODO: Implement method
	return "kind/feature", nil
}
