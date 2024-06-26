package issue_types

import "slices"

type IssueType string

func (it IssueType) String() string {
	return string(it)
}

// Valid returns true if the IssueType is a valid one.
func (it IssueType) Valid() bool {
	return slices.Contains(GetValidIssueTypes(), it)
}

const (
	Bug           IssueType = "bug"
	Bugfix        IssueType = "bugfix"
	Dependency    IssueType = "dependency"
	Deprecation   IssueType = "deprecation"
	Documentation IssueType = "documentation"
	Feature       IssueType = "feature"
	Hotfix        IssueType = "hotfix"
	Improvement   IssueType = "improvement"
	Internal      IssueType = "internal"
	Refactoring   IssueType = "refactoring"
	Release       IssueType = "release"
	Removal       IssueType = "removal"
	Revert        IssueType = "revert"
	Security      IssueType = "security"
	Other         IssueType = "other"
	Unknown       IssueType = "unknown"
)

func GetValidIssueTypes() []IssueType {
	return []IssueType{
		Bugfix,
		Dependency,
		Deprecation,
		Documentation,
		Feature,
		Hotfix,
		Improvement,
		Internal,
		Refactoring,
		Release,
		Removal,
		Revert,
		Security,
	}
}

func GetBugValues() []IssueType {
	return []IssueType{
		Bugfix,
		Hotfix,
		Other,
	}
}

func GetAllValues() []IssueType {
	return append(GetValidIssueTypes(), Other)
}

// ParseIssueType parses a string into an IssueType.
func ParseIssueType(s string) IssueType {
	result := IssueType(s)

	for _, issueType := range GetAllValues() {
		if result == issueType {
			return result
		}
	}

	return Unknown
}
