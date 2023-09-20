package use_cases

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcIssueContextMaxLen(t *testing.T) {
	tests := []struct {
		name       string
		repository string
		branchType string
		issueId    string
		want       int
	}{
		{
			name:       "gh-sherpa -> feature/GH-1",
			repository: "gh-sherpa",
			branchType: "feature",
			issueId:    "GH-1",
			want:       41,
		},
		{
			name:       "reallyasuperlongrepositorynamethatisnotreallycommon -> feature/GH-1",
			repository: "reallyasuperlongrepositorynamethatisnotreallycommon",
			branchType: "feature",
			issueId:    "GH-1",
			want:       0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := calcIssueContextMaxLen(tt.repository, tt.branchType, tt.issueId)

			assert.Equal(t, tt.want, context)
		})
	}
}
