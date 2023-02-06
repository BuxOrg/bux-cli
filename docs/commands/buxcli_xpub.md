## buxcli xpub

manage your xpubs in BUX

### Synopsis

```
____  _____________ ____ _____________ 
\   \/  /\______   \    |   \______   \
 \     /  |     ___/    |   /|    |  _/
 /     \  |    |   |    |  / |    |   \
/___/\  \ |____|   |______/  |______  /
      \_/                           \/
```

This command is for xpub (HD-Key) related commands.

new: creates a new xpub in BUX (xpub new <xpriv>)
get: get a xpub from BUX (xpub get <xpub> | <xpub_id> -m=<metadata_json>)


```
buxcli xpub [flags]
```

### Examples

```
buxcli xpub new <xpriv>
```

### Options

```
  -h, --help              help for xpub
  -m, --metadata string   Model Metadata
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

