<!-- markdownlint-disable MD033 -->
<!-- omit from toc -->
# Sherpa extension for GitHub CLI

**Sherpa** extension for [GitHub CLI](https://github.com/cli/cli) helps you to **automate certain operations of the
development life cycle of a task**.

<img src="docs/images/create-pr.svg" alt="alt text" width="700" height="450"/>

<!-- omit from toc -->
## Table of contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Contribute](#contribute)

## Prerequisites

- An available GitHub account.
- [**GitHub CLI**](https://github.com/cli/cli) configured (and [authenticated](https://cli.github.com/manual/gh_auth_login)) in your development environment (version `2.0.0` or higher).
- **Bash**: Supports Linux, MacOS and Windows (for the latter, we recommend using [`WSL2`](https://learn.microsoft.com/en-us/windows/wsl/install)).

## Installation

Make sure you meet the [prerequisites](#prerequisites) first.

You can **install** this extension just running this command from your terminal:

```sh
gh extension install InditexTech/gh-sherpa
```

<!-- omit from toc -->
### Upgrade

If you have already installed this extension and you want to **upgrade** it, so, you should run this command from your terminal:

```sh
gh extension upgrade sherpa
```

<!-- omit from toc -->
### Remove

To **remove** this extension just run:

```sh
gh extension remove sherpa
```

</details>

## Usage

After installing this extension in your development environment, you can know the available commands in the [`USAGE.md`](docs/USAGE.md) file.

## Configuration

Sherpa CLI can be configured using its own configuration file, stored in `$HOME/.config/sherpa/config.yml`.

If no configuration file is found, the first time you run a command it will ask you to configure your Jira credentials (if you want to use Jira integration) and then proceed to create the configuration file.

Check the [`internal/config/default-config.yml`](internal/config/default-config.yml) file to see the available configuration parameters as well as the default values and some examples.

<!-- omit from toc -->
### Jira configuration

>NOTE: This configuration is only required if you want to use Jira integration.

| Parameter                      | Description                                   | Default value |
| ------------------------------ | --------------------------------------------- | ------------- |
| `jira.auth.host`               | Jira host to connect to.                      | `""`          |
| `jira.auth.token`              | Jira already generated PAT                    | `""`          |
| `jira.auth.skip_tls_verify`    | Skip TLS verification for the given hos       | `false`       |
| `jira.issue_types.bugfix`      | List of Jira types ID related to bugfixes     | `["1"]`       |
| `jira.issue_types.feature`     | List of Jira types ID related to features     | `["3", "5"]`  |
| `jira.issue_types.improvement` | List of Jira types ID related to improvements | `["4"]`       |

<!-- omit from toc -->
### GitHub configuration

| Parameter                           | Description                             | Default value            |
| ----------------------------------- | --------------------------------------- | ------------------------ |
| `github.issue_labels`               | Github issue labels related to tasks    | *See lines below*        |
| `github.issue_labels.bugfix`        | List of labels related to bugfixes      | `["kind/bug]`            |
| `github.issue_labels.feature`       | List of labels related to features      | `["kind/feature"]`       |
| `github.issue_labels.refactoring`   | List of labels related to refactoring   | `["kind/refactoring"]`   |
| `github.issue_labels.documentation` | List of labels related to documentation | `["kind/documentation"]` |
| `github.issue_labels.improvement`   | List of labels related to improvements  | `["kind/improvement"]`   |

<!-- omit from toc -->
### Branches configuration

| Parameter           | Description                             | Default value |
| ------------------- | --------------------------------------- | ------------- |
| `branches.prefixes` | Branch prefix related to the issue type | `{}`          |

>NOTE: By default it will match the issue type name with the branch prefix. For example, if the issue type name is `bugfix` it will match the branch prefix `bugfix`.

## Contribute

Before developing any new feature or fix, please, check the [`CONTRIBUTING.md`](CONTRIBUTING.md) file. You will find there the steps to contribute along with development and testing guidelines.
