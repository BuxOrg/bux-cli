## buxcli completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(buxcli completion bash)

To load completions for every new session, execute once:

#### Linux:

	buxcli completion bash > /etc/bash_completion.d/buxcli

#### macOS:

	buxcli completion bash > $(brew --prefix)/etc/bash_completion.d/buxcli

You will need to start a new shell for this setup to take effect.


```
buxcli completion bash
```

### Options

```
  -h, --help              help for bash
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

