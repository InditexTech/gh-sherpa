# Usage

While using GitHub issues you **must** be within a working github repository, so it can interact with the repository's issues and pull requests.

## TL;DR

```
$ gh sherpa --help

Usage:
  sherpa [command]

Available Commands:
  create-branch Create a local branch from an issue type (alias: cb)
  create-pr     Create a pull request from the current local branch or issue type (alias: cpr)
  help          Help about any command

Flags:
  -h, --help      help for sherpa
  -v, --version   version for sherpa
  -y, --yes       use the default proposed fields

Use "sherpa [command] --help" for more information about a command.
```

## Create branch

Create a git branch associated to a GitHub or Jira issue.

### Synopsis

```sh
gh sherpa create-branch, cb [flags]
```

#### Required parameters

* `--issue, -i`: GitHub or Jira issue identifier.

#### Optional parameters

* `--base, -b`: Base branch for checkout. By default is the default branch.
* `--no-fetch`: Remote branches will not be fetched.
* `--yes, -y`: The branch will be created without confirmation.
* `--fork`: Automatically set up fork for external contributors.
* `--fork-name`: Specify custom fork organization/user (e.g. MyOrg/gh-sherpa).
* `--prefer-hotfix`: Prefer hotfix branch prefix for bug issues when using non-interactive mode (`-y`). For GitHub issues, this flag checks if the `kind/bug` label is present **anywhere** in the issue's label list (not just as the first or primary label). When found, it creates a `hotfix/` branch instead of `bugfix/`, regardless of the issue's detected type or other labels present.
* `--branch-type`: Force a specific branch type prefix (e.g. `feature`, `bugfix`, `hotfix`). Bypasses issue label detection and works in both interactive and non-interactive mode.
* `--branch-description`: Force a specific branch description slug instead of deriving it from the issue title. Works in both interactive and non-interactive mode.
* `--branch-name`: Use exactly this branch name without any auto-generation. Takes priority over all other naming flags.
* `--dry-run`: Print what would happen without actually creating the branch.
* `--output`: Output format. Use `json` to get machine-readable output `{"branch":"<name>"}`. Default is human-readable text.

### Possible scenarios

#### Create a branch name associated to an issue

```sh
# Create a branch name associated to a GitHub issue
gh sherpa create-branch --issue 17

# Create a branch name associated to a Jira issue
gh sherpa create-branch --issue SHERPA-31
```

#### Create a branch name without confirmation

```sh
gh sherpa create-branch --issue 17 --yes
```

#### Create a branch from a release branch and does not git fetch

```sh
gh sherpa create-branch --issue SHERPA-31 --base release/1.3.5 --no-fetch
```

#### Create a hotfix branch from a bug issue without confirmation

When using `--prefer-hotfix`, the tool scans **all labels** on the issue looking for `kind/bug`. If found anywhere in the label list, the branch will be created with a `hotfix/` prefix instead of the default `bugfix/` prefix. This works even if `kind/bug` is not the first label or if other type labels (like `kind/feature`, `kind/internal`) are also present.

```sh
# For a GitHub issue #17 with kind/bug label (in any position)
# Creates: hotfix/GH-17-issue-description
gh sherpa create-branch --issue 17 --yes --prefer-hotfix

# Works even if the issue has multiple labels like:
# ["priority/high", "kind/feature", "kind/bug", "component/api"]
```

#### Create a branch with automatic fork setup for external contributors

```sh
# One-command fork setup and branch creation
gh sherpa create-branch --issue 32 --fork

# Custom fork organization
gh sherpa create-branch --issue 45 --fork --fork-name MyOrg/gh-sherpa
```

#### Create a branch with a specific type and description (AI-friendly)

```sh
# Force branch type and description without interactive prompts
gh sherpa create-branch --issue 42 --yes --branch-type feature --branch-description "add-auth-endpoint"
# Creates: feature/GH-42-add-auth-endpoint

# Use an exact branch name
gh sherpa create-branch --issue 42 --yes --branch-name "feature/GH-42-my-exact-name"

# Preview what would be created
gh sherpa create-branch --issue 42 --yes --branch-type feature --dry-run

# Get machine-readable output for scripting
gh sherpa create-branch --issue 42 --yes --branch-type feature --output json
# Output: {"branch":"feature/GH-42-issue-title"}
```

## Create pull request

Create a pull request associated to a GitHub or Jira issue.

### Synopsis

```sh
gh sherpa create-pr, cpr [flags]
```

#### Optional parameters

* `--issue, -i`: GitHub or Jira issue identifier.
* `--base, -b`: Base branch for checkout. By default is the default branch.
* `--no-fetch`: Remote branches will not be fetched.
* `--yes, -y`: The pull request will be created without confirmation.
* `--no-draft`: The pull request will be created in ready for review mode. By default is in draft mode.
* `--no-close-issue, -n`: The GitHub issue will not be closed when the pull request is merged. By default is closed.
* `--template`: Path to a pull request template file.
* `--fork`: Automatically set up fork for external contributors.
* `--fork-name`: Specify custom fork organization/user (e.g. MyOrg/gh-sherpa).
* `--prefer-hotfix`: Prefer hotfix branch prefix for bug issues when using non-interactive mode (`-y`). For GitHub issues, this flag checks if the `kind/bug` label is present **anywhere** in the issue's label list (not just as the first or primary label). When found, it creates a `hotfix/` branch instead of `bugfix/`, regardless of the issue's detected type or other labels present.
* `--branch-type`: Force a specific branch type prefix (e.g. `feature`, `bugfix`, `hotfix`). Bypasses issue label detection.
* `--branch-description`: Force a specific branch description slug instead of deriving it from the issue title.
* `--branch-name`: Use exactly this branch name without any auto-generation.
* `--dry-run`: Print what would happen without actually creating the PR.
* `--output`: Output format. Use `json` to get machine-readable output `{"branch":"<name>","pr_url":"<url>","draft":<bool>}`. Default is human-readable text.
* `--pr-title`: Override the auto-generated PR title.
* `--pr-body`: Override the auto-generated PR body.
* `--pr-body-file`: Read the PR body from a file (overrides `--pr-body` and `--template`).
* `--no-use-existing-branch`: Fail if a branch for this issue already exists (default non-interactive behavior is to reuse the existing branch).
* `--label`: Additional label to apply to the PR. Can be repeated: `--label bug --label priority/high`.
* `--reviewer`: Request a review from this user or team. Can be repeated: `--reviewer alice --reviewer org/team`.
* `--assignee`: Assign this user to the PR. Can be repeated: `--assignee alice`.

### Possible scenarios

#### Create a branch and pull request in draft-mode associated to an issue

```sh
# Create a pull request in draft-mode associated to a GitHub issue
gh sherpa create-pr -i 750

# Create a pull request in draft-mode associated to a Jira issue
gh sherpa create-pr -i SHERPA-71
```

#### Create a branch and pull request associated to an issue without confirmation

```sh
gh sherpa create-pr -i SHERPA-71 --yes
```

#### Create a branch and pull request given a template

```sh
gh sherpa create-pr --issue 750 --template docs/pull_request_template.md
```

#### Create a pull request associated to an existing local branch

```sh
gh sherpa create-pr
```

#### Create a branch and pull request in ready for review mode

```sh
gh sherpa create-pr -i 750 --no-draft
```

#### Create a branch and pull request with target branch main

```sh
gh sherpa create-pr --issue SHERPA-81 --base main
```

#### Create a branch and pull request with no auto close issue

```sh
gh sherpa create-pr --issue 750 --no-close-issue
```

#### Create a hotfix pull request from a bug issue without confirmation

When using `--prefer-hotfix`, the tool scans **all labels** on the issue looking for `kind/bug`. If found anywhere in the label list, the branch and pull request will be created with a `hotfix/` prefix instead of the default `bugfix/` prefix. This works even if `kind/bug` is not the first label or if other type labels are also present.

```sh
# For a GitHub issue #750 with kind/bug label (in any position)
# Creates: hotfix/GH-750-issue-description
gh sherpa create-pr --issue 750 --yes --prefer-hotfix

# Works even if the issue has multiple labels like:
# ["priority/high", "kind/feature", "kind/bug", "component/api"]
```

#### Create a pull request with automatic fork setup for external contributors

```sh
# One-command fork setup and pull request creation
gh sherpa create-pr --issue 32 --fork

# Custom fork organization
gh sherpa create-pr --issue 45 --fork --fork-name MyOrg/gh-sherpa
```

#### Create a fully scripted PR (AI-friendly, zero interactive prompts)

```sh
# All parameters specified — no interactive prompts, no stdin required
gh sherpa create-pr \
  --issue 42 \
  --yes \
  --branch-type feature \
  --branch-description "add-auth-endpoint" \
  --pr-title "feat: add authentication endpoint" \
  --pr-body "Closes #42" \
  --reviewer alice \
  --reviewer org/backend-team \
  --assignee bob \
  --label "priority/high" \
  --no-draft

# With JSON output for scripting
gh sherpa create-pr --issue 42 --yes --branch-type bugfix --output json
# Output: {"branch":"bugfix/GH-42-issue-title","pr_url":"https://...","draft":true}

# Preview without executing
gh sherpa create-pr --issue 42 --yes --branch-type feature --dry-run

# Fail if a branch already exists (instead of reusing it)
gh sherpa create-pr --issue 42 --yes --no-use-existing-branch
```

## Fork Configuration

For external contributors working via forks, Sherpa provides seamless fork management through the `--fork` flag. This feature automates the entire fork setup process.

### Configuration

You can set a default fork organization in your configuration file (`~/.config/sherpa/config.yml`):

```yaml
github:
  fork_organization: "MyOrg"
```

### Fork Workflow Examples

**First-time Fork Setup:**
```bash
❯ gh sherpa create-branch -i 32 --fork
=> Detecting repository setup...
=> No fork detected. Creating fork danielfn/gh-sherpa...
✓ Fork created successfully
=> Setting up remotes (origin: fork, upstream: original)...
=> Setting default repository to upstream...
=> Fetching branches from fork...
=> Creating branch bugfix/GH-32-fix-link-to-cla...
✓ Ready to start working on issue #32!
```

**Subsequent Usage:**
```bash
❯ gh sherpa create-branch -i 45 --fork
=> Fork already configured, creating branch...
=> Creating branch feature/GH-45-new-feature...
✓ Ready to start working on issue #45!
```

### What the `--fork` flag does:

1. **Detects repository state** - Checks if already in a fork setup
2. **Creates fork if needed** - Runs `gh repo fork --remote` with user confirmation
3. **Sets upstream as default** - Runs `gh repo set-default <upstream-repo>`
4. **Fetches from fork** - Runs `git fetch origin` to sync branches
5. **Proceeds with standard operation** - Creates branch/PR with correct remotes
