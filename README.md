# BUX: CLI
> Command line application for interacting with [BUX](https://getbux.io)

[![Release](https://img.shields.io/github/release-pre/BuxOrg/bux-cli.svg?logo=github&style=flat&v=1)](https://github.com/BuxOrg/bux-cli/releases)
[![Downloads](https://img.shields.io/github/downloads/BuxOrg/bux-cli/total.svg?logo=github&style=flat&v=1)](https://github.com/BuxOrg/bux-cli/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/BuxOrg/bux-cli/run-tests.yml?branch=master&logo=github&v=3)](https://github.com/BuxOrg/bux-cli/actions)
[![Report](https://goreportcard.com/badge/github.com/BuxOrg/bux-cli?style=flat&v=1)](https://goreportcard.com/report/github.com/BuxOrg/bux-cli)
[![Go](https://img.shields.io/github/go-mod/go-version/BuxOrg/bux-cli?v=1)](https://golang.org/)
<br>
[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://api.mergify.com/v1/badges/BuxOrg/bux-cli&style=flat&v=1)](https://mergify.io)
[![Makefile Included](https://img.shields.io/badge/Makefile-Supported%20-brightgreen?=flat&logo=probot&v=1)](Makefile)
[![Sponsor](https://img.shields.io/badge/sponsor-BuxOrg-181717.svg?logo=github&style=flat&v=1)](https://github.com/sponsors/BuxOrg)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat&v=1)](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=bux-server&utm_term=bux-server&utm_content=bux-server)

<br/>

## Table of Contents
- [Installation](#installation)
- [What is BUX?](#what-is-bux)
- [Getting Started](#getting-started)
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
curl -LkSs https://github.com/BuxOrg/bux-cli/releases/download/v0.1.0/bux-cli_macOS_64-bit.tar.gz -o app.tar.gz
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

## What is BUX?

[Read more about BUX](https://getbux.io)

<br/>

## Getting Started
The default configuration will use a [`config.json`](config-example.json) and `datastore.db` file.

These files are located in your home directory (`~/buxcli/`).

It is recommended to make changes to the `~/buxcli/config.json` file.

> Start by creating a new xpriv using the `xpriv` command.
```shell script
buxcli xpriv new
```

> Next, create a new xpub using the `xpub` command.
```shell script 
buxcli xpub new <xpriv> --metadata='{ "name": "xpub_1", "description": "my xpub description"}'
```

> Now you can create a new destination using the `destination` command.
```shell script
buxcli destination new <xpub> --metadata='{ "name": "destination_1", "description": "my destination description"}'
```

> Finally, you can record a transaction using the `transaction` command. (IE: an incoming transaction)
```shell script
buxcli transaction record <xpub> --txid=<tx_id> --metadata='{ "name": "transaction_1", "description": "my transaction description"}'
```

> Create a draft transaction using the `transaction` command. (IE: an outgoing transaction)
```shell script
buxcli transaction new <xpub> --txconfig='{"send_all_to":{"to":"1L6Tqxe..."},"expires_in":300000000000}' --metadata='{"name":"test draft tx"}'
```

<br/>

## Commands

### `destination`
> Create a new destination using optional metadata ([view example](docs/commands/buxcli_destination.md))
```shell script
buxcli destination new <xpub> --metadata='{ "name": "destination_1", "description": "my destination description"}'
```
<br/>

> Get an existing destination from id ([view example](docs/commands/buxcli_destination.md))
```shell script
buxcli destination get <destination_id> -x=<xpub_id>
```
<br/>

> Get an existing destination from locking script ([view example](docs/commands/buxcli_destination.md))
```shell script
buxcli destination get <locking_script> -x=<xpub_id>
```
<br/>

> Get an existing destination from address ([view example](docs/commands/buxcli_destination.md))
```shell script
buxcli destination get <address> -x=<xpub_id>
```
<br/>

> Get address information from BUX and [WhatsOnChain](https://whatsonchain.com) ([view example](docs/commands/buxcli_destination.md))
```shell script 
buxcli destination get <address> -x=<xpub_id> -w
```
<br/>

> Get help for the destination command
```shell script
buxcli destination --help
```

<br/>

___

<br/>

### `transaction`
> Start a new draft transaction in BUX ([view example](docs/commands/buxcli_transaction.md))
```shell script
buxcli transaction new <xpub> --txconfig='{"send_all_to":{"to":"1L6Tqxe..."},"expires_in":300000000000}' --metadata='{"name":"test draft tx"}'
```
<br/>

> Record a transaction using a Transaction ID into BUX ([view example](docs/commands/buxcli_transaction.md))
```shell script
buxcli transaction record <xpub> --txid=<tx_id> --metadata='{ "name": "transaction_1", "description": "my transaction description"}'
```
<br/>

> Record a transaction using hex into BUX ([view example](docs/commands/buxcli_transaction.md))
```shell script 
buxcli transaction record <xpub> --hex=<tx_hex> --metadata='{ "name": "transaction_1", "description": "my transaction description"}'
```
<br/>

> Record a transaction using hex and a previously generated draft id into BUX ([view example](docs/commands/buxcli_transaction.md))
```shell script 
buxcli transaction record <xpub> --hex=<tx_hex> --draft=<draft_id> --metadata='{ "name": "transaction_1", "description": "my transaction description"}'
```
<br/>

> Get transaction information from BUX ([view example](docs/commands/buxcli_transaction.md))
```shell script 
buxcli transaction info <xpub_id> --txid=<tx_id>
```
<br/>

> Get transaction information from BUX and [WhatsOnChain](https://whatsonchain.com) ([view example](docs/commands/buxcli_transaction.md))
```shell script 
buxcli transaction info <xpub_id> --txid=<tx_id> -w
```
<br/>

> Get help for the transaction command
```shell script
buxcli transaction --help
```

<br/>

___

<br/>

### `xpub`
> Create a new xpub with optional metadata ([view example](docs/commands/buxcli_xpub.md))
```shell script
buxcli xpub new <xpriv> --metadata='{ "name": "xpub_1", "description": "my xpub description"}'
```
<br/>

> Get an existing xpub record from key ([view example](docs/commands/buxcli_xpub.md))
```shell script
buxcli xpub get <xpub>
```
<br/>

> Get an existing xpub record from xpub id ([view example](docs/commands/buxcli_xpub.md))
```shell script
buxcli xpub get <xpub_id>
```
<br/>

> Get an existing xpub record using metadata ([view example](docs/commands/buxcli_xpub.md))
```shell script
buxcli xpub get -m=<metadata_json>
```
<br/>

> Get help for the xpub command
```shell script
buxcli xpub --help
```

<br/>

___

<br/>

### `xpriv`
> Create a new xpriv key ([view example](docs/commands/buxcli_xpriv.md))
```shell script
buxcli xpriv new
```
<br/>

> Get information about an existing xpriv key ([view example](docs/commands/buxcli_xpriv.md))
```shell script
buxcli xpriv info <xpriv>
```
<br/>

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
- [bux](https://github.com/BuxOrg/bux) for using the BUX interface
- [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper) for an easy configuration & CLI application development
- [color](https://github.com/fatih/color) for colorful logs
- [go-bitcoin](https://github.com/bitcoinschema/go-bitcoin) for helping with bitcoin related operations
- [go-bk](https://github.com/libsv/go-bk) for bitcoin key operations
- [go-homedir](https://github.com/mitchellh/go-homedir) to find the home directory
- [go-minercraft](https://github.com/tonicpow/go-minercraft) for MinerCraft support
- [go-whatsonchain](https://github.com/mrz1836/go-whatsonchain) for [WhatsOnChain](https://whatsonchain.com) support
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

[![Stars](https://img.shields.io/github/stars/BuxOrg/bux-cli?label=Please%20like%20us&style=social&v=1)](https://github.com/BuxOrg/bux-cli/stargazers)

<br/>

## License

[![License](https://img.shields.io/github/license/BuxOrg/bux-cli.svg?style=flat&v=1)](LICENSE)
