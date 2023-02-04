## BUX CLI: Examples & Docs
Below are some examples using the **buxcli** app

### View All Commands (Help)
```shell script
buxcli
```

```text
Available Commands:
  completion  generate the autocompletion script for the specified shell
  destination manage your destinations in BUX
  help        help about any command
  xpriv       create or derive xpriv keys
  xpub        manage your xpubs in BUX
```

Global flags for the entire application [(command specs)](commands/buxcli.md)
```text
      --config string   custom config file (default is $HOME/buxcli/config.json)
      --docs            generate docs from all commands (./docs/commands)
      --flush-cache     flushes ALL cache, empties local temporary database
  -h, --help            help for buxcli
      --no-cache        turn off caching for this specific command
      --verbose         enable verbose logging
  -v, --version         version for buxcli

```