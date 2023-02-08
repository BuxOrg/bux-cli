## buxcli xpriv

create new xpriv keys and see additional info

### Synopsis

```
____  _______________________._______   ____
\   \/  /\______   \______   \   \   \ /   /
 \     /  |     ___/|       _/   |\   Y   / 
 /     \  |    |    |    |   \   | \     /  
/___/\  \ |____|    |____|_  /___|  \___/   
      \_/                  \/
```

This command is for xpriv key related commands. These commands are read-only and no data is stored on the BUX servers.

new: creates a new xpriv key (xpriv new)
info: gets the xpub, WIF and other info from the xpriv key (xpriv info <xpriv>)


```
buxcli xpriv [flags]
```

### Examples

```
buxcli xpriv new
```

### Options

```
  -h, --help   help for xpriv
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

