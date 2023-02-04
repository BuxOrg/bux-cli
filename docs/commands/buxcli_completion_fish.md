## buxcli completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	buxcli completion fish | source

To load completions for every new session, execute once:

	buxcli completion fish > ~/.config/fish/completions/buxcli.fish

You will need to start a new shell for this setup to take effect.


```
buxcli completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string   custom config file (default is $HOME/buxcli/config.json)
      --docs            generate docs from all commands (./docs/commands)
      --flush-cache     flushes ALL cache, empties local temporary database
      --no-cache        turn off caching for this specific command
      --verbose         enable verbose logging
```

### SEE ALSO

* [buxcli completion](buxcli_completion.md)	 - Generate the autocompletion script for the specified shell

