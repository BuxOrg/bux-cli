## buxcli transaction

manage your transactions in BUX

### Synopsis

```
_____________________    _____    _______    _________   _____  ____________________.___________    _______   
\__    ___/\______   \  /  _  \   \      \  /   _____/  /  _  \ \_   ___ \__    ___/|   \_____  \   \      \  
  |    |    |       _/ /  /_\  \  /   |   \ \_____  \  /  /_\  \/    \  \/ |    |   |   |/   |   \  /   |   \ 
  |    |    |    |   \/    |    \/    |    \/        \/    |    \     \____|    |   |   /    |    \/    |    \
  |____|    |____|_  /\____|__  /\____|__  /_______  /\____|__  /\______  /|____|   |___\_______  /\____|__  /
                   \/         \/         \/        \/         \/        \/                      \/         \/
```

This command is for transaction related commands.

record: records a new transaction in BUX (transaction record <xpub> -i=<tx_id>)


```
buxcli transaction [flags]
```

### Examples

```
buxcli record <xpub> -i=<tx_id>
```

### Options

```
  -d, --draft string      Draft ID (optional)
  -h, --help              help for transaction
  -x, --hex string        Transaction Hex
  -m, --metadata string   Model Metadata
  -i, --txid string       Transaction ID
  -w, --woc               Optional flag to use WhatsOnChain for additional transaction data
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

* [buxcli](buxcli.md)	 - Command line app for interacting with a BUX database or server

