// SPDX-FileCopyrightText: 2023 INDITEX S.A
//
// SPDX-License-Identifier: Apache-2.0

package domain

type Issue struct {
	ID           string
	Title        string
	Body         string
	Url          string
	IssueTracker IssueTrackerType
	// Used in GitHub
	Labels []Label
	// Used in Jira
	Type IssueType
}

type IssueType struct {
	Id          string
	Name        string
	Description string
}
