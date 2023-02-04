## buxcli xpriv

create or derive xpriv keys

### Synopsis

```
____  _______________________._______   ____
\   \/  /\______   \______   \   \   \ /   /
 \     /  |     ___/|       _/   |\   Y   / 
 /     \  |    |    |    |   \   | \     /  
/___/\  \ |____|    |____|_  /___|  \___/   
      \_/                  \/
```

This command is for xpriv key related commands.

create: creates a new xpriv key (xpriv create)
wif: gets the WIF from the xpriv key (xpriv wif <xpriv>)
xpub: gets the xpub from the xpriv key (xpriv xpub <xpriv>)


```
buxcli xpriv [flags]
```

### Examples

```
buxcli xpriv create
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

