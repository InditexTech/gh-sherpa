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
- [Development](#development)
- [Testing the application](#testing-the-application)

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

Before developing any new feature or fix, please, check the [`CONTRIBUTING.md`](CONTRIBUTING.md) file.

## Development

Make sure that you have:

- Read [`CONTRIBUTING.md`](CONTRIBUTING.md)
- Meet the [prerequisites](#prerequisites).
- [Golang](https://golang.org/doc/install) (version `1.20.4` or higher).
- [GNU Make](https://www.gnu.org/software/make/) (version `4.2.1` or higher).
- [Git](https://git-scm.com/downloads) (version `2.25.1` or higher).

Activate the development mode setting `GH_SHERPA_DEV` environment variable:

```sh
export GH_SHERPA_DEV=1
```

Install the extension using the local path:

```sh
git clone https://github.com/InditexTech/gh-sherpa.git
cd gh-sherpa
gh extension remove sherpa && gh extension install .
```

>NOTE: You can also use `make install` to install the extension as a binary in your `$GOPATH/bin` or just run the generated binary after a `make build` execution with `./bin/gh-sherpa`.

## Testing the application

You can run the tests with the following command:

```sh
make test
```

<!-- omit from toc -->
### Writing tests

We use [stretchr/testify suite package](https://github.com/stretchr/testify#suite-package) for testing when needed. You can also write regular tests without using the suite package.

<!-- omit from toc -->
### Mocking interfaces

We use [vektra/mockery](https://github.com/vektra/mockery) for mocking interfaces. You can generate the mocks with the following command:

```sh
make generate-mocks
```

This command will generate the mocks in the `internal/mocks` directory, as configured in the [`.mockery.yaml`](.mockery.yaml) file.

>NOTE: Please, refrain from using the generated `NewMockXXXX` constructors. Instead instantiate the mocks using `&MockXXXX{}`. This is needed because the generated constructors will always execute `mock.AssertExpectation(t)` on cleanup, which will fail if the test did not expect a call to the mock.

<!-- omit from toc -->
### Coverage report

You can also run the tests with coverage with the following command:

```sh
make coverage
```

It will generate a `coverage.out` file in the `.local` directory. You can see the coverage report running the following command:

```sh
go tool cover -html=.local/coverage.out
```

It will generate an HTML file with the coverage report that you can open in your browser.

>NOTE: For Windows WSL users, you may need to convert the `coverage.out` file to a Windows compatible path. You can do it with the following command:
> ```sh
> wslpath -w PATH/TO/GENERATED/HTML/FILE
> ```
