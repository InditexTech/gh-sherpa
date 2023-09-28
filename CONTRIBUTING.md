# Contributing

Thank you for your interest in contributing to this project! We value and appreciate any contributions you can make. To maintain a
collaborative and respectful environment, please consider the following guidelines when contributing to this project.

## How to Contribute

1. Open an issue to discuss and gather feedback on the feature or fix you wish to address.
2. Fork the repository and clone it to your local machine.
3. Create a new branch to work on your contribution: `git checkout -b your-branch-name`.
4. Make the necessary changes in your local branch.
5. Ensure that your code follows the established project style and formatting guidelines.
6. Perform testing to ensure your changes do not introduce errors.
7. Make clear and descriptive commits that explain your changes.
8. Push your branch to the remote repository: `git push origin your-branch-name`.
9. Open a pull request describing your changes and linking the corresponding issue.
10. Await comments and discussions on your pull request. Make any necessary modifications based on the received feedback.
11. Once your pull request is approved, your contribution will be merged into the main branch.

## Contribution Guidelines

- All contributors are expected to follow the project's [code of conduct](CODE_OF_CONDUCT.md). Please be respectful and considerate towards other contributors.
- Before starting work on a new feature or fix, check existing [issues](../../issues) and [pull requests](../../pulls) to avoid duplications and unnecessary discussions.
- If you wish to work on an existing issue, comment on the issue to inform other contributors that you are working on it. This will help coordinate efforts and prevent conflicts.
- It is always advisable to discuss and gather feedback from the community before making significant changes to the project's structure or architecture.
- Ensure a clean and organized commit history. Divide your changes into logical and descriptive commits.
- Document any new changes or features you add. This will help other contributors and project users understand your work and its purpose.
- Be sure to link the corresponding issue in your pull request to maintain proper tracking of contributions.

## Development

Make sure that you have:

- Read the rest of the [`CONTRIBUTING.md`](CONTRIBUTING.md) sections.
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

### Writing tests

We use [stretchr/testify suite package](https://github.com/stretchr/testify#suite-package) for testing when needed. You can also write regular tests without using the suite package.

### Mocking interfaces

We use [vektra/mockery](https://github.com/vektra/mockery) for mocking interfaces. You can generate the mocks with the following command:

```sh
make generate-mocks
```

This command will generate the mocks in the `internal/mocks` directory, as configured in the [`.mockery.yaml`](.mockery.yaml) file.

>NOTE: Please, refrain from using the generated `NewMockXXXX` constructors. Instead instantiate the mocks using `&MockXXXX{}`. This is needed because the generated constructors will always execute `mock.AssertExpectation(t)` on cleanup, which will fail if the test did not expect a call to the mock.

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

## Helpful Resources

- [Project documentation](README.md): Refer to our documentation for more information on the project structure and how to contribute.
- [Use cases](docs/USAGE.md): Check out the available use cases and examples to learn how to use this extension.
- [Architecture](docs/ARCHITECTURE.md): Learn more about the project's architecture and how it works.
- [Issues](../../issues): Check open issues and look for opportunities to contribute. Make sure to open an issue before starting work on a new feature or fix.

Thank you for your time and contribution! Your work helps to grow and improve this project. If you have any questions, feel free to reach out to us.
