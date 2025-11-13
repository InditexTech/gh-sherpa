package config

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPatResponseBody_Int64Support(t *testing.T) {
	tests := []struct {
		name        string
		jsonPayload string
		expectedId  int64
		wantType    string
		description string
	}{
		{
			name:        "PAT ID within int32 range",
			jsonPayload: `{"id": 1234567890, "name": "test-token", "rawToken": "abc123"}`,
			expectedId:  1234567890,
			wantType:    "int64",
			description: "Typical PAT ID",
		},
		{
			name:        "PAT ID at int32 boundary",
			jsonPayload: `{"id": 2147483647, "name": "test-token", "rawToken": "abc123"}`,
			expectedId:  2147483647,
			wantType:    "int64",
			description: "Maximum int32 value",
		},
		{
			name:        "PAT ID exceeding int32 max",
			jsonPayload: `{"id": 2147483648, "name": "test-token", "rawToken": "abc123"}`,
			expectedId:  2147483648,
			wantType:    "int64",
			description: "One above int32 max",
		},
		{
			name:        "Large PAT ID",
			jsonPayload: `{"id": 7486581160, "name": "test-token", "rawToken": "abc123"}`,
			expectedId:  7486581160,
			wantType:    "int64",
			description: "Real-world large ID that caused the bug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var patResp patResponseBody

			err := json.Unmarshal([]byte(tt.jsonPayload), &patResp)
			require.NoError(t, err, "Should unmarshal JSON without error")

			// Verify the PAT ID is stored correctly as int64
			assert.Equal(t, tt.expectedId, patResp.Id, tt.description)

			// Verify the type is int64 using reflection
			idType := fmt.Sprintf("%T", patResp.Id)
			assert.Equal(t, tt.wantType, idType, "patResponseBody.Id should be int64 type")
		})
	}
}

func TestPatResponseBody_UnmarshalLargeValues(t *testing.T) {
	// This test specifically ensures we don't get the original error:
	// "json: cannot unmarshal number 7486581160 into Go struct field Label.Labels.Id of type int"

	largeIdJSON := `{
		"id": 7486581160,
		"name": "gh-sherpa-token",
		"createdAt": "2025-01-01T00:00:00Z",
		"expiringAt": "2025-12-31T23:59:59Z",
		"rawToken": "sample-token-value"
	}`

	var patResp patResponseBody
	err := json.Unmarshal([]byte(largeIdJSON), &patResp)

	require.NoError(t, err, "Should not fail when unmarshaling large ID values")
	assert.Equal(t, int64(7486581160), patResp.Id)
	assert.Equal(t, "gh-sherpa-token", patResp.Name)
	assert.Equal(t, "sample-token-value", patResp.RawToken)
}
