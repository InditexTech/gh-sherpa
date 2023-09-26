# Sherpa extension for GitHub CLI

**Sherpa** extension for [GitHub CLI](https://github.com/cli/cli) helps you to **automate certain operations of the
development life cycle of a task**.

<img src="docs/images/create-pr.svg" alt="alt text" width="700" height="450"/>

## Prerequisites

* An available GitHub account.
* [**GitHub CLI**](https://github.com/cli/cli) configured in your development environment (version `2.0.0` or higher).
* **Bash**: Supports Linux, MacOS and Windows (for the latter, we recommend using `WSL2`).

## Installation

Make sure you meet the [prerequisites](#prerequisites) first.

### Option 1: Install as a GitHub CLI extension (recommended)

<details open>
<summary>Show/close</summary>

You can **install** this extension just running this command from your terminal:

```sh
gh extension install InditexTech/gh-sherpa
```

#### Upgrade

If you have already installed this extension and you want to **upgrade** it, so, you should run this command from your terminal:

```sh
gh extension upgrade sherpa
```

#### Remove

To **remove** this extension just run:

```sh
gh extension remove sherpa
```

</details>

### Option 2: From source

Check the corresponding section in [`Build / install from source`](#build--install-from-source).

### Option 3: Using `go install`

```sh
go install github.com/InditexTech/gh-sherpa@latest
```
### Option 4: Download the binary file

You can download the binary file from the [releases page](https://github.com/InditexTech/gh-sherpa/releases)


## Usage

After installing this extension in your development environment, you can know the available commands in the [`USAGE.md`](docs/USAGE.md) file.

## Build / install from source

Make sure that you have:

* Meet the [prerequisites](#prerequisites).
* [Golang](https://golang.org/doc/install) (version `1.20.4` or higher).
* [GNU Make](https://www.gnu.org/software/make/) (version `4.2.1` or higher).
* [Git](https://git-scm.com/downloads) (version `2.25.1` or higher).

Download the source code:

```
git clone https://github.com/InditexTech/gh-sherpa.git
cd gh-sherpa
```

Then, you can **build** this extension just running this command from your terminal:

```sh
# This will generate a binary file called `gh-sherpa` in the `bin` directory.
make build
```

This command will generate a binary file called `gh-sherpa` in the `bin` directory.

You can also install it as a go binary file with the following command:

```sh
# This will install `gh-sherpa` in your `$GOPATH/bin` directory.
make install
```

## Contribute

Before developing any new feature or fix, please, check the [`CONTRIBUTING.md`](CONTRIBUTING.md) file.

## Development

First read [`CONTRIBUTING.md`](CONTRIBUTING.md) file and the [`Build / install from source`](#build--install-from-source) section.

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
