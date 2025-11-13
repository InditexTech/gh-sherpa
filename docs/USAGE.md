# Usage

While using GitHub issues you **must** be within a working github repository, so it can interact with the repository's issues and pull requests.

## TL;DR

```
$ gh sherpa --help

Usage:
  sherpa [command]

Available Commands:
  create-branch Create a local branch from an issue type
  create-pr     Create a pull request from the current local branch or issue type
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

* `--base`: Base branch for checkout. By default is the default branch.
* `--no-fetch`: Remote branches will not be fetched.
* `--yes, -y`: The branch will be created without confirmation.
* `--fork`: Automatically set up fork for external contributors.
* `--fork-name`: Specify custom fork organization/user (e.g. MyOrg/gh-sherpa).
* `--prefer-hotfix`: Prefer hotfix branch prefix for bug issues when using non-interactive mode (`-y`).

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

```sh
gh sherpa create-branch --issue 17 --yes --prefer-hotfix
```

#### Create a branch with automatic fork setup for external contributors

```sh
# One-command fork setup and branch creation
gh sherpa create-branch --issue 32 --fork

# Custom fork organization
gh sherpa create-branch --issue 45 --fork --fork-name MyOrg/gh-sherpa
```

## Create pull request

Create a pull request associated to a GitHub or Jira issue.

### Synopsis

```sh
gh sherpa create-pr, cpr [flags]
```

#### Optional parameters

* `--issue, -i`: GitHub or Jira issue identifier.
* `--base`: Base branch for checkout. By default is the default branch.
* `--no-fetch`: Remote branches will not be fetched.
* `--yes, -y`: The pull request will be created without confirmation.
* `--no-draft`: The pull request will be created in ready for review mode. By default is in draft mode.
* `--no-close-issue`: The GitHub issue will not be closed when the pull request is merged. By default is closed.
* `--template`: Path to a pull request template file
* `--fork`: Automatically set up fork for external contributors.
* `--fork-name`: Specify custom fork organization/user (e.g. MyOrg/gh-sherpa).
* `--prefer-hotfix`: Prefer hotfix branch prefix for bug issues when using non-interactive mode (`-y`).

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

```sh
gh sherpa create-pr --issue 750 --yes --prefer-hotfix
```

#### Create a pull request with automatic fork setup for external contributors

```sh
# One-command fork setup and pull request creation
gh sherpa create-pr --issue 32 --fork

# Custom fork organization
gh sherpa create-pr --issue 45 --fork --fork-name MyOrg/gh-sherpa
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
