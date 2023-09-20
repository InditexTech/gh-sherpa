# Sherpa extension for GitHub CLI

**Sherpa** extension for [GitHub CLI](https://github.com/cli/cli) helps you to **automate certain operations of the
development life cycle of a task**.

<img src="docs/images/create-pr.svg" alt="alt text" width="700" height="450">

## Getting Started

### Prerequisites

* An available GitHub account.
* [**GitHub CLI**](https://github.com/cli/cli) configured in your development environment (version `2.0.0` or higher).
* **Bash**: Supports Linux, MacOS and Windows (for the latter, we recommend using `WSL2`).

### Install

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

## Contribute

Before developing any new feature or fix, please, check the [`CONTRIBUTING.md`](CONTRIBUTING.md) file.

## Development

Once time you have read the [`CONTRIBUTING.md`](CONTRIBUTING.md) file, you will need:

* Install Golang (version `1.20.4` or higher).
* Activate the development mode setting `GH_SHERPA_DEV` environment variable:

```sh
export GH_SHERPA_DEV=1
```

* Install the extension using the local path:

```sh
git clone https://github.com/InditexTech/gh-sherpa.git
cd gh-sherpa
gh extension remove sherpa && gh extension install .
```
