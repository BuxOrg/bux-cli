## buxcli completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(buxcli completion zsh); compdef _buxcli buxcli

To load completions for every new session, execute once:

#### Linux:

	buxcli completion zsh > "${fpath[1]}/_buxcli"

#### macOS:

	buxcli completion zsh > $(brew --prefix)/share/zsh/site-functions/_buxcli

You will need to start a new shell for this setup to take effect.


```
buxcli completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
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

