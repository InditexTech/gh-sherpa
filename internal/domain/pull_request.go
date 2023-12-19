// SPDX-FileCopyrightText: 2023 INDITEX S.A
//
// SPDX-License-Identifier: Apache-2.0

package domain

type PullRequest struct {
	Title       string
	Number      int
	State       string
	Closed      bool
	Url         string
	HeadRefName string
	BaseRefName string
	Labels      []Label
}
