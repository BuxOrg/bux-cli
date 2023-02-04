## buxcli

Command line app for interacting with a BUX database or server

### Synopsis

```
__________ ____ _______  ___         _________ .____    .___ 
\______   \    |   \   \/  /         \_   ___ \|    |   |   |
 |    |  _/    |   /\     /   ______ /    \  \/|    |   |   |
 |    |   \    |  / /     \  /_____/ \     \___|    |___|   |
 |______  /______/ /___/\  \          \______  /_______ \___|
        \/               \_/                 \/        \/  v0.1.0
```
Author: MrZ © 2023 github.com/BuxOrg/bux-cli

This CLI app is used for interacting with BUX databases or servers.

Learn more about BUX: https://GetBux.io


### Examples

```
buxcli -h
```

### Options

```
      --config string   custom config file (default is $HOME/buxcli/config.json)
      --docs            generate docs from all commands (./docs/commands)
      --flush-cache     flushes ALL cache, empties local temporary database
  -h, --help            help for buxcli
      --no-cache        turn off caching for this specific command
      --verbose         enable verbose logging
  -v, --version         version for buxcli
```

### SEE ALSO

* [buxcli completion](buxcli_completion.md)	 - Generate the autocompletion script for the specified shell
* [buxcli destination](buxcli_destination.md)	 - manage your destinations in BUX
* [buxcli transaction](buxcli_transaction.md)	 - manage your transactions in BUX
* [buxcli xpriv](buxcli_xpriv.md)	 - create or derive xpriv keys
* [buxcli xpub](buxcli_xpub.md)	 - manage your xpubs in BUX

