[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=InditexTech_gh-sherpa&metric=bugs)](https://sonarcloud.io/summary/new_code?id=InditexTech_gh-sherpa)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=InditexTech_gh-sherpa&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=InditexTech_gh-sherpa)[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=InditexTech_gh-sherpa&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=InditexTech_gh-sherpa)

# Sherpa extension for GitHub CLI

Sherpa for [GitHub CLI](https://github.com/cli/cli) makes it easy for you to **create branches** and **pull requests**
associated with any **GitHub or Jira issue**.

This extension retrieves the type of issue (_User Story_, _Bug_, _Technical Improvement_, etc) and creates a branch or
pull request associated with that issue, following the contribution model you define in a
[configuration file](#configuration).

![Create PR](docs/images/create-pr.svg)

## Table of contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [AI-assisted development](#ai-assisted-development)
- [Contribute](#contribute)

## Prerequisites

- An available GitHub account.
- [**GitHub CLI**](https://github.com/cli/cli) (version `2.0.0` or higher) configured and
[authenticated](https://cli.github.com/manual/gh_auth_login) in your development environment.
- **Bash**: Supports Linux, MacOS and Windows (for the latter, we recommend using
[`WSL2`](https://learn.microsoft.com/en-us/windows/wsl/install)).

## Installation

Make sure you meet the [prerequisites](#prerequisites) first.

You can **install** this extension just running this command from your terminal:

```sh
gh extension install InditexTech/gh-sherpa
```

### Upgrade

If you have already installed this extension and you want to **upgrade** it, so, you should run this command from your
terminal:

```sh
gh extension upgrade sherpa
```

### Remove

To **remove** this extension just run:

```sh
gh extension remove sherpa
```

## Configuration

Sherpa uses different issue types (`feature`, `bugfix`, `hotfix`, `refactoring`, etc) when mapping an issue with its
corresponding branch prefix.

### Default configuration

By default, it will use the [`internal/config/default-config.yml`](internal/config/default-config.yml) configuration
file to perform these mappings.

### Custom configuration

Otherwise, if you wish customize the different issue types, branch prefixes, etc, so, you can **create your own
configuration file** located in `$HOME/.config/sherpa/config.yml` from the
[default config file](internal/config/default-config.yml).

> If you are **using Jira as issue tracker**, so, the first time you run a command it will ask you to configure Jira
credentials and then proceed to create the custom configuration file with the provided Jira credentials.

## Usage

After installing this extension in your development environment, you can know the available commands in the
[`USAGE.md`](docs/USAGE.md) file.

## AI-assisted development

Sherpa is designed to be fully driven by AI coding agents without requiring any interactive terminal input.

All branch naming decisions and PR metadata can be specified via CLI flags — no stdin prompts are needed when using `-y`/`--yes` together with the new flags.

### Typical AI agent workflow

**Create a branch** for an issue with a known type:
```sh
gh sherpa create-branch --issue 42 --yes \
  --branch-type feature \
  --branch-description "implement-oauth"
# Output: feature/GH-42-implement-oauth
```

**Create a fully configured draft PR** in one command:
```sh
gh sherpa create-pr --issue 42 --yes \
  --branch-type feature \
  --branch-description "implement-oauth" \
  --pr-title "feat(auth): implement OAuth2 login" \
  --pr-body "Closes #42\n\nImplements OAuth2 login flow." \
  --reviewer alice \
  --assignee bob \
  --label "priority/high"
```

**Get machine-readable output** for chaining with other tools:
```sh
BRANCH=$(gh sherpa create-branch --issue 42 --yes --branch-type bugfix --output json | jq -r .branch)
```

**Preview without side effects**:
```sh
gh sherpa create-pr --issue 42 --yes --branch-type feature --dry-run
```

When a branch already exists for the issue, Sherpa automatically reuses it in non-interactive mode (`-y`). To opt out of this behavior and force an error instead, use `--no-use-existing-branch`.

## Contribute

Before developing any new feature or fix, please, check the [`CONTRIBUTING.md`](CONTRIBUTING.md) file. You will find
there the steps to contribute along with development and testing guidelines.

## Security

If you find a security vulnerability in this project, please, check the [`SECURITY.md`](SECURITY.md) file to know how to
report it.
