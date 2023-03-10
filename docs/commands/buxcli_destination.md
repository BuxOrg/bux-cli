## buxcli destination

manage and interact with destinations in BUX

### Synopsis

```
________  ___________ ____________________.___ _______      ________________.___________    _______   
\______ \ \_   _____//   _____/\__    ___/|   |\      \    /  _  \__    ___/|   \_____  \   \      \  
 |    |  \ |    __)_ \_____  \   |    |   |   |/   |   \  /  /_\  \|    |   |   |/   |   \  /   |   \ 
 |    '   \|        \/        \  |    |   |   /    |    \/    |    \    |   |   /    |    \/    |    \
 /_______  /_______  /_______  /  |____|   |___\____|__  /\____|__  /____|   |___\_______  /\____|__  /
		 \/        \/        \/                        \/         \/                     \/         \/
```

This command is for destination (address, locking script) related commands.

new: creates a new destination in BUX (destination new <xpub>)
get: gets an existing destination in BUX (destination get <destination_id | address | locking_script> -x=<xpub_id>)


```
buxcli destination [flags]
```

### Examples

```
buxcli destination new <xpub>
```

### Options

```
  -h, --help              help for destination
  -m, --metadata string   Model Metadata
  -w, --woc               Optional flag to use WhatsOnChain for additional address data
  -x, --xpubid string     Xpub ID
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

