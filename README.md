# Sherpa extension for GitHub CLI

**Sherpa** extension for [GitHub CLI](https://github.com/cli/cli) helps you to **automate certain operations of the
development life cycle of a task**.

![Create PR](docs/images/create-pr.svg)

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

### Upgrade

If you have already installed this extension and you want to **upgrade** it, so, you should run this command from your terminal:

```sh
gh extension upgrade sherpa
```

### Remove

To **remove** this extension just run:

```sh
gh extension remove sherpa
```

## Usage

After installing this extension in your development environment, you can know the available commands in the [`USAGE.md`](docs/USAGE.md) file.

## Configuration

Check the [`internal/config/default-config.yml`](internal/config/default-config.yml) file to see the available configuration parameters as well as the default values and some examples.
You can also find here the available GH Sherpa issue types. These values will be set if no other configuration override those values.

In order to override the default values, you can use your own configuration file located in `$HOME/.config/sherpa/config.yml` for this.

If no configuration file is found, the first time you run a command it will ask you to configure your Jira credentials (if you want to use Jira integration) and then proceed to create the configuration file with the provided Jira credentials.

## Contribute

Before developing any new feature or fix, please, check the [`CONTRIBUTING.md`](CONTRIBUTING.md) file. You will find there the steps to contribute along with development and testing guidelines.
