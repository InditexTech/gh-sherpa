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

### Possible scenarios

#### Create a pull request associated to an existing local branch

```sh
gh sherpa create-pr
```

#### Create a branch and pull request in draft-mode associated to an issue

```sh
# Create a pull request in draft-mode associated to a GitHub issue
gh sherpa create-pr -i 750

# Create a pull request in draft-mode associated to a Jira issue
gh sherpa create-pr -i SHERPA-71
```

#### Create a branch and pull request in ready for review mode associated to an issue without confirmation

```sh
gh sherpa create-pr -i 750 --yes --no-draft
```

#### Create a branch and pull request with target branch main

```sh
gh sherpa create-pr --issue SHERPA-81 --base main
```

#### Create a branch and pull request with no auto close issue

```sh
gh sherpa create-pr --issue 750 --no-close-issue
```
