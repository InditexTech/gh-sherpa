package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPullRequest_Int64Support(t *testing.T) {
	tests := []struct {
		name        string
		prNumber    int64
		wantType    string
		description string
	}{
		{
			name:        "PR number within int32 range",
			prNumber:    1234567,
			wantType:    "int64",
			description: "Typical PR number",
		},
		{
			name:        "PR number at int32 boundary",
			prNumber:    2147483647,
			wantType:    "int64",
			description: "Maximum int32 value",
		},
		{
			name:        "PR number exceeding int32 max",
			prNumber:    2147483648,
			wantType:    "int64",
			description: "One above int32 max",
		},
		{
			name:        "Large PR number",
			prNumber:    9999999999,
			wantType:    "int64",
			description: "Very large PR number for future-proofing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PullRequest{
				Title:       "Test PR",
				Number:      tt.prNumber,
				State:       "open",
				Closed:      false,
				Url:         "https://github.com/example/repo/pull/123",
				HeadRefName: "feature/test",
				BaseRefName: "main",
				Labels:      []Label{},
				Body:        "Test body",
			}

			// Verify the PR number is stored correctly as int64
			assert.Equal(t, tt.prNumber, pr.Number, tt.description)

			// Verify the type is int64 using reflection
			prNumberType := fmt.Sprintf("%T", pr.Number)
			assert.Equal(t, tt.wantType, prNumberType, "PR.Number should be int64 type")
		})
	}
}
