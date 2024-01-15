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
	cfg                  Configuration
	issueTrackerProvider domain.IssueTrackerProvider
}

var _ domain.LabelProvider = (*LabelsProvider)(nil)

// New returns a new labels provider
func New(cfg Configuration, issueTrackerProvider domain.IssueTrackerProvider) (*LabelsProvider, error) {
	return &LabelsProvider{
		cfg:                  cfg,
		issueTrackerProvider: issueTrackerProvider,
	}, nil
}

// GetIssueTypeLabel returns the label for the issue type
func (p LabelsProvider) GetIssueTypeLabel(issue domain.Issue) (issueTypeLabel string, err error) {
	issueTracker, err := p.issueTrackerProvider.GetIssueTracker(issue.ID)
	if err != nil {
		return
	}

	issueType, err := issueTracker.GetIssueType(issue)
	if err != nil {
		return
	}

	for mappedIssueType, labels := range p.cfg.IssueLabels {
		if issueType == mappedIssueType {
			if len(labels) > 0 {
				return labels[0], nil
			}
		}
	}

	return
}
