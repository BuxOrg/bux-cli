# BUX: CLI
> Command line application for interacting with BUX

[![Build Status](https://img.shields.io/github/actions/workflow/status/BuxOrg/bux-cli/run-tests.yml?branch=master&logo=github&v=3)](https://github.com/BuxOrg/bux-cli/actions)
[![Report](https://goreportcard.com/badge/github.com/BuxOrg/bux-cli?style=flat&v=1)](https://goreportcard.com/report/github.com/BuxOrg/bux-cli)
[![Go](https://img.shields.io/github/go-mod/go-version/BuxOrg/bux-cli?v=1)](https://golang.org/)
<br>
[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://api.mergify.com/v1/badges/BuxOrg/bux-cli&style=flat&v=1)](https://mergify.io)
[![Sponsor](https://img.shields.io/badge/sponsor-mrz1836-181717.svg?logo=github&style=flat&v=2)](https://github.com/sponsors/mrz1836)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat&v=2)](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=bux&utm_term=bux&utm_content=bux)

<br/>

## Table of Contents
- [Installation](#installation)
- [Commands](#commands)
- [Documentation](#documentation)
- [Examples & Tests](#examples--tests)
- [Code Standards](#code-standards)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

<br/>

## Installation

**Install with [brew](https://github.com/mrz1836/homebrew-bux-cli)**
```shell script
brew tap BuxOrg/bux-cli && brew install bux-cli
buxcli
```

**Install using a [compiled binary](https://github.com/BuxOrg/bux-cli/releases)** on Linux or Mac _(Mac example)_
```shell script
curl -LkSs https://github.com/BuxOrg/bux-cli/releases/download/v0.3.24/bux-cli_macOS_64-bit.tar.gz -o app.tar.gz
tar -zxf app.tar.gz && cd ./app/
./buxcli
```

**Install with [go](https://formulae.brew.sh/formula/go)**
```shell script
go get github.com/BuxOrg/bux-cli
cd /$GOPATH/src/github.com/BuxOrg/bux-cli && make install
buxcli
```

<br/>

## Commands

### `xpub`
> Create a new xpub ([view example](docs/commands/buxcli_xpub.md))
```shell script
buxcli xpub create <xpriv>
```

> Get help for the xpub command
```shell script
buxcli xpub --help
```

<br/>

___

<br/>

### `xpriv`
> Create a xpriv key ([view example](docs/commands/buxcli_xpriv.md))
```shell script
buxcli xpriv create
```

> Generate a WIF from a xpriv ([view example](docs/commands/buxcli_xpriv.md))
```shell script
buxcli xpriv wif <xpriv>
```

> Generate a WIF from a xpriv ([view example](docs/commands/buxcli_xpriv.md))
```shell script
buxcli xpriv wif <xpriv>
```

> Get help for the xpriv command
```shell script
buxcli xpriv --help
```

<br/>

## Documentation
Get started with the [examples](docs/examples.md). View the generated golang [godocs](https://pkg.go.dev/github.com/BuxOrg/bux-cli?tab=subdirectories).

All the generated command documentation can be found in [docs/commands](docs/commands).

### Supported Operating Systems
- [x] Linux
- [x] Mac
- [ ] Windows _(coming soon)_

<br/>

<details>
<summary><strong><code>Custom Configuration</code></strong></summary>
<br/>

The configuration file should be located in your `$HOME/buxcli` folder and named `config.json`.

View the [example config file](config-example.json).

You can also specify a custom configuration file using `--config "~/folder/path/config.json"`
</details>

<details>
<summary><strong><code>Local Database (Cache)</code></strong></summary>
<br/>

The database is located in your `$HOME/buxcli` folder.

To clear the entire database:
```shell script
buxcli --flush-cache
```

Run commands _ignoring_ local cache:
```shell script
buxcli some_command --no-cache
```
</details>

<details>
<summary><strong><code>Package Dependencies</code></strong></summary>
<br/>

- [badger](https://github.com/dgraph-io/badger) for persistent database storage
- [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper) for an easy configuration & CLI application development
- [color](https://github.com/fatih/color) for colorful logs
- [columnize](https://github.com/ryanuber/columnize) for displaying terminal data in columns
- [go-homedir](https://github.com/mitchellh/go-homedir) to find the home directory
- [go-sanitize](https://github.com/mrz1836/go-sanitize) for sanitation and data formatting
- [go-validate](https://github.com/mrz1836/go-validate) for domain/email/ip validations
- [resty](https://github.com/go-resty/resty) for custom HTTP client support
</details>

<details>
<summary><strong><code>Application Deployment</code></strong></summary>
<br/>

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary deployment to GitHub and can be installed via: `brew install goreleaser`.

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production.

The release can also be deployed to a `homebrew` repository: [homebrew-bux-cli](https://github.com/mrz1836/homebrew-bux-cli).
</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>
<br/>

View all `makefile` commands
```shell script
make help
```

List of all current commands:
```text
all                      Runs multiple commands
build                    Build all binaries (darwin, linux, windows)
clean                    Remove previous builds and any test cache data
clean-mods               Remove all the Go mod cache
coverage                 Shows the test coverage
darwin                   Build for Darwin (macOS amd64)
diff                     Show the git diff
gen-docs                 Generate documentation from all available commands (fresh install)
generate                 Runs the go generate command in the base of the repo
gif-render               Render gifs in .github dir (find/replace text etc)
godocs                   Sync the latest tag with GoDocs
help                     Show this help message
install                  Install the application
install-go               Install the application (Using Native Go)
install-releaser         Install the GoReleaser application
lint                     Run the golangci-lint application (install if not found)
linux                    Build for Linux (amd64)
release                  Full production release (creates release in Github)
release                  Runs common.release then runs godocs
release-snap             Test the full release (build binaries)
release-test             Full production test release (everything except deploy)
replace-version          Replaces the version in HTML/JS (pre-deploy)
tag                      Generate a new tag and push (tag version=0.0.0)
tag-remove               Remove a tag if found (tag-remove version=0.0.0)
tag-update               Update an existing tag to current commit (tag-update version=0.0.0)
test                     Runs lint and ALL tests
test-ci                  Runs all tests via CI (exports coverage)
test-ci-no-race          Runs all tests via CI (no race) (exports coverage)
test-ci-short            Runs unit tests via CI (exports coverage)
test-no-lint             Runs just tests
test-short               Runs vet, lint and tests (excludes integration tests)
test-unit                Runs tests and outputs coverage
uninstall                Uninstall the application (and remove files)
update-linter            Update the golangci-lint package (macOS only)
update-terminalizer      Update the terminalizer application
vet                      Run the Go vet application
windows                  Build for Windows (amd64)
```
</details>

<br/>

## Examples & Tests
All unit tests and [examples](docs/examples.md) run via [GitHub Actions](https://github.com/BuxOrg/bux-cli/actions) and
uses [Go version 1.18.x](https://golang.org/doc/go1.18). View the [configuration file](.github/workflows/run-tests.yml).

Run all tests (including integration tests)
```shell script
make test
```

<br/>

## Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## Usage
View all the [examples](docs/examples.md) and see the [commands above](#commands)

All the generated command documentation can be found in [docs/commands](docs/commands).

<br/>

## Maintainers
| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) | 
|:------------------------------------------------------------------------------------------------:|
|                                [MrZ](https://github.com/mrz1836)                                 |

<br/>

## Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and please follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/mrz1836) :clap:
or by making a [**bitcoin donation**](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=bux&utm_term=bux&utm_content=bux) to ensure this journey continues indefinitely! :rocket:

[![Stars](https://img.shields.io/github/stars/BuxOrg/bux-cli?label=Please%20like%20us&style=social&v=2)](https://github.com/BuxOrg/bux-cli/stargazers)

<br/>

## License

[![License](https://img.shields.io/github/license/BuxOrg/bux-cli.svg?style=flat)](LICENSE)
